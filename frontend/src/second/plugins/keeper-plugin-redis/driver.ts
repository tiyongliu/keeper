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
  dialect,
  engine: 'redis@dbgate-plugin-redis',
  title: 'Redis',
  defaultPort: 6379,
  editorMode: 'text',
  databaseEngineTypes: ['keyvalue'],
  supportedCreateDatabase: false,
  databaseUrlPlaceholder: 'e.g. redis://:authpassword@127.0.0.1:6380/4',
  supportedKeyTypes: [
    {
      name: 'string',
      label: 'String',
      dbKeyFields: [{ name: 'value' }],
      addMethod: 'set',
    },
    {
      name: 'list',
      label: 'List',
      dbKeyFields: [{ name: 'value' }],
      addMethod: 'rpush',
      showItemList: true,
    },
    {
      name: 'set',
      label: 'Set',
      dbKeyFields: [{ name: 'value' }],
      keyColumn: 'value',
      addMethod: 'sadd',
      showItemList: true,
    },
    {
      name: 'zset',
      label: 'Sorted Set',
      dbKeyFields: [{ name: 'score' }, { name: 'value' }],
      keyColumn: 'value',
      addMethod: 'zadd',
      showItemList: true,
    },
    {
      name: 'hash',
      label: 'Hash',
      dbKeyFields: [{ name: 'key' }, { name: 'value' }],
      keyColumn: 'key',
      addMethod: 'hset',
      showItemList: true,
    },
    {
      name: 'stream',
      label: 'Stream',
      dbKeyFields: [{ name: 'id' }, { name: 'value' }],
      keyColumn: 'id',
      addMethod: 'xaddjson',
      showItemList: true,
    },
  ],

  showConnectionField: (field, values) => {
    if (field == 'useDatabaseUrl') return true;
    if (values.useDatabaseUrl) {
      return ['databaseUrl', 'isReadOnly'].includes(field);
    }
    return ['server', 'port', 'password', 'isReadOnly'].includes(field);
  },
};

export default driver
