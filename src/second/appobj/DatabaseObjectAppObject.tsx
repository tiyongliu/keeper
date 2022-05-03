import {computed, defineComponent, onMounted, PropType, unref, watch, ref, toRefs} from 'vue'
import {isNaN} from 'lodash-es'
import AppObjectCore from '/@/second/appobj/AppObjectCore.vue'
import {filterName} from '/@/packages/tools/src'
import { dataBaseStore } from "/@/store/modules/dataBase";

export const extractKey = ({ schemaName, pureName }) => (schemaName ? `${schemaName}.${pureName}` : pureName);
export const createMatcher = ({ schemaName, pureName }) => filter => filterName(filter, pureName, schemaName);

export const databaseObjectIcons = {
  tables: 'img table',
  collections: 'img collection',
  views: 'img view',
  matviews: 'img view',
  procedures: 'img procedure',
  functions: 'img function',
  queries: 'img query-data',
}

const defaultTabs = {
  tables: 'TableDataTab',
  collections: 'CollectionDataTab',
  views: 'ViewDataTab',
  matviews: 'ViewDataTab',
  queries: 'QueryDataTab',
};

const menus = {
  tables: [
    {
      label: 'Open data',
      tab: 'TableDataTab',
      forceNewTab: true,
    },
    {
      label: 'Open form',
      tab: 'TableDataTab',
      forceNewTab: true,
      initialData: {
        grid: {
          isFormView: true,
        },
      },
    },
    {
      label: 'Open structure',
      tab: 'TableStructureTab',
      icon: 'img table-structure',
    },
    {
      label: 'Drop table',
      isDrop: true,
    },
    {
      label: 'Rename table',
      isRename: true,
    },
    {
      label: 'Query designer',
      isQueryDesigner: true,
    },
    {
      label: 'Show diagram',
      isDiagram: true,
    },
    {
      divider: true,
    },
    {
      label: 'Export',
      functionName: 'tableReader',
      isExport: true,
    },
    {
      label: 'Open as data sheet',
      isOpenFreeTable: true,
    },
    {
      label: 'Open active chart',
      isActiveChart: true,
    },
    {
      divider: true,
    },
    {
      label: 'SQL: CREATE TABLE',
      scriptTemplate: 'CREATE TABLE',
    },
    {
      label: 'SQL: SELECT',
      scriptTemplate: 'SELECT',
    },
    {
      label: 'SQL Generator: CREATE TABLE',
      sqlGeneratorProps: {
        createTables: true,
        createIndexes: true,
        createForeignKeys: true,
      },
    },
    {
      label: 'SQL Generator: DROP TABLE',
      sqlGeneratorProps: {
        dropTables: true,
        dropReferences: true,
      },
    },
    {
      label: 'SQL Generator: INSERT',
      sqlGeneratorProps: {
        insert: true,
      },
    },
  ],
  views: [
    {
      label: 'Open data',
      tab: 'ViewDataTab',
      forceNewTab: true,
    },
    {
      label: 'Open structure',
      tab: 'TableStructureTab',
      icon: 'img view-structure',
    },
    {
      label: 'Drop view',
      isDrop: true,
    },
    {
      label: 'Query designer',
      isQueryDesigner: true,
    },
    {
      divider: true,
    },
    {
      label: 'Export',
      isExport: true,
      functionName: 'tableReader',
    },
    {
      label: 'Open as data sheet',
      isOpenFreeTable: true,
    },
    {
      label: 'Open active chart',
      isActiveChart: true,
    },
    {
      divider: true,
    },
    {
      label: 'SQL: CREATE VIEW',
      scriptTemplate: 'CREATE OBJECT',
    },
    {
      label: 'SQL: CREATE TABLE',
      scriptTemplate: 'CREATE TABLE',
    },
    {
      label: 'SQL: SELECT',
      scriptTemplate: 'SELECT',
    },
    {
      label: 'SQL Generator: CREATE VIEW',
      sqlGeneratorProps: {
        createViews: true,
      },
    },
    {
      label: 'SQL Generator: DROP VIEW',
      sqlGeneratorProps: {
        dropViews: true,
      },
    },
  ],
  matviews: [
    {
      label: 'Open data',
      tab: 'ViewDataTab',
      forceNewTab: true,
    },
    {
      label: 'Open structure',
      tab: 'TableStructureTab',
    },
    {
      label: 'Drop view',
      isDrop: true,
    },
    {
      label: 'Query designer',
      isQueryDesigner: true,
    },
    {
      divider: true,
    },
    {
      label: 'Export',
      isExport: true,
      functionName: 'tableReader',
    },
    {
      label: 'Open as data sheet',
      isOpenFreeTable: true,
    },
    {
      label: 'Open active chart',
      isActiveChart: true,
    },
    {
      divider: true,
    },
    {
      label: 'SQL: CREATE MATERIALIZED VIEW',
      scriptTemplate: 'CREATE OBJECT',
    },
    {
      label: 'SQL: CREATE TABLE',
      scriptTemplate: 'CREATE TABLE',
    },
    {
      label: 'SQL: SELECT',
      scriptTemplate: 'SELECT',
    },
    {
      label: 'SQL Generator: CREATE MATERIALIZED VIEW',
      sqlGeneratorProps: {
        createMatviews: true,
      },
    },
    {
      label: 'SQL Generator: DROP MATERIALIZED VIEW',
      sqlGeneratorProps: {
        dropMatviews: true,
      },
    },
  ],
  queries: [
    {
      label: 'Open data',
      tab: 'QueryDataTab',
      forceNewTab: true,
    },
  ],
  procedures: [
    {
      label: 'Drop procedure',
      isDrop: true,
    },
    {
      label: 'SQL: CREATE PROCEDURE',
      scriptTemplate: 'CREATE OBJECT',
    },
    {
      label: 'SQL: EXECUTE',
      scriptTemplate: 'EXECUTE PROCEDURE',
    },
    {
      label: 'SQL Generator: CREATE PROCEDURE',
      sqlGeneratorProps: {
        createProcedures: true,
      },
    },
    {
      label: 'SQL Generator: DROP PROCEDURE',
      sqlGeneratorProps: {
        dropProcedures: true,
      },
    },
  ],
  functions: [
    {
      label: 'Drop function',
      isDrop: true,
    },
    {
      label: 'SQL: CREATE FUNCTION',
      scriptTemplate: 'CREATE OBJECT',
    },
    {
      label: 'SQL Generator: CREATE FUNCTION',
      sqlGeneratorProps: {
        createFunctions: true,
      },
    },
    {
      label: 'SQL Generator: DROP FUNCTION',
      sqlGeneratorProps: {
        dropFunctions: true,
      },
    },
  ],
  collections: [
    {
      label: 'Open data',
      tab: 'CollectionDataTab',
      forceNewTab: true,
    },
    {
      label: 'Open JSON',
      tab: 'CollectionDataTab',
      forceNewTab: true,
      initialData: {
        grid: {
          isJsonView: true,
        },
      },
    },
    {
      label: 'Export',
      isExport: true,
      functionName: 'tableReader',
    },
    {
      label: 'Drop collection',
      isDropCollection: true,
    },
    {
      label: 'Rename collection',
      isRenameCollection: true,
    },
    {
      divider: true,
    },
    {
      label: 'JS: dropCollection()',
      scriptTemplate: 'dropCollection',
    },
    {
      label: 'JS: find()',
      scriptTemplate: 'findCollection',
    },
  ],
};

