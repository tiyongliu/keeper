export function buildExtensions(plugins = [
                                  {
                                    "name": "dbgate-plugin-csv",
                                    "main": "dist/backend.js",
                                    "version": "4.1.1",
                                    "homepage": "https://dbgate.org",
                                    "description": "CSV import/export plugin for DbGate",
                                    "repository": {
                                      "type": "git",
                                      "url": "https://github.com/dbgate/dbgate"
                                    },
                                    "author": "Jan Prochazka",
                                    "license": "MIT",
                                    "keywords": [
                                      "csv",
                                      "import",
                                      "export",
                                      "dbgate",
                                      "dbgateplugin"
                                    ],
                                    "files": [
                                      "dist",
                                      "icon.svg"
                                    ],
                                    "scripts": {
                                      "build:frontend": "webpack --config webpack-frontend.config",
                                      "build:frontend:watch": "webpack --watch --config webpack-frontend.config",
                                      "build:backend": "webpack --config webpack-backend.config.js",
                                      "build": "yarn build:frontend && yarn build:backend",
                                      "plugin": "yarn build && yarn pack && dbgate-plugin dbgate-plugin-csv",
                                      "copydist": "yarn build && yarn pack && dbgate-copydist ../dist/dbgate-plugin-csv",
                                      "plugout": "dbgate-plugout dbgate-plugin-csv",
                                      "prepublishOnly": "yarn build"
                                    },
                                    "devDependencies": {
                                      "csv": "^5.3.2",
                                      "dbgate-plugin-tools": "^1.0.7",
                                      "lodash": "^4.17.21",
                                      "webpack": "^4.42.0",
                                      "webpack-cli": "^3.3.11"
                                    },
                                    "readme": "[![styled with prettier](https://img.shields.io/badge/styled_with-prettier-ff69b4.svg)](https://github.com/prettier/prettier)\r\n[![NPM version](https://img.shields.io/npm/v/dbgate-plugin-csv.svg)](https://www.npmjs.com/package/dbgate-plugin-csv)\r\n\r\n# dbgate-plugin-csv\r\n\r\nCSV import/export plugin for DbGate\r\n\r\n## Usage without DbGate\r\n\r\nExport from fake object reader into CSV file. Fake object file can be replaced with other reader/writer factory functions, as described in \r\n[dbgate-api package](https://www.npmjs.com/package/dbgate-api)\r\n\r\n```javascript\r\nconst dbgateApi = require('dbgate-api');\r\nconst dbgatePluginCsv = require(\"dbgate-plugin-csv\");\r\n\r\ndbgateApi.registerPlugins(dbgatePluginCsv);\r\n\r\n\r\nasync function run() {\r\n  const reader = await dbgateApi.fakeObjectReader();\r\n  const writer = await dbgatePluginCsv.shellApi.writer({ fileName: 'myfile1.csv', separator: ';' });\r\n  await dbgateApi.copyStream(reader, writer);\r\n  \r\n  console.log('Finished job script');\r\n}\r\ndbgateApi.runScript(run);\r\n\r\n\r\n```\r\n\r\n## Factory functions\r\n\r\n### shellApi.reader\r\nReads CSV file\r\n```js\r\n  const dbgatePluginCsv = require(\"dbgate-plugin-csv\");\r\n  const reader = await dbgatePluginCsv.shellApi.reader({\r\n    fileName: 'test.csv',\r\n    encoding: 'utf-8',\r\n    header: true,\r\n    delimiter: ',',\r\n    quoted: false,\r\n    limitRows: null\r\n  });\r\n```\r\n\r\n### shellApi.writer\r\nWrites CSV file\r\n```js\r\n  const dbgatePluginCsv = require(\"dbgate-plugin-csv\");\r\n  const writer = await dbgatePluginCsv.shellApi.writer({\r\n    fileName: 'test.csv',\r\n    encoding: 'utf-8',\r\n    header: true,\r\n    delimiter: ',',\r\n    quoted: false\r\n  });\r\n```\r\n",
                                    "isPackaged": true
                                  },
                                  {
                                    "name": "dbgate-plugin-excel",
                                    "main": "dist/backend.js",
                                    "version": "4.1.1",
                                    "description": "MS Excel import/export plugin for DbGate",
                                    "homepage": "https://dbgate.org",
                                    "repository": {
                                      "type": "git",
                                      "url": "https://github.com/dbgate/dbgate"
                                    },
                                    "author": "Jan Prochazka",
                                    "license": "MIT",
                                    "keywords": [
                                      "excel",
                                      "import",
                                      "export",
                                      "dbgate",
                                      "dbgateplugin"
                                    ],
                                    "files": [
                                      "dist",
                                      "icon.svg"
                                    ],
                                    "scripts": {
                                      "build:frontend": "webpack --config webpack-frontend.config",
                                      "build:frontend:watch": "webpack --watch --config webpack-frontend.config",
                                      "build:backend": "webpack --config webpack-backend.config.js",
                                      "build": "yarn build:frontend && yarn build:backend",
                                      "plugin": "yarn build && dbgate-plugin dbgate-plugin-excel",
                                      "copydist": "yarn build && yarn pack && dbgate-copydist ../dist/dbgate-plugin-excel",
                                      "plugout": "dbgate-plugout dbgate-plugin-excel",
                                      "prepublishOnly": "yarn build"
                                    },
                                    "devDependencies": {
                                      "lodash": "^4.17.21",
                                      "xlsx": "^0.16.8",
                                      "dbgate-plugin-tools": "^1.0.7",
                                      "webpack": "^4.42.0",
                                      "webpack-cli": "^3.3.11"
                                    },
                                    "readme": "[![styled with prettier](https://img.shields.io/badge/styled_with-prettier-ff69b4.svg)](https://github.com/prettier/prettier)\r\n[![NPM version](https://img.shields.io/npm/v/dbgate-plugin-excel.svg)](https://www.npmjs.com/package/dbgate-plugin-excel)\r\n\r\n# dbgate-plugin-excel\r\n\r\nMS Excel import/export plugin for DbGate\r\n\r\n\r\n## Usage without DbGate\r\n\r\nExport from fake object reader into MS Excel file. Fake object file can be replaced with other reader/writer factory functions, as described in \r\n[dbgate-api package](https://www.npmjs.com/package/dbgate-api)\r\n\r\n```javascript\r\nconst dbgateApi = require('dbgate-api');\r\nconst dbgatePluginExcel = require(\"dbgate-plugin-excel\");\r\n\r\ndbgateApi.registerPlugins(dbgatePluginExcel);\r\n\r\n\r\nasync function run() {\r\n  const reader = await dbgateApi.fakeObjectReader();\r\n  const writer = await dbgatePluginExcel.shellApi.writer({ fileName: 'myfile1.xlsx', sheetName: 'Sheet 1' });\r\n  await dbgateApi.copyStream(reader, writer);\r\n  console.log('Finished job script');\r\n}\r\ndbgateApi.runScript(run);\r\n\r\n\r\n```\r\n\r\n## Factory functions\r\n\r\n### shellApi.reader\r\nReads tabular data from one sheet in MS Excel file.\r\n```js\r\n  const reader = await dbgatePluginExcel.shellApi.reader({\r\n    fileName: 'test.xlsx',\r\n    sheetName: 'Album',\r\n    limitRows: null\r\n  });\r\n```\r\n\r\n### shellApi.writer\r\nWrites tabular data into MS excel file. There could be more writes into the some file in one script, if property sheetName is different.\r\n```js\r\n  const reader = await dbgatePluginExcel.shellApi.writer({\r\n    fileName: 'test.xlsx',\r\n    sheetName: 'Album',\r\n  });\r\n```\r\n",
                                    "isPackaged": true
                                  },
                                  {
                                    "name": "dbgate-plugin-mongo",
                                    "main": "dist/backend.js",
                                    "version": "4.1.1",
                                    "license": "MIT",
                                    "author": "Jan Prochazka",
                                    "homepage": "https://dbgate.org",
                                    "description": "MongoDB connect plugin for DbGate",
                                    "repository": {
                                      "type": "git",
                                      "url": "https://github.com/dbgate/dbgate"
                                    },
                                    "keywords": [
                                      "dbgate",
                                      "dbgateplugin",
                                      "mongo",
                                      "mongodb"
                                    ],
                                    "files": [
                                      "dist",
                                      "icon.svg"
                                    ],
                                    "scripts": {
                                      "build:frontend": "webpack --config webpack-frontend.config",
                                      "build:frontend:watch": "webpack --watch --config webpack-frontend.config",
                                      "build:backend": "webpack --config webpack-backend.config.js",
                                      "build": "yarn build:frontend && yarn build:backend",
                                      "plugin": "yarn build && yarn pack && dbgate-plugin dbgate-plugin-mongo",
                                      "copydist": "yarn build && yarn pack && dbgate-copydist ../dist/dbgate-plugin-mongo",
                                      "plugout": "dbgate-plugout dbgate-plugin-mongo",
                                      "prepublishOnly": "yarn build"
                                    },
                                    "devDependencies": {
                                      "dbgate-plugin-tools": "^1.0.7",
                                      "dbgate-query-splitter": "^4.1.1",
                                      "webpack": "^4.42.0",
                                      "webpack-cli": "^3.3.11",
                                      "dbgate-tools": "^4.1.1",
                                      "is-promise": "^4.0.0",
                                      "mongodb": "^3.6.5",
                                      "mongodb-client-encryption": "^1.2.3"
                                    },
                                    "readme": "[![styled with prettier](https://img.shields.io/badge/styled_with-prettier-ff69b4.svg)](https://github.com/prettier/prettier)\r\n[![NPM version](https://img.shields.io/npm/v/dbgate-plugin-mongo.svg)](https://www.npmjs.com/package/dbgate-plugin-mongo)\r\n\r\n# dbgate-plugin-mongo\r\n\r\nUse DbGate for install of this plugin\r\n",
                                    "isPackaged": true
                                  },
                                  {
                                    "name": "dbgate-plugin-mssql",
                                    "main": "dist/backend.js",
                                    "version": "4.1.1",
                                    "homepage": "https://dbgate.org",
                                    "description": "MS SQL connect plugin for DbGate",
                                    "repository": {
                                      "type": "git",
                                      "url": "https://github.com/dbgate/dbgate"
                                    },
                                    "author": "Jan Prochazka",
                                    "license": "MIT",
                                    "keywords": [
                                      "sql",
                                      "mssql",
                                      "dbgate",
                                      "dbgateplugin"
                                    ],
                                    "files": [
                                      "dist",
                                      "icon.svg"
                                    ],
                                    "scripts": {
                                      "build:frontend": "webpack --config webpack-frontend.config",
                                      "build:frontend:watch": "webpack --watch --config webpack-frontend.config",
                                      "build:backend": "webpack --config webpack-backend.config.js",
                                      "build": "yarn build:frontend && yarn build:backend",
                                      "prepublishOnly": "yarn build",
                                      "plugin": "yarn build && yarn pack && dbgate-plugin dbgate-plugin-mssql",
                                      "copydist": "yarn build && yarn pack && dbgate-copydist ../dist/dbgate-plugin-mssql",
                                      "plugout": "dbgate-plugout dbgate-plugin-mssql"
                                    },
                                    "devDependencies": {
                                      "dbgate-plugin-tools": "^1.0.7",
                                      "dbgate-query-splitter": "^4.1.1",
                                      "webpack": "^4.42.0",
                                      "webpack-cli": "^3.3.11",
                                      "dbgate-tools": "^4.1.1",
                                      "tedious": "^9.2.3",
                                      "async-lock": "^1.2.6"
                                    },
                                    "readme": "[![styled with prettier](https://img.shields.io/badge/styled_with-prettier-ff69b4.svg)](https://github.com/prettier/prettier)\r\n[![NPM version](https://img.shields.io/npm/v/dbgate-plugin-mssql.svg)](https://www.npmjs.com/package/dbgate-plugin-mssql)\r\n\r\n# dbgate-plugin-mssql\r\n\r\nMS SQL connector plugin for DbGate\r\n",
                                    "isPackaged": true
                                  },
                                  {
                                    "name": "dbgate-plugin-mysql",
                                    "main": "dist/backend.js",
                                    "version": "4.1.1",
                                    "homepage": "https://dbgate.org",
                                    "description": "MySQL connect plugin for DbGate",
                                    "repository": {
                                      "type": "git",
                                      "url": "https://github.com/dbgate/dbgate"
                                    },
                                    "author": "Jan Prochazka",
                                    "license": "MIT",
                                    "keywords": [
                                      "sql",
                                      "dbgate",
                                      "dbgateplugin",
                                      "mysql"
                                    ],
                                    "files": [
                                      "dist",
                                      "icon.svg"
                                    ],
                                    "scripts": {
                                      "build:frontend": "webpack --config webpack-frontend.config",
                                      "build:frontend:watch": "webpack --watch --config webpack-frontend.config",
                                      "build:backend": "webpack --config webpack-backend.config.js",
                                      "build": "yarn build:frontend && yarn build:backend",
                                      "plugin": "yarn build && yarn pack && dbgate-plugin dbgate-plugin-mysql",
                                      "copydist": "yarn build && yarn pack && dbgate-copydist ../dist/dbgate-plugin-mysql",
                                      "plugout": "dbgate-plugout dbgate-plugin-mysql",
                                      "prepublishOnly": "yarn build"
                                    },
                                    "devDependencies": {
                                      "dbgate-plugin-tools": "^1.0.7",
                                      "dbgate-query-splitter": "^4.1.1",
                                      "webpack": "^4.42.0",
                                      "webpack-cli": "^3.3.11",
                                      "dbgate-tools": "^4.1.1",
                                      "mysql2": "^2.2.5"
                                    },
                                    "readme": "[![styled with prettier](https://img.shields.io/badge/styled_with-prettier-ff69b4.svg)](https://github.com/prettier/prettier)\r\n[![NPM version](https://img.shields.io/npm/v/dbgate-plugin-mysql.svg)](https://www.npmjs.com/package/dbgate-plugin-mysql)\r\n\r\n# dbgate-plugin-mysql\r\n\r\nUse DbGate for install of this plugin\r\n",
                                    "isPackaged": true
                                  },
                                  {
                                    "name": "dbgate-plugin-postgres",
                                    "main": "dist/backend.js",
                                    "version": "4.1.1",
                                    "license": "MIT",
                                    "description": "PostgreSQL connector plugin for DbGate",
                                    "homepage": "https://dbgate.org",
                                    "repository": {
                                      "type": "git",
                                      "url": "https://github.com/dbgate/dbgate"
                                    },
                                    "author": "Jan Prochazka",
                                    "keywords": [
                                      "dbgate",
                                      "dbgateplugin",
                                      "postgresql"
                                    ],
                                    "files": [
                                      "dist",
                                      "icon.svg"
                                    ],
                                    "scripts": {
                                      "build:frontend": "webpack --config webpack-frontend.config",
                                      "build:frontend:watch": "webpack --watch --config webpack-frontend.config",
                                      "build:backend": "webpack --config webpack-backend.config.js",
                                      "build": "yarn build:frontend && yarn build:backend",
                                      "plugin": "yarn build && yarn pack && dbgate-plugin dbgate-plugin-postgres",
                                      "copydist": "yarn build && yarn pack && dbgate-copydist ../dist/dbgate-plugin-postgres",
                                      "plugout": "dbgate-plugout dbgate-plugin-postgres",
                                      "prepublishOnly": "yarn build"
                                    },
                                    "devDependencies": {
                                      "dbgate-plugin-tools": "^1.0.7",
                                      "dbgate-query-splitter": "^4.1.1",
                                      "dbgate-tools": "^4.1.1",
                                      "lodash": "^4.17.21",
                                      "pg": "^8.7.1",
                                      "webpack": "^4.42.0",
                                      "webpack-cli": "^3.3.11"
                                    },
                                    "readme": "[![styled with prettier](https://img.shields.io/badge/styled_with-prettier-ff69b4.svg)](https://github.com/prettier/prettier)\r\n[![NPM version](https://img.shields.io/npm/v/dbgate-plugin-postgres.svg)](https://www.npmjs.com/package/dbgate-plugin-postgres)\r\n\r\n# dbgate-plugin-postgres\r\n\r\nUse DbGate for install of this plugin\r\n",
                                    "isPackaged": true
                                  },
                                  {
                                    "name": "dbgate-plugin-redis",
                                    "main": "dist/backend.js",
                                    "version": "1.0.0",
                                    "license": "MIT",
                                    "description": "Redis connector plugin for DbGate",
                                    "homepage": "https://dbgate.org",
                                    "repository": {
                                      "type": "git",
                                      "url": "https://github.com/dbgate/dbgate"
                                    },
                                    "author": "Jan Prochazka",
                                    "keywords": [
                                      "dbgate",
                                      "dbgateplugin",
                                      "redis"
                                    ],
                                    "files": [
                                      "dist"
                                    ],
                                    "scripts": {
                                      "build:frontend": "webpack --config webpack-frontend.config",
                                      "build:frontend:watch": "webpack --watch --config webpack-frontend.config",
                                      "build:backend": "webpack --config webpack-backend.config.js",
                                      "build": "yarn build:frontend && yarn build:backend",
                                      "plugin": "yarn build && yarn pack && dbgate-plugin dbgate-plugin-redis",
                                      "plugout": "dbgate-plugout dbgate-plugin-redis",
                                      "copydist": "yarn build && yarn pack && dbgate-copydist ../dist/dbgate-plugin-redis",
                                      "prepublishOnly": "yarn build"
                                    },
                                    "devDependencies": {
                                      "dbgate-plugin-tools": "^1.0.7",
                                      "dbgate-query-splitter": "^4.1.1",
                                      "dbgate-tools": "^4.1.1",
                                      "lodash": "^4.17.21",
                                      "webpack": "^4.42.0",
                                      "webpack-cli": "^3.3.11",
                                      "async": "^3.2.3",
                                      "ioredis": "^4.28.5",
                                      "node-redis-dump2": "^0.5.0"
                                    },
                                    "readme": "[![styled with prettier](https://img.shields.io/badge/styled_with-prettier-ff69b4.svg)](https://github.com/prettier/prettier)\r\n[![NPM version](https://img.shields.io/npm/v/dbgate-plugin-redis.svg)](https://www.npmjs.com/package/dbgate-plugin-redis)\r\n\r\n# dbgate-plugin-redis\r\n\r\nUse DbGate for install of this plugin\r\n",
                                    "isPackaged": true
                                  },
                                  {
                                    "name": "dbgate-plugin-sqlite",
                                    "main": "dist/backend.js",
                                    "version": "4.1.1",
                                    "homepage": "https://dbgate.org",
                                    "description": "SQLite connect plugin for DbGate",
                                    "repository": {
                                      "type": "git",
                                      "url": "https://github.com/dbgate/dbgate"
                                    },
                                    "author": "Jan Prochazka",
                                    "license": "MIT",
                                    "keywords": [
                                      "dbgate",
                                      "dbgateplugin",
                                      "sqlite"
                                    ],
                                    "files": [
                                      "dist",
                                      "icon.svg"
                                    ],
                                    "scripts": {
                                      "build:frontend": "webpack --config webpack-frontend.config",
                                      "build:frontend:watch": "webpack --watch --config webpack-frontend.config",
                                      "build:backend": "webpack --config webpack-backend.config.js",
                                      "build": "yarn build:frontend && yarn build:backend",
                                      "plugin": "yarn build && yarn pack && dbgate-plugin dbgate-plugin-sqlite",
                                      "copydist": "yarn build && yarn pack && dbgate-copydist ../dist/dbgate-plugin-sqlite",
                                      "plugout": "dbgate-plugout dbgate-plugin-sqlite",
                                      "prepublishOnly": "yarn build"
                                    },
                                    "devDependencies": {
                                      "dbgate-tools": "^4.1.1",
                                      "dbgate-plugin-tools": "^1.0.4",
                                      "dbgate-query-splitter": "^4.1.1",
                                      "byline": "^5.0.0",
                                      "webpack": "^4.42.0",
                                      "webpack-cli": "^3.3.11"
                                    },
                                    "readme": "[![styled with prettier](https://img.shields.io/badge/styled_with-prettier-ff69b4.svg)](https://github.com/prettier/prettier)\r\n[![NPM version](https://img.shields.io/npm/v/dbgate-plugin-sqlite.svg)](https://www.npmjs.com/package/dbgate-plugin-sqlite)\r\n\r\n# dbgate-plugin-sqlite\r\n\r\nUse DbGate for install of this plugin\r\n",
                                    "isPackaged": true
                                  },
                                  {
                                    "name": "dbgate-plugin-xml",
                                    "main": "dist/backend.js",
                                    "version": "1.0.0",
                                    "homepage": "https://dbgate.org",
                                    "description": "XML import/export plugin for DbGate",
                                    "repository": {
                                      "type": "git",
                                      "url": "https://github.com/dbgate/dbgate"
                                    },
                                    "author": "Jan Prochazka",
                                    "license": "MIT",
                                    "keywords": [
                                      "xml",
                                      "import",
                                      "export",
                                      "dbgate",
                                      "dbgateplugin"
                                    ],
                                    "files": [
                                      "dist"
                                    ],
                                    "scripts": {
                                      "build:frontend": "webpack --config webpack-frontend.config",
                                      "build:frontend:watch": "webpack --watch --config webpack-frontend.config",
                                      "build:backend": "webpack --config webpack-backend.config.js",
                                      "build": "yarn build:frontend && yarn build:backend",
                                      "plugin": "yarn build && yarn pack && dbgate-plugin dbgate-plugin-xml",
                                      "copydist": "yarn build && yarn pack && dbgate-copydist ../dist/dbgate-plugin-xml",
                                      "plugout": "dbgate-plugout dbgate-plugin-xml",
                                      "prepublishOnly": "yarn build"
                                    },
                                    "devDependencies": {
                                      "node-xml-stream-parser": "^1.0.12",
                                      "dbgate-plugin-tools": "^1.0.7",
                                      "webpack": "^4.42.0",
                                      "webpack-cli": "^3.3.11"
                                    },
                                    "readme": "[![styled with prettier](https://img.shields.io/badge/styled_with-prettier-ff69b4.svg)](https://github.com/prettier/prettier)\r\n[![NPM version](https://img.shields.io/npm/v/dbgate-plugin-xml.svg)](https://www.npmjs.com/package/dbgate-plugin-xml)\r\n\r\n# dbgate-plugin-xml\r\n\r\nUse DbGate for install of this plugin\r\n",
                                    "isPackaged": true
                                  }
                                ]
) {
  const extensions = {
    plugins,
    drivers: buildDrivers(plugins),
  };
  return extensions
}

function buildDrivers(plugins) {
  const res = [];
  for (const { content } of plugins) {
    console.log(plugins)


    if (content.drivers) { // @ts-ignore
      res.push(...content.drivers);
    }
  }
  return res;
}