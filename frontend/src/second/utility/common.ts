import {startCase} from 'lodash-es'
import {useLocaleStore} from '/@/store/modules/locale'

const locale = useLocaleStore()

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

export function resolveKeyText(keyText: string): string {
  if (isMac()) {
    return keyText.replace('CtrlOrCommand+', 'Command+');
  }
  return keyText.replace('CtrlOrCommand+', 'Ctrl+');
}

export function isCtrlOrCommandKey(event) {
  if (isMac()) {
    return event.metaKey;
  }
  return event.ctrlKey;
}

export function setSelectedTabFunc(files, tabid) {
  return [
    ...(files || []).filter(x => x.tabid != tabid).map(x => ({ ...x, selected: false })),
    ...(files || []).filter(x => x.tabid == tabid).map(x => ({ ...x, selected: true })),
  ];
}

export function setSelectedTab(tabid) {
  locale.updateOpenedTabs(tabs => setSelectedTabFunc(tabs, tabid))
}
