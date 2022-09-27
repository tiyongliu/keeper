import dayjs from 'dayjs'
import {isTypeDateTime} from '/@/second/keeper-tools'

export function getFilterValueExpression(value, dataType?) {
  if (value == null) return 'NULL'
  if (isTypeDateTime(dataType)) return dayjs(value).format('YYYY-MM-DD HH:mm:ss')
  if (value === true) return 'TRUE'
  if (value === false) return 'FALSE'
  if (value.$oid) return `ObjectId("${value.$oid}")`
  return `="${value}"`
}
