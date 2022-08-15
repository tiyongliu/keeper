const driver = {
  engine: 'mongo@dbgate-plugin-mongo',
  title: 'MongoDB',
  editorMode: 'javascript',
  defaultPort: 27017,
  supportsDatabaseUrl: true,
  databaseUrlPlaceholder: 'e.g. mongodb://username:password@mongodb.mydomain.net/dbname',
  showConnectionField: (field: string, values) => {
    if (field == 'useDatabaseUrl') return true;
    if (values.useDatabaseUrl) {
      return ['databaseUrl', 'defaultDatabase', 'singleDatabase', 'isReadOnly'].includes(field);
    }
    return ['server', 'port', 'user', 'password', 'defaultDatabase', 'singleDatabase', 'isReadOnly'].includes(field);
  },
  importExportArgs: [
    {
      type: 'checkbox',
      name: 'createStringId',
      label: 'Create string _id attribute',
      apiName: 'createStringId',
      direction: 'target',
    },
  ],

}

export default driver
