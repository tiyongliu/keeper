import {Component, defineComponent, PropType, ref, toRefs, unref, watch} from 'vue'
import {cloneDeepWith, zipObject} from 'lodash-es'
import LoadingDataGridCore from './LoadingDataGridCore'
import {GridConfig, GridDisplay, MacroDefinition} from '/@/second/keeper-datalib'
import {databaseConnectionsCollectionDataApi} from '/@/api/simpleApis'
import {parseFilter} from '/@/second/keeper-filterparser'
import ChangeSetGrider from './ChangeSetGrider'

function buildGridMongoCondition(props) {
  const filters = props?.display?.config?.filters;

  const conditions: any[] = [];
  for (const uniqueName in filters || {}) {
    if (!filters[uniqueName]) continue;
    try {
      const ast = parseFilter(filters[uniqueName], 'mongo');
      // console.log('AST', ast);
      const cond = cloneDeepWith(ast, expr => {
        if (expr.__placeholder__) {
          return {
            [uniqueName]: expr.__placeholder__,
          };
        }
      });
      conditions.push(cond);
    } catch (err) {
      // error in filter
    }
  }

  return conditions.length > 0
    ? {
      $and: conditions,
    }
    : undefined;
}

function buildMongoSort(props) {
  const sort = props?.display?.config?.sort;

  if (sort?.length > 0) {
    return zipObject(
      sort.map(col => col.uniqueName),
      sort.map(col => (col.order == 'DESC' ? -1 : 1))
    );
  }

  return null;
}

export async function loadCollectionDataPage(props, offset, limit): Promise<any[]> {
  const {conid, database} = props;
  const response = await databaseConnectionsCollectionDataApi({
    conid,
    database,
    options: {
      pureName: props.pureName,
      limit,
      skip: offset,
      condition: buildGridMongoCondition(props),
      sort: buildMongoSort(props),
    },
  }) as any

  if (response.errorMessage) return response
  return response.rows;
}

function dataPageAvailable() {
  return true;
  // const { display } = props;
  // const sql = display.getPageQuery(0, 1);
  // return !!sql;
}

async function loadRowCount(props) {
  const {conid, database} = props;
  const response = await databaseConnectionsCollectionDataApi({
    conid,
    database,
    options: {
      pureName: props.pureName,
      countDocuments: true,
      condition: buildGridMongoCondition(props),
    },
  }) as any

  return response.count;
}

export default defineComponent({
  name: 'CollectionDataGridCore',
  props: {
    conid: {
      type: String as PropType<string>
    },
    display: {
      type: Object as PropType<GridDisplay>
    },
    database: {
      type: String as PropType<string>
    },
    schemaName: {
      type: String as PropType<string>
    },
    pureName: {
      type: String as PropType<string>
    },
    config: {
      type: Object as PropType<GridConfig>,
    },
    changeSetState: {
      type: Object as PropType<any>
    },
    dispatchChangeSet: {
      type: Function as PropType<(action: any) => void>
    },
    macroPreview: {
      type: [String, Object] as PropType<string | Component | MacroDefinition>,
    },
    macroValues: {
      type: Object as PropType<any>
    },
    selectedCellsPublished: {
      type: Function as PropType<() => []>,
      default: () => []
    },
    loadedRows: {
      type: Array as PropType<any[]>,
      default: () => []
    },
  },
  emits: ['update:loadedRows', 'update:selectedCellsPublished'],
  setup(props, {attrs, emit}) {
    const {
      display,
      macroPreview,
      macroValues,
      loadedRows,
      changeSetState,
      dispatchChangeSet,
      selectedCellsPublished
    } = toRefs(props)

    const grider = ref()
    const loadedRowsRW = ref(loadedRows.value)
    const selectedCellsPublishedRW = ref(selectedCellsPublished.value)

    watch(() => [
      ...loadedRowsRW.value,
      changeSetState.value,
      dispatchChangeSet.value,
      display.value!,
      macroPreview.value,
      macroValues.value,
      selectedCellsPublished.value,
    ], () => {
      grider.value = new ChangeSetGrider(
        loadedRowsRW.value,
        changeSetState.value,
        dispatchChangeSet.value,
        display.value!,
        macroPreview.value as Nullable<MacroDefinition>,
        macroValues.value,
        selectedCellsPublished.value()
      )
    })

    watch(() => [...loadedRowsRW.value], () => {
      emit('update:loadedRows', unref(loadedRowsRW.value))
    })

    watch(() => selectedCellsPublishedRW.value, () => {
      emit('update:selectedCellsPublished', selectedCellsPublishedRW.value)
    })

    return () => (
      <LoadingDataGridCore
        {...Object.assign({}, props, attrs)}
        loadDataPage={loadCollectionDataPage}
        dataPageAvailable={dataPageAvailable}
        loadRowCount={loadRowCount}
        vModel:loadedRows={loadedRowsRW.value}
        vModel:selectedCellsPublished={selectedCellsPublishedRW.value}
        frameSelection={!!macroPreview.value}
        grider={grider.value}
      />
    )
  }
})

