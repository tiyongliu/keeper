import {isNaN, isNumber} from 'lodash-es'
import {getCurrentSettings} from '/@/store'

export function getBoolSettingsValue(name, defaultValue) {
  const settings = getCurrentSettings();
  const res = settings[name];
  if (res == null) return defaultValue;
  return !!res;
}

export function getIntSettingsValue(name, defaultValue, min: null | number = null, max: null | number = null) {
  const settings = getCurrentSettings();
  const parsed = parseInt(settings[name]);
  if (isNaN(parsed)) {
    return defaultValue;
  }
  if (isNumber(parsed)) {
    if (min != null && parsed < min) return min;
    if (max != null && parsed > max) return max;
    return parsed;
  }
  return defaultValue;
}
