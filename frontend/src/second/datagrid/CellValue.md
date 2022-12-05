```tsx
import {defineComponent, PropType, toRefs} from 'vue'
import {isArray, isDate, isNumber, isPlainObject, isString, pad} from 'lodash-es'

const dateTimeRegex =
  /^([0-9]+)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])[Tt]([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]|60)(\.[0-9]+)?(([Zz])|()|([\+|\-]([01][0-9]|2[0-3]):[0-5][0-9]))$/;

function formatNumber(value) {
  if (value >= 10000 || value <= -10000) {
    // if (getBoolSettingsValue('dataGrid.thousandsSeparator', false)) {
    //   return value.toLocaleString();
    // } else {
    //   return value.toString();
    // }
  }

  return value.toString();
}

function formatDateTime(testedString) {
  const m = testedString.match(dateTimeRegex);
  return `${m[1]}-${m[2]}-${m[3]} ${m[4]}:${m[5]}:${m[6]}`;
}

function makeBulletString(value) {
  return pad('', value.length, '•');
}

function highlightSpecialCharacters(value) {
  value = value.replace(/\n/g, '↲');
  value = value.replace(/\r/g, '');
  value = value.replace(/^(\s+)/, makeBulletString);
  value = value.replace(/(\s+)$/, makeBulletString);
  value = value.replace(/(\s\s+)/g, makeBulletString);
  return value;
}

export default defineComponent({
  name: "CellValue",
  props: {
    rowData: {
      type: Object as PropType<object>
    },
    value: {
      type: [String, Number, Boolean, Object, Array] as PropType<string | boolean | number | object | string[]>,
    },
    jsonParsedValue: {
      type: [Object, Array] as PropType<object | string[]>,
    }
  },
  setup(props) {
    const {rowData, value, jsonParsedValue} = toRefs(props)

    const renderValue = function () {
      if (rowData.value == null)
        return <span class='td-cell null'>(No row)</span>
      else if (value.value === null)
        return <span class='td-cell null'>(NULL)</span>
      else if (value.value === undefined)
        return <span class='td-cell null'>(No field)</span>
      else if (isDate(value.value))
        return <>{value.toString()}</>
      else if (value.value === true)
        return <span class='td-cell value'>true</span>
      else if (value.value === false)
        return <span class='td-cell value'>false</span>
      else if (isNumber(value.value))
        return <span class='td-cell value'>{formatNumber(value.value)}</span>
      else if (isString(value.value) && !jsonParsedValue.value) {
        if (dateTimeRegex.test(value.value)) {
          return <span class='td-cell value'>{formatDateTime(value.value)}</span>
        } else {
          return highlightSpecialCharacters(value.value)
        }
      } else if (value.value['$oid'])
        return <span class='value'>ObjectId("{value.value['$oid']}")</span>
      else if (isPlainObject(value.value))
        return <span class='td-cell null'
                     title={`${JSON.stringify(value.value, undefined, 2)}`}>(JSON)</span>
      else if (isArray(value.value))
        return <span class='td-cell null'
                     title={`${value.value.map(x => JSON.stringify(x)).join('\n')}`}>[{value.value.length} items]</span>
      else if (isPlainObject(jsonParsedValue.value))
        return <span class='td-cell null'
                     title={`${JSON.stringify(jsonParsedValue.value, undefined, 2)}`}>(JSON)</span>
      else if (isArray(jsonParsedValue.value))
        return <span class='td-cell null'
                     title={`${jsonParsedValue.value.map(x => JSON.stringify(x)).join('\n')}`}>
          [{jsonParsedValue.value.length} items]</span>
    }

    return () => renderValue()
  }
})
```
