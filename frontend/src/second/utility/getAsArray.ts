import {isArray} from 'lodash-es'

export default function getAsArray(obj) {
  if (isArray(obj)) return obj;
  if (obj != null) return [obj];
  return [];
}
