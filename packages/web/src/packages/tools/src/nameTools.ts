import {TableInfo, ColumnInfo} from '../types/dbinfo'

export function findForeignKeyForColumn(table: TableInfo, column: ColumnInfo) {
  return (table.foreignKeys || []).find(fk => fk.columns.find(col => col.columnName == column.columnName));
}
