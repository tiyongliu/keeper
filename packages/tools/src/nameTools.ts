import {TableInfo, ColumnInfo} from 'dbbox-types'

export function findForeignKeyForColumn(table: TableInfo, column: ColumnInfo) {
  return (table.foreignKeys || []).find(fk => fk.columns.find(col => col.columnName == column.columnName));
}
