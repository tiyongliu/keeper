import {fromPairs} from 'lodash-es'
import openNewTab from '/@/second/utility/openNewTab'
export default function openReferenceForm(rowData, column, conid, database) {
  const formViewKey = fromPairs(
    column.foreignKey.columns.map(({ refColumnName, columnName }) => [refColumnName, rowData[columnName]])
  );
  void openNewTab(
    {
      title: column.foreignKey.refTableName,
      icon: 'img table',
      tabComponent: 'TableDataTab',
      props: {
        schemaName: column.foreignKey.refSchemaName,
        pureName: column.foreignKey.refTableName,
        conid,
        database,
        objectTypeField: 'tables',
      },
    },
    {
      grid: {
        isFormView: true,
        formViewKey,
      },
    },
    {
      forceNewTab: true,
    }
  );
}

export function openPrimaryKeyForm(rowData, baseTable, conid, database) {
  const formViewKey = fromPairs(
    baseTable.primaryKey.columns.map(({ columnName }) => [columnName, rowData[columnName]])
  );
  void openNewTab(
    {
      title: baseTable.pureName,
      icon: 'img table',
      tabComponent: 'TableDataTab',
      props: {
        schemaName: baseTable.schemaName,
        pureName: baseTable.pureName,
        conid,
        database,
        objectTypeField: 'tables',
      },
    },
    {
      grid: {
        isFormView: true,
        formViewKey,
      },
    },
    {
      forceNewTab: true,
    }
  );
}
