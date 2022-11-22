import {FormViewDisplay} from './FormViewDisplay'
import {GridCache, GridConfig} from './GridConfig'
import {ChangeCacheFunc, ChangeConfigFunc, DisplayColumn} from './GridDisplay'
import {DictionaryDescriptionFunc, TableGridDisplay} from './TableGridDisplay'
import {Condition, mergeConditions, OrderByExpression} from '/@/second/keeper-sqltree'
import {DatabaseInfo, EngineDriver, NamedObjectInfo} from '/@/second/keeper-types'
import stableStringify from 'json-stable-stringify'
import {ChangeSetFieldDefinition, ChangeSetRowDefinition} from './ChangeSet';

export class TableFormViewDisplay extends FormViewDisplay {
  // use utility functions from GridDisplay and publish result in FromViewDisplay interface
  private gridDisplay: TableGridDisplay;

  constructor(
    public tableName: NamedObjectInfo,
    driver: EngineDriver,
    config: GridConfig,
    setConfig: ChangeConfigFunc,
    cache: GridCache,
    setCache: ChangeCacheFunc,
    dbinfo: DatabaseInfo,
    displayOptions,
    serverVersion,
    getDictionaryDescription: Nullable<DictionaryDescriptionFunc> = null,
    isReadOnly = false
  ) {
    super(config, setConfig, cache, setCache, driver, dbinfo, serverVersion);
    this.gridDisplay = new TableGridDisplay(
      tableName,
      driver,
      config,
      setConfig,
      cache,
      setCache,
      dbinfo,
      displayOptions,
      serverVersion,
      getDictionaryDescription,
      isReadOnly
    );
    this.gridDisplay.addAllExpandedColumnsToSelected = true;

    this.isLoadedCorrectly = this.gridDisplay.isLoadedCorrectly && !!this.driver;
    this.columns = [];
    this.addDisplayColumns(this.gridDisplay.columns);
    this.baseTable = this.gridDisplay.baseTable!;
    this.gridDisplay.hintBaseColumns = this.columns;
  }

  addDisplayColumns(columns: DisplayColumn[]) {
    for (const col of columns) {
      this.columns.push(col);
      if (this.gridDisplay.isExpandedColumn(col.uniqueName)) {
        const table = this.gridDisplay.getFkTarget(col);
        if (table) {
          const subcolumns = this.gridDisplay.getDisplayColumns(table, col.uniquePath) as DisplayColumn[];
          this.addDisplayColumns(subcolumns);
        }
      }
    }
  }

  getPrimaryKeyEqualCondition(row = null): Nullable<Condition> {
    // @ts-ignore
    if (!row) row = this.config.formViewKeyRequested || this.config.formViewKey;
    if (!row) return null;
    // @ts-ignore
    const {primaryKey} = this.gridDisplay.baseTable;
    if (!primaryKey) return null;
    return {
      conditionType: 'and',
      conditions: primaryKey.columns.map(({columnName}) => ({
        conditionType: 'binary',
        operator: '=',
        left: {
          exprType: 'column',
          columnName,
          source: {
            alias: 'basetbl',
          },
        },
        right: {
          exprType: 'value',
          value: row![columnName],
        },
      })),
    };
  }

  getPrimaryKeyOperatorCondition(operator): Nullable<Condition> {
    if (!this.config.formViewKey) return null;
    const conditions: any[] = [];
    // @ts-ignore
    const {primaryKey} = this.gridDisplay.baseTable;
    if (!primaryKey) return null;
    for (let index = 0; index < primaryKey.columns.length; index++) {

      conditions.push({
        conditionType: 'and',
        conditions: [
          ...primaryKey.columns.slice(0, index).map(({columnName}) => ({
            conditionType: 'binary',
            operator: '=',
            left: {
              exprType: 'column',
              columnName,
              source: {
                alias: 'basetbl',
              },
            },
            right: {
              exprType: 'value',
              value: this.config.formViewKey![columnName],
            },
          })),
          ...primaryKey.columns.slice(index).map(({columnName}) => ({
            conditionType: 'binary',
            operator: operator,
            left: {
              exprType: 'column',
              columnName,
              source: {
                alias: 'basetbl',
              },
            },
            right: {
              exprType: 'value',
              value: this.config.formViewKey![columnName],
            },
          })),
        ],
      });
    }

    if (conditions.length == 1) {
      return conditions[0];
    }

    return {
      conditionType: 'or',
      conditions,
    };
  }

