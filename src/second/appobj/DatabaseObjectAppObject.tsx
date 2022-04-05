import {computed, defineComponent, onMounted, PropType, unref, watch, ref} from 'vue'
import AppObjectCore from '/@/second/appobj/AppObjectCore.vue'
import {filterName} from '/@/packages/tools/src'

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

export default defineComponent({
  name: 'DatabaseObjectAppObject',
  props: {
    data: {
      type: Object as PropType<{name: string}>,
      default: {"name":"performance_schema","connection":{"server":"localhost","engine":"mysql@dbgate-plugin-mysql","sshMode":"userPassword","sshPort":"22","sshKeyfile":"C:\\Users\\Administrator\\.ssh\\id_rsa","port":"5001","user":"root","password":"crypt:6f4a500c408bce8261606389954de288bdb5d66c02d0e8b1ffe4275cba54d24a5942f62d29ae1703d1f9de29e468af72adlfTFoqebJEvRNqKoigxQ==","_id":"b9c51b10-b354-11ec-812a-3d58c681a37b","status":{"name":"ok"}}}
    }
  },
  components: {
    AppObjectCore
  },
  setup(props) {
    const {data} = props
    return () => <AppObjectCore
      {...props}
      title={unref(data).name}
      icon="img database"
    />
  }
})
