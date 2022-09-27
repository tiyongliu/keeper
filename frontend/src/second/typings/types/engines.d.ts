export interface EngineDriver {
  engine: string;
  title: string;
  defaultPort?: number;
  databaseEngineTypes: string[];
  editorMode?: string;
  readOnlySessions: boolean;
  supportedKeyTypes: SupportedDbKeyType[];
  supportsDatabaseUrl?: boolean;
  supportsDatabaseDump?: boolean;
  isElectronOnly?: boolean;
  supportedCreateDatabase?: boolean;
  showConnectionField?: (field: string, values: any) => boolean;
  showConnectionTab?: (tab: 'ssl' | 'sshTunnel', values: any) => boolean;
  beforeConnectionSave?: (values: any) => any;
  databaseUrlPlaceholder?: string;
  defaultAuthTypeName?: string;
  defaultSocketPath?: string;
  authTypeLabel?: string;
  importExportArgs?: any[];
  connect({ server, port, user, password, database }): Promise<any>;
  close(pool): Promise<any>;
  query(pool: any, sql: string, options?: QueryOptions): Promise<QueryResult>;
  stream(pool: any, sql: string, options: StreamOptions);
  readQuery(pool: any, sql: string, structure?: TableInfo): Promise<stream.Readable>;
  readJsonQuery(pool: any, query: any, structure?: TableInfo): Promise<stream.Readable>;
  writeTable(pool: any, name: NamedObjectInfo, options: WriteTableOptions): Promise<stream.Writeable>;
  analyseSingleObject(
    pool: any,
    name: NamedObjectInfo,
    objectTypeField: keyof DatabaseInfo
  ): Promise<TableInfo | ViewInfo | ProcedureInfo | FunctionInfo | TriggerInfo>;
  analyseSingleTable(pool: any, name: NamedObjectInfo): Promise<TableInfo>;
  getVersion(pool: any): Promise<{ version: string }>;
  listDatabases(pool: any): Promise<
    {
      name: string;
    }[]
    >;
  loadKeys(pool, root: string, filter?: string): Promise;
  exportKeys(pool, options: {}): Promise;
  loadKeyInfo(pool, key): Promise;
  loadKeyTableRange(pool, key, cursor, count): Promise;
  loadFieldValues(pool: any, name: NamedObjectInfo, field: string, search: string): Promise;
  analyseFull(pool: any, serverVersion): Promise<DatabaseInfo>;
  analyseIncremental(pool: any, structure: DatabaseInfo, serverVersion): Promise<DatabaseInfo>;
  dialect: SqlDialect;
  dialectByVersion(version): SqlDialect;
  createDumper(options = null): SqlDumper;
  createBackupDumper(pool: any, options): Promise<SqlBackupDumper>;
  getAuthTypes(): EngineAuthType[];
  readCollection(pool: any, options: ReadCollectionOptions): Promise<any>;
  updateCollection(pool: any, changeSet: any): Promise<any>;
  getCollectionUpdateScript(changeSet: any): string;
  createDatabase(pool: any, name: string): Promise;
  dropDatabase(pool: any, name: string): Promise;
  getQuerySplitterOptions(usage: 'stream' | 'script' | 'editor'): any;
  script(pool: any, sql: string): Promise;
  getNewObjectTemplates(): NewObjectTemplate[];
  // direct call of pool method, only some methods could be supported, on only some drivers
  callMethod(pool, method, args);

  analyserClass?: any;
  dumperClass?: any;
}
