import {SqlDialect} from './dialect'
import {SqlDumper} from './dumper'
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


  dialect: SqlDialect;
  dialectByVersion(version): SqlDialect;
  createDumper(options = null): SqlDumper;


  analyserClass?: any;
  dumperClass?: any;
}
