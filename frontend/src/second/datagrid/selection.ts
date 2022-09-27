import {isNumber} from 'lodash-es'
export type CellAddress = [number | 'header' | 'filter' | undefined, number | 'header' | undefined]
export type RegularCellAddress = [number, number];

export const topLeftCell: CellAddress = [0, 0];
export const undefinedCell: CellAddress = [undefined, undefined];
export const nullCell: CellAddress = null;
export const emptyCellArray: CellAddress[] = [];

export function isRegularCell(cell: CellAddress): cell is RegularCellAddress {
  if (!cell) return false;
  const [row, col] = cell;
  return isNumber(row) && isNumber(col);
}

function normalizeHeaderForSelection(addr: CellAddress): CellAddress {
  if (addr[0] == 'filter') return ['header', addr[1]];
  return addr;
}

export function convertCellAddress(row, col): CellAddress {
  const rowNumber = parseInt(row);
  const colNumber = parseInt(col);
  return [isNaN(rowNumber) ? row : rowNumber, isNaN(colNumber) ? col : colNumber];
}

export function cellFromEvent(event): CellAddress {
  const cell = event.target.closest('td');
  if (!cell) return undefinedCell;
  const col = cell.getAttribute('data-col');
  const row = cell.getAttribute('data-row');
  return convertCellAddress(row, col);
}

