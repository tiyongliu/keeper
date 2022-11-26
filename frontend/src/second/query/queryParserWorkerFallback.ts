import {splitQuery} from '/@/second/keeper-splitter'

export default function c(data) {
  return splitQuery(data.text, data.options);
}
