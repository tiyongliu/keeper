import {flatten} from 'lodash-es'
import {DictionaryDescription} from '/@/second/keeper-datalib'
import {ApplicationDefinition, TableInfo} from '/@/second/keeper-types'
import {filterAppsForDatabase} from './appTools'

function checkDescriptionColumns(columns: string[], table: TableInfo) {
  if (!columns?.length) return false;
  if (!columns.every(x => table.columns.find(y => y.columnName == x))) return false;
  if (table.primaryKey?.columns?.find(x => columns.includes(x.columnName))) return false;
  return true;
}

export function getDictionaryDescription(
  table: TableInfo,
  conid: string,
  database: string,
  apps: ApplicationDefinition[],
  connections,
  skipCheckSaved = false
): DictionaryDescription | null {
  const conn = connections?.find(x => x._id == conid);

  if (!conn) {
    return null;
  }

  const dbApps = filterAppsForDatabase(conn, database, apps);

  if (!dbApps) {
    return null;
  }

  const cached = flatten(dbApps.map(x => x.dictionaryDescriptions || [])).find(
    x => x.pureName == table.pureName && x.schemaName == table.schemaName
  );

  if (cached && (skipCheckSaved || checkDescriptionColumns(cached.columns, table))) {
    return cached;
  }

  const descColumn = table.columns.find(x => x?.dataType?.toLowerCase()?.includes('char'));
  if (descColumn) {
    return {
      columns: [descColumn.columnName],
      delimiter: null,
      expression: descColumn.columnName,
    };
  }

  return null;
}
