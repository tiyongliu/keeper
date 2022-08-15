import {startCase} from 'lodash-es';

export function getObjectTypeFieldLabel(objectTypeField) {
  if (objectTypeField == 'matviews') return 'Materialized Views';
  return startCase(objectTypeField)
}

export function isMac() {
  // @ts-ignore
  const platform = navigator?.platform || navigator?.userAgentData?.platform || 'unknown';
  return platform.toUpperCase().indexOf('MAC') >= 0;
}

export function formatKeyText(keyText: string): string {
  if (isMac()) {
    return keyText
      .replace('CtrlOrCommand+', '⌘ ')
      .replace('Shift+', '⇧ ')
      .replace('Alt+', '⌥ ')
      .replace('Command+', '⌘ ')
      .replace('Ctrl+', '⌃ ')
      .replace('Backspace', '⌫ ');
  }
  return keyText.replace('CtrlOrCommand+', 'Ctrl+');
}
