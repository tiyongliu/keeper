import {EngineDriver, SqlDumper} from '/@/second/keeper-types'
import {Condition} from './types';

export function treeToSql<T>(driver: EngineDriver, object: T, func: (dmp: SqlDumper, obj: T) => void) {
  const dmp = driver.createDumper();
  func(dmp, object);
  return dmp.s;
}

export function mergeConditions(condition1: Condition, condition2: Condition): Condition {
  if (!condition1) return condition2;
  if (!condition2) return condition1;
  if (condition1.conditionType == 'and' && condition2.conditionType == 'and') {
    return {
      conditionType: 'and',
      conditions: [...condition1.conditions, ...condition2.conditions],
    };
  }
  if (condition1.conditionType == 'and') {
    return {
      conditionType: 'and',
      conditions: [...condition1.conditions, condition2],
    };
  }
  if (condition2.conditionType == 'and') {
    return {
      conditionType: 'and',
      conditions: [condition1, ...condition2.conditions],
    };
  }
  return {
    conditionType: 'and',
    conditions: [condition1, condition2],
  };
}
