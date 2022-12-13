import {defineComponent, PropType, toRefs} from 'vue'
import {TableFormViewDisplay} from "/@/second/keeper-datalib";
import {ColumnReference, EngineDriver} from "/@/second/keeper-types";
import FontIcon from '/@/second/icons/FontIcon.vue'
import ColumnLabel from '/@/second/elements/ColumnLabel.vue'
import InlineButton from '/@/second/buttons/InlineButton.vue'
import DataFilterControl from '/@/second/datagrid/DataFilterControl.vue'
import {getFilterType} from '/@/second/keeper-filterparser'

export default defineComponent({
  name: 'FormViewFilterColumn',
  props: {
    column: {
      type: Object as PropType<ColumnReference>,
    },
    formDisplay: {
      type: Object as PropType<TableFormViewDisplay>
    },
    filters: {
      type: Object as PropType<{ [uniqueName: string]: string }>
    },
    driver: {
      type: Object as PropType<EngineDriver>
    },
    conid: {
      type: String as PropType<string>
    },
    database: {
      type: String as PropType<string>
    },
    schemaName: {
      type: String as PropType<string>
    },
    pureName: {
      type: String as PropType<string>
    }
  },
  setup(props) {
    const {
      column,
      formDisplay,
      filters,
      driver,
      conid,
      database,
      schemaName,
      pureName
    } = toRefs(props)

    return () => (
      column.value && <div class="m-1">
        <div class="space-between">
          <ColumnLabel {...column.value}/>
          <InlineButton
            square
            narrow
            onClick={() => {
              formDisplay.value && formDisplay.value.removeFilter(column.value!['uniqueName'])
            }}>
            <FontIcon icon="icon close"/>
          </InlineButton>
        </div>
        <DataFilterControl
          filterType={getFilterType(column.value!['dataType'])}
          filter={filters.value![column.value['uniqueName']]}
          setFilter={value => formDisplay.value && formDisplay.value.setFilter(column.value!['uniqueName'], value)}
          driver={driver.value}
          conid={conid.value}
          database={database.value}
          schemaName={schemaName.value}
          pureName={pureName.value}
          columnName={column.value['uniquePath'].length == 1 ? column.value['uniquePath'][0] : null}
          foreignKey={column.value['foreignKey']}
        />
      </div>
    )
  }
})
