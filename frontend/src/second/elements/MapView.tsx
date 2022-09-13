import {find} from 'lodash-es'
import {isWktGeometry} from '/@/second/keeper-tools'
export function selectionCouldBeShownOnMap(selection: {column: string; value: any}[]) {
  if (selection.length > 0 && find(selection, x => isWktGeometry(x.value))) {
    return true;
  }

  if (
    selection.find(x => x.column.toLowerCase().includes('lat')) &&
    (selection.find(x => x.column.toLowerCase().includes('lon')) ||
      selection.find(x => x.column.toLowerCase().includes('lng')))
  ) {
    return true;
  }
  return false;
}
