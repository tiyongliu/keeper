import {TableInfo, ColumnInfo} from '/@/second/keeper-types'

export function findForeignKeyForColumn(table: TableInfo, column: ColumnInfo) {
  return (table.foreignKeys || []).find(fk => fk.columns.find(col => col.columnName == column.columnName));
}
