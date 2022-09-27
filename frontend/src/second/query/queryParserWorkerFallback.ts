import { splitQuery } from '/@/second/keeper-splitter'

export default function c(data) {
  const result = splitQuery(data.text, data.options);
  return result;
}
