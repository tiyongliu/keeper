import {MacroDefinition, MacroSelectedCell} from './MacroDefinition'
import {ChangeSet, setChangeSetValue} from './ChangeSet';
import {GridDisplay} from './GridDisplay';

const getMacroFunction = {
  transformValue: code => `
(value, args, modules, rowIndex, row, columnName) => {
    ${code}
}
`,
  transformRows: code => `
(rows, args, modules, selectedCells, cols, columns) => {
  ${code}
}
`,
  transformData: code => `
(rows, args, modules, selectedCells, cols, columns) => {
  ${code}
}
`,
};

export function compileMacroFunction(macro: MacroDefinition, errors = []) {
  if (!macro) return null;
  let func;
  try {
    func = eval(getMacroFunction[macro.type](macro.code));
    return func;
  } catch (err) {
    // @ts-ignore
    errors.push(`Error compiling macro ${macro.name}: ${err.message}`);
    return null;
  }
}

const modules = {};

export function runMacroOnValue(compiledFunc, macroArgs, value, rowIndex, row, column, errors = []) {
  if (!compiledFunc) return value;
  try {
    const res = compiledFunc(value, macroArgs, modules, rowIndex, row, column);
    return res;
  } catch (err) {
    // @ts-ignore
    errors.push(`Error processing column ${column} on row ${rowIndex}: ${err.message}`);
    return value;
  }
}

export function runMacroOnChangeSet(
  macro: MacroDefinition,
  macroArgs: {},
  selectedCells: MacroSelectedCell[],
  changeSet: ChangeSet,
  display: GridDisplay
): Nullable<ChangeSet> {
  const errors = [];
  const compiledMacroFunc = compileMacroFunction(macro, errors);
  if (!compiledMacroFunc) return null;

  let res = changeSet;
  for (const cell of selectedCells) {
    const definition = display.getChangeSetField(cell.rowData, cell.column, undefined);
    const macroResult = runMacroOnValue(
      compiledMacroFunc,
      macroArgs,
      cell.value,
      cell.row,
      cell.rowData,
      cell.column,
      errors
    );
    res = setChangeSetValue(res, definition!, macroResult);
  }

  return res;
}
