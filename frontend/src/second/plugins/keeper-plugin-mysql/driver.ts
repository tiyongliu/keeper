
/** @type {import('dbgate-types').SqlDialect} */
const dialect = {
  rangeSelect: true,
  stringEscapeChar: '\\',
  fallbackDataType: 'longtext',
  enableConstraintsPerTable: false,
  anonymousPrimaryKey: true,
  explicitDropConstraint: true,
  quoteIdentifier(s) {
    return '`' + s + '`';
  },

  createColumn: true,
  dropColumn: true,
  changeColumn: true,
  createIndex: true,
  dropIndex: true,
  createForeignKey: true,
  dropForeignKey: true,
  createPrimaryKey: true,
  dropPrimaryKey: true,
  dropIndexContainsTableSpec: true,
  createUnique: true,
  dropUnique: true,
  createCheck: true,
  dropCheck: true,

  dropReferencesWhenDropTable: false,

  columnProperties: {
    columnComment: true,
    isUnsigned: true,
    isZerofill: true,
  },
};

const mysqlDriverBase = {
  showConnectionField: (field, values) =>
    ['server', 'port', 'user', 'password', 'defaultDatabase', 'singleDatabase', 'isReadOnly'].includes(field),
  dialect,
  defaultPort: 3306,
  readOnlySessions: true,
  supportsDatabaseDump: true,

  getNewObjectTemplates() {
    return [
      { label: 'New view', sql: 'CREATE VIEW myview\nAS\nSELECT * FROM table1' },
      {
        label: 'New procedure',
        sql:
          'DELIMITER //\n\nCREATE PROCEDURE myproc (IN arg1 INT)\nBEGIN\n  SELECT * FROM table1;\nEND\n\nDELIMITER ;',
      },
      { label: 'New function', sql: 'CREATE FUNCTION myfunc (arg1 INT)\nRETURNS INT DETERMINISTIC\nRETURN 1' },
    ];
  },
};

/** @type {import('dbgate-types').EngineDriver} */
const mysqlDriver = {
  ...mysqlDriverBase,
  engine: 'mysql@dbgate-pluginMysql',
  title: 'MySQL',
};

/** @type {import('dbgate-types').EngineDriver} */
 const mariaDriver = {
  ...mysqlDriverBase,
  engine: 'mariadb@dbgate-pluginMysql',
  title: 'MariaDB',
}

export default [mysqlDriver, mariaDriver]
