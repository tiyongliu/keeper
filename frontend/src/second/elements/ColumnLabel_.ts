export function _getColumnIcon(column, forceIcon = false) {
  if (column.autoIncrement) return 'img autoincrement';
  if (column.foreignKey) return 'img foreign-key';
  if (forceIcon) return 'img column';
  return null;
}
