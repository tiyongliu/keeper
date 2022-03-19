/**
 * Independent time operation tool to facilitate subsequent switch to dayjs
 */
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime'
const DATE_TIME_FORMAT = 'YYYY-MM-DD HH:mm:ss';
const DATE_FORMAT = 'YYYY-MM-DD';

// 使用中文语言包。固定格式
// dayjs.locale('zh-cn')
// 使用对时间插件插件。固定格式dayjs.extend(插件)
dayjs.extend(relativeTime)

export function formatToDateTime(
  date: dayjs.Dayjs | undefined = undefined,
  format = DATE_TIME_FORMAT,
): string {
  return dayjs(date).format(format);
}

export function formatToDate(
  date: dayjs.Dayjs | undefined = undefined,
  format = DATE_FORMAT,
): string {
  return dayjs(date).format(format);
}

export function fromNow(
  date: dayjs.Dayjs | undefined = undefined,
  withoutSuffix?: boolean): string {
  return dayjs(date).fromNow(withoutSuffix)
}

export const dateUtil = dayjs;
