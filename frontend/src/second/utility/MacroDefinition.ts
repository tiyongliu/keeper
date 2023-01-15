export interface MacroArgument {
  type: 'text' | 'select';
  label: string;
  name: string;
}

export interface MacroDefinition {
  title: string;
  name: string;
  group: string;
  description?: string;
  type: 'transformValue';
  code: string;
  args?: MacroArgument[];
}

export interface MacroSelectedCell {
  column: string;
  row: number;
  rowData: any;
  value: any;
}
