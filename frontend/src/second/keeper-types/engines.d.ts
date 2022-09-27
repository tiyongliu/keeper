import {SqlDialect} from './dialect'
import {SqlDumper} from './dumper'
export interface EngineDriver {
  engine: string;
  title: string;

  dialect: SqlDialect;
  dialectByVersion(version): SqlDialect;
  createDumper(options = null): SqlDumper;
}