function testEqual(a, b) {
  return (
    a.conid == b.conid &&
    a.database == b.database &&
    a.objectTypeField == b.objectTypeField &&
    a.pureName == b.pureName &&
    a.schemaName == b.schemaName
  );
}

function formatRowCount(value) {
  const num = parseInt(value);
  if (isNaN(num)) return value;
  return num.toLocaleString();
}

export default defineComponent({
  name: 'DatabaseObjectAppObject',
  props: {
    data: {
      type: Object as PropType<{name: string, schemaName: string, objectTypeField: string, pureName: string, tableRowCount?: null}>,
    },
    passProps: {
      type: Object as PropType<{
        showPinnedInsteadOfUnpin: boolean
      }>,
    },
  },
  components: {
    AppObjectCore
  },
  setup(props, {attrs}) {
    const dataBase = dataBaseStore()

    const {data, passProps} = toRefs(props)

    const isPinned = computed(() => !!dataBase.getPinnedTables.find(x => testEqual(data.value, x)))

    function handleClick() {
      handleDatabaseObjectClick()
    }

    return () => <AppObjectCore
      {...attrs}
      data={data.value}
      title={unref(data)!.schemaName ? `${unref(data)!.schemaName}.${unref(data)!.pureName}` : unref(data)!.pureName }
      icon={databaseObjectIcons[data.value!.objectTypeField]}
      showPinnedInsteadOfUnpin={passProps.value?.showPinnedInsteadOfUnpin}
      onPin={unref(isPinned) ? null : () => dataBase.subscribePinnedTables([
        ...unref(dataBase.$state.pinnedTables),
        unref(data)!
      ])}
      onUnpin={unref(isPinned) ? () => dataBase.subscribePinnedTables(
        unref(dataBase.$state.pinnedTables).filter(x => !testEqual(x, data.value))
      ) : null}
      extInfo={unref(data)!.tableRowCount != null ? `${formatRowCount(unref(data)!.tableRowCount)} rows` : null}
      onClick={() => handleClick()}
    />
  }
})


export async function openDatabaseObjectDetail() {
}

export function handleDatabaseObjectClick() {

}
