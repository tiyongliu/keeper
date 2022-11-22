import {omit, pick} from 'lodash-es'
import {createGridCache, GridCache, GridConfig} from './GridConfig'
import {ChangeCacheFunc, ChangeConfigFunc, DisplayColumn} from './GridDisplay'
import {DatabaseInfo, EngineDriver, SqlDialect, TableInfo} from '/@/second/keeper-types'
import {getFilterValueExpression} from '/@/second/keeper-filterparser';


export class FormViewDisplay {
  isLoadedCorrectly = true;
  columns: DisplayColumn[];
  public baseTable: TableInfo;
  dialect: SqlDialect | undefined;

  constructor(
    public config: GridConfig,
    protected setConfig: ChangeConfigFunc,
    public cache: GridCache,
    protected setCache: ChangeCacheFunc,
    public driver?: EngineDriver,
    public dbinfo: Nullable<DatabaseInfo> = null,
    public serverVersion = null
  ) {
    this.dialect = (driver?.dialectByVersion && driver?.dialectByVersion(serverVersion)) || driver?.dialect;
  }

  addFilterColumn(column) {
    if (!column) return;
    this.setConfig(cfg => ({
      ...cfg,
      formFilterColumns: [...(cfg.formFilterColumns || []), column.uniqueName],
    }));
  }

  filterCellValue(column, rowData) {
    if (!column || !rowData) return;
    const value = rowData[column.uniqueName];
    const expr = getFilterValueExpression(value, column.dataType);
    if (expr) {
      this.setConfig(cfg => ({
        ...cfg,
        filters: {
          ...cfg.filters,
          [column.uniqueName]: expr,
        },
        addedColumns: cfg.addedColumns.includes(column.uniqueName)
          ? cfg.addedColumns
          : [...cfg.addedColumns, column.uniqueName],
      }));
      this.reload();
    }
  }

  setFilter(uniqueName, value) {
    this.setConfig(cfg => ({
      ...cfg,
      filters: {
        ...cfg.filters,
        [uniqueName]: value,
      },
    }));
    this.reload();
  }

  removeFilter(uniqueName) {
    const reloadRequired = !!this.config.filters[uniqueName];
    this.setConfig(cfg => ({
      ...cfg,
      formFilterColumns: (cfg.formFilterColumns || []).filter(x => x != uniqueName),
      filters: omit(cfg.filters || [], uniqueName),
    }));
    if (reloadRequired) this.reload();
  }

  reload() {
    this.setCache(_ => ({
      // ...cache,
      ...createGridCache(),
      refreshTime: new Date().getTime(),
    }));
  }

  getKeyValue(columnName) {
    const {formViewKey, formViewKeyRequested} = this.config;
    if (formViewKeyRequested && formViewKeyRequested[columnName]) return formViewKeyRequested[columnName];
    if (formViewKey && formViewKey[columnName]) return formViewKey[columnName];
    return null;
  }

  requestKeyValue(columnName, value) {
    if (this.getKeyValue(columnName) == value) return;

    this.setConfig(cfg => ({
      ...cfg,
      formViewKeyRequested: {
        ...cfg.formViewKey,
        ...cfg.formViewKeyRequested,
        [columnName]: value,
      },
    }));
    this.reload();
  }

  extractKey(row) {
    if (!row || !this.baseTable || !this.baseTable.primaryKey) {
      return null;
    }
    return pick(
      row,
      this.baseTable.primaryKey.columns.map(x => x.columnName)
    );
  }

  cancelRequestKey(rowData) {
    this.setConfig(cfg => ({
      ...cfg,
      formViewKeyRequested: null,
      formViewKey: rowData ? this.extractKey(rowData) : cfg.formViewKey,
    } as any));
  }
}
