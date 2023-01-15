import { ColumnInfo, ConstraintInfo, NamedObjectInfo, TableInfo, SqlObjectInfo } from './dbinfo';

export interface AlterProcessor {
  createTable(table: TableInfo);
  dropTable(table: TableInfo);
  createColumn(column: ColumnInfo, constraints: ConstraintInfo[]);
  changeColumn(oldColumn: ColumnInfo, newColumn: ColumnInfo, constraints?: ConstraintInfo[]);
  dropColumn(column: ColumnInfo);
  createConstraint(constraint: ConstraintInfo);
  changeConstraint(oldConstraint: ConstraintInfo, newConstraint: ConstraintInfo);
  dropConstraint(constraint: ConstraintInfo);
  renameTable(table: TableInfo, newName: string);
  renameColumn(column: ColumnInfo, newName: string);
  renameConstraint(constraint: ConstraintInfo, newName: string);
  recreateTable(oldTable: TableInfo, newTable: TableInfo);
  createSqlObject(obj: SqlObjectInfo);
  dropSqlObject(obj: SqlObjectInfo);
  fillPreloadedRows(table: NamedObjectInfo, oldRows: any[], newRows: any[], key: string[], insertOnly: string[]);
}
