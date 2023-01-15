import dayjs from 'dayjs'
import {isTypeDateTime} from '/@/second/keeper-tools'

export type FilterMultipleValuesMode = 'is' | 'is_not' | 'contains' | 'begins' | 'ends';

export function getFilterValueExpression(value, dataType?) {
  if (value == null) return 'NULL'
  if (isTypeDateTime(dataType)) return dayjs(value).format('YYYY-MM-DD HH:mm:ss')
  if (value === true) return 'TRUE'
  if (value === false) return 'FALSE'
  if (value.$oid) return `ObjectId("${value.$oid}")`
  return `="${value}"`
}

export function createMultiLineFilter(mode: FilterMultipleValuesMode, text: string) {
  let res = '';
  for (let line of text.split('\n')) {
    line = line.trim();
    if (line.length == 0) continue;

    if (res.length > 0) {
      switch (mode) {
        case 'is_not':
          res += ' ';
          break;
        default:
          res += ',';
          break;
      }
    }

    switch (mode) {
      case 'is':
        res += "='" + line + "'";
        break;
      case 'is_not':
        res += "<>'" + line + "'";
        break;
      case 'contains':
        res += "'" + line + "'";
        break;
      case 'begins':
        res += "^'" + line + "'";
        break;
      case 'ends':
        res += "$'" + line + "'";
        break;
    }
  }
  return res;
}
