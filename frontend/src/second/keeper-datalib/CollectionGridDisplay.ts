import {compact, findLastIndex, isArray, isNumber, isPlainObject} from 'lodash-es'
import {ChangeCacheFunc, ChangeConfigFunc, GridDisplay} from './GridDisplay'
import type {CollectionInfo, EngineDriver} from '/@/second/keeper-types'
import {GridCache, GridConfig} from './GridConfig'


function getObjectKeys(obj) {
  if (isArray(obj)) {
    return Object.keys(obj)
      .slice(0, 10)
      .map(x => parseInt(x));
  }
  if (isPlainObject(obj)) {
    return Object.keys(obj);
  }
  return [];
}

function createHeaderText(path) {
  let res = `${path[0]}`;
  for (let i = 1; i < path.length; i++) {
    const name = path[i];
    if (isNumber(name)) res += `[${name}]`;
    else res += `.${name}`;
  }
  return res;
}

function getColumnsForObject(basePath, obj, res: any[], display) {
  for (const name of getObjectKeys(obj)) {
    const uniqueName = [...basePath, name].join('.');
    let column = res.find(x => x.uniqueName == uniqueName);
    if (!column) {
      column = getDisplayColumn(basePath, name, display);
      if (basePath.length > 0) {
        const lastIndex1 = findLastIndex(res, x => x.parentHeaderText.startsWith(column.parentHeaderText));
        const lastIndex2 = findLastIndex(res, x => x.headerText == column.parentHeaderText);
        // console.log(uniqueName, lastIndex1, lastIndex2);
        if (lastIndex1 >= 0) res.splice(lastIndex1 + 1, 0, column);
        else if (lastIndex2 >= 0) res.splice(lastIndex2 + 1, 0, column);
        else res.push(column);
      } else {
        res.push(column);
      }
    }
    if (isPlainObject(obj[name]) || isArray(obj[name])) {
      column.isExpandable = true;
    }

    if (display.isExpandedColumn(column.uniqueName)) {
      getColumnsForObject([...basePath, name], obj[name], res, display);
    }
  }
}

function getDisplayColumn(basePath, columnName, display) {
  const uniquePath = [...basePath, columnName];
  const uniqueName = uniquePath.join('.');
  return {
    columnName,
    headerText: createHeaderText(uniquePath),
    uniqueName,
    uniquePath,
    isStructured: true,
    parentHeaderText: createHeaderText(basePath),
    filterType: 'mongo',
    pureName: display.collection?.pureName,
    schemaName: display.collection?.schemaName,
  };
}

export function analyseCollectionDisplayColumns(rows, display) {
  const res: any[] = [];
  const addedColumns = display?.config?.addedColumns;
  for (const row of rows || []) {
    getColumnsForObject([], row, res, display);
  }
  for (const added of addedColumns || []) {
    // @ts-ignore
    if (res.find(x => x.uniqueName == added)) continue;
    // @ts-ignore
    res.push(getDisplayColumn([], added, display));
  }
  return (
    res.map(col => ({
      ...col,
      isChecked: display.isColumnChecked(col),
    })) || []
  );
}

export class CollectionGridDisplay extends GridDisplay {
  constructor(
    public collection: CollectionInfo,
    driver: EngineDriver,
    config: GridConfig,
    setConfig: ChangeConfigFunc,
    cache: GridCache,
    setCache: ChangeCacheFunc,
    loadedRows,
    changeSet,
    readOnly = false
  ) {
    super(config, setConfig, cache, setCache, driver);
    const changedDocs = compact(changeSet.updates.map(chs => chs.document));
    const insertedDocs = compact(changeSet.inserts.map(chs => chs.fields));
    this.columns = analyseCollectionDisplayColumns([...(loadedRows || []), ...changedDocs, ...insertedDocs], this);
    this.filterable = true;
    this.sortable = true;
    this.editable = !readOnly;
    this.supportsReload = true;
    this.isDynamicStructure = true;
    this.changeSetKeyFields = ['_id'];
    this.baseCollection = collection;
  }
}
