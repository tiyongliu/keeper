function jsonStringifyWithObjectId(obj) {
  return JSON.stringify(obj, undefined, 2).replace(
    /\{\s*\"\$oid\"\s*\:\s*\"([0-9a-f]+)\"\s*\}/g,
    (m, id) => `ObjectId("${id}")`
  );
}

/** @type {import('dbgate-types').SqlDialect} */
const dialect = {
  limitSelect: true,
  rangeSelect: true,
  offsetFetchRangeSyntax: true,
  stringEscapeChar: "'",
  fallbackDataType: 'nvarchar(max)',
  quoteIdentifier(s) {
    return `[${s}]`;
  },
};

/** @type {import('dbgate-types').EngineDriver} */
const driver = {
  databaseEngineTypes: ['document'],
  dialect,
  engine: 'mongo',
  title: 'MongoDB',
  editorMode: 'javascript',
  defaultPort: 27017,
  supportsDatabaseUrl: true,
  databaseUrlPlaceholder: 'e.g. mongodb://username:password@mongodb.mydomain.net/dbname',

  getQuerySplitterOptions: () => mongoSplitterOptions,

  showConnectionField: (field, values) => {
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

  getCollectionUpdateScript(changeSet) {
    let res = '';
    for (const insert of changeSet.inserts) {
      res += `db.${insert.pureName}.insert(${jsonStringifyWithObjectId({
        ...insert.document,
        ...insert.fields,
      })});\n`;
    }
    for (const update of changeSet.updates) {
      if (update.document) {
        res += `db.${update.pureName}.replaceOne(${jsonStringifyWithObjectId(
          update.condition
        )}, ${jsonStringifyWithObjectId({
          ...update.document,
          ...update.fields,
        })});\n`;
      } else {
        res += `db.${update.pureName}.updateOne(${jsonStringifyWithObjectId(
          update.condition
        )}, ${jsonStringifyWithObjectId({
          $set: update.fields,
        })});\n`;
      }
    }
    for (const del of changeSet.deletes) {
      res += `db.${del.pureName}.deleteOne(${jsonStringifyWithObjectId(del.condition)});\n`;
    }
    return res;
  },
};

export default driver;
