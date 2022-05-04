import {startCase} from 'lodash-es';

export function getObjectTypeFieldLabel(objectTypeField) {
  if (objectTypeField == 'matviews') return 'Materialized Views';
  return startCase(objectTypeField)
}