  getSelect() {
    if (!this.driver) return null;
    const select = this.gridDisplay.createSelect();
    if (!select) return null;
    select.topRecords = 1;
    return select;
  }

  getCurrentRowQuery() {
    const select = this.getSelect()!;
    if (!select) return null;

    select.where = mergeConditions(select.where!, this.getPrimaryKeyEqualCondition()!);
    return select;
  }

  getCountSelect() {
    const select = this.getSelect();
    if (!select) return null;
    // @ts-ignore
    select.orderBy = null;
    select.columns = [
      {
        exprType: 'raw',
        sql: 'COUNT(*)',
        alias: 'count',
      },
    ];
    // @ts-ignore
    select.topRecords = null;
    return select;
  }

  getCountQuery() {
    if (!this.driver) return null;
    const select = this.getCountSelect();
    if (!select) return null;
    return select;
  }

  getBeforeCountQuery() {
    if (!this.driver) return null;
    const select = this.getCountSelect();
    if (!select) return null;
    select.where = mergeConditions(select.where!, this.getPrimaryKeyOperatorCondition('<')!)
    return select;
  }

  navigate(row) {
    const formViewKey = this.extractKey(row);
    this.setConfig(cfg => ({
      ...cfg,
      formViewKey,
    } as any));
  }

  isLoadedCurrentRow(row) {
    if (!row) return false;
    const formViewKey = this.extractKey(row);
    return stableStringify(formViewKey) == stableStringify(this.config.formViewKey);
  }

  navigateRowQuery(commmand: 'begin' | 'previous' | 'next' | 'end') {
    if (!this.driver) return null;
    const select = this.gridDisplay.createSelect();
    if (!select) return null;
    // @ts-ignore
    const {primaryKey} = this.gridDisplay.baseTable;

    function getOrderBy(direction): OrderByExpression[] {
      return primaryKey.columns.map(({columnName}) => ({
        exprType: 'column',
        columnName,
        direction,
      }));
    }

    select.topRecords = 1;
    switch (commmand) {
      case 'begin':
        select.orderBy = getOrderBy('ASC');
        break;
      case 'end':
        select.orderBy = getOrderBy('DESC');
        break;
      case 'previous':
        select.orderBy = getOrderBy('DESC');
        select.where = mergeConditions(select.where!, this.getPrimaryKeyOperatorCondition('<')!);
        break;
      case 'next':
        select.orderBy = getOrderBy('ASC');
        select.where = mergeConditions(select.where!, this.getPrimaryKeyOperatorCondition('>')!);
        break;
    }

    return select;
  }

  getChangeSetRow(row): Nullable<ChangeSetRowDefinition> {
    if (!this.baseTable) return null;
    return <ChangeSetRowDefinition>{
      pureName: this.baseTable.pureName,
      schemaName: this.baseTable.schemaName,
      condition: this.extractKey(row),
    };
  }

  getChangeSetField(row, uniqueName): Nullable<ChangeSetFieldDefinition> {
    const col = this.columns.find(x => x.uniqueName == uniqueName);
    if (!col) return null;
    if (!this.baseTable) return null;
    if (this.baseTable.pureName != col.pureName || this.baseTable.schemaName != col.schemaName) return null;
    return <ChangeSetFieldDefinition>{
      ...this.getChangeSetRow(row),
      uniqueName: uniqueName,
      columnName: col.columnName,
    };
  }

  toggleExpandedColumn(uniqueName: string, value?: boolean) {
    this.gridDisplay.toggleExpandedColumn(uniqueName, value);
    this.gridDisplay.reload();
  }

  isExpandedColumn(uniqueName: string) {
    return this.gridDisplay.isExpandedColumn(uniqueName);
  }

  get editable() {
    return this.gridDisplay.editable;
  }
}
