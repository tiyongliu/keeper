import { isTypeNumber, isTypeString, isTypeLogical, isTypeDateTime } from '/@/second/keeper-tools';
import { FilterType } from './types';

export function getFilterType(dataType: string): FilterType {
  if (!dataType) return 'string';
  if (isTypeNumber(dataType)) return 'number';
  if (isTypeString(dataType)) return 'string';
  if (isTypeLogical(dataType)) return 'logical';
  if (isTypeDateTime(dataType)) return 'datetime';
  return 'string';
}
