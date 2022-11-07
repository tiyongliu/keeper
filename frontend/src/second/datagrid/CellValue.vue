<template>
  <span class="null" v-if="rowData == null">(No row)</span>
  <span class="null" v-else-if="value === null">(NULL)</span>
  <span class="null" v-else-if="value === undefined">(No field)</span>
  <template v-else-if="isDate(value)">{{value.toString()}}</template>
  <span class="value" v-else-if="value === true">true</span>
  <span class="value" v-else-if="value === false">false</span>
  <span class="value" v-else-if="isNumber(value)">{{formatNumber(value)}}</span>
  <template v-else-if="isString(value) && !jsonParsedValue">
    <span class="value" v-if="dateTimeRegex.test(value)">
      {{formatDateTime(value)}}
    </span>
    <template v-else>
      {{highlightSpecialCharacters(value)}}
    </template>
  </template>
  <template v-else-if="value?.type == 'Buffer' && isArray(value.data)">
    <span class="value" v-if="value.data.length <= 16">{{'0x' + arrayToHexString(value.data)}}</span>
    <span class="null" v-else>({{value.data.length}}} bytes)</span>
  </template>
  <span class="value" v-else-if="value.$oid">ObjectId("{{value.$oid}}}")</span>
  <span class="null" v-else-if="isPlainObject(value)" :title="`${JSON.stringify(value, undefined, 2)}`">(JSON)</span>
  <span class="null" v-else-if="isArray(value)" :title="`${value.map(x => JSON.stringify(x)).join('\n')}`">[{{value.length}}} items]</span>
  <span class="null" v-else-if="isPlainObject(jsonParsedValue)" :title="`${JSON.stringify(jsonParsedValue, undefined, 2)}`">(JSON)</span>
  <span class="null" v-else-if="isArray(jsonParsedValue)" :title="`${jsonParsedValue.map(x => JSON.stringify(x)).join('\n')}`">[{jsonParsedValue.length} items]</span>
  <template v-else>{{value.toString() || ''}}</template>
</template>

<script lang="ts">
import {defineComponent, PropType, toRefs} from 'vue'
import {isDate, isNumber, isString, isArray, isPlainObject} from 'lodash-es'
import {arrayToHexString} from '/@/second/keeper-tools'
import {pad} from 'lodash-es'
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

// const dateTimeRegex = /^\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\d(\.\d\d\d)?Z?$/;
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

export default defineComponent({
  name: "CellValue",
  props: {
    rowData: {
      type: [Boolean, String, Number, Object, Array] as PropType<boolean | string | number | object | string[]>
    },
    value: {
      type: [Boolean, String, Number] as PropType<boolean | string | number>
    },
    jsonParsedValue: {
      type: [Boolean, String, Number, Object, Array] as PropType<boolean | string | number | object | string[]>,
    }
  },
  setup(props) {


    return {
      ...toRefs(props),
      isDate, isNumber, isString, isArray, isPlainObject,
      dateTimeRegex,
      formatDateTime,
      formatNumber,
      highlightSpecialCharacters,
      arrayToHexString,
    }
  }
})
</script>

<style scoped>
.null {
  color: var(--theme-font-3);
  font-style: italic;
}
.value {
  color: var(--theme-icon-green);
}
</style>
