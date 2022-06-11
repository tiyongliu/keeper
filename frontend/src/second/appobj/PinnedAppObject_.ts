export function getDatabaseMenuItems() {
  const handleNewQuery = () => {
  }

  const handleNewTable = () => {
  }

  const handleNewCollection = () => {
  }

  const handleImport = () => {

  }
  const handleExport = () => {

  }

  const handleSqlRestore = () => {
  }

  const handleSqlDump = () => {

  }

  return [
    {onClick: handleNewQuery, text: 'New query', isNewQuery: true},
    {onClick: handleNewTable, text: 'New table'},
    {onClick: handleNewCollection, text: 'New collection'},
    {onClick: handleImport, text: 'Import wizard'},
    {onClick: handleExport, text: 'Export wizard'},
    {onClick: handleSqlRestore, text: 'Restore/import SQL dump'},
    {onClick: handleSqlDump, text: 'Backup/export SQL dump'},
    { divider: true },
  ]
}
