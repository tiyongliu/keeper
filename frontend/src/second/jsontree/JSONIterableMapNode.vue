<template>
  <JSONNested
    :name="name"
    :isParentExpanded="isParentExpanded"
    :isParentArray="isParentArray"
    :keys="keys"
    :getKey="getKey"
    :getValue="getValue"
    :label="nodeType`(${keys.length})`"
    colon=""
    bracketOpen="{"
    bracketClose="}"
  />
</template>

<script lang="ts" setup>
import {defineProps, PropType, ref, toRefs, watchEffect} from "vue";
import JSONNested from './JSONNested.vue'
import MapEntry from './utils/MapEntry'

const props = defineProps({
  name: {
    type: String as PropType<string>,
  },
  value: {
    type: [String, Object, Array] as PropType<string | object | string[]>,
  },
  isParentExpanded: {
    type: Boolean as PropType<boolean>,
  },
  isParentArray: {
    type: Boolean as PropType<boolean>,
  },
  nodeType: {
    type: String as PropType<string>,
  },
})

const {
  name,
  value,
  isParentExpanded,
  isParentArray,
  nodeType,
} = toRefs(props)
const keys = ref<any[]>([])

watchEffect(() => {
  try {
    let result: any[] = []
    let i = 0;
    for (const entry of value?.value) {
      result.push([i++, new MapEntry(entry[0], entry[1])]);
    }
    keys.value = result;
  } catch (e) {
    console.log(e)
  }
})

function getKey(entry) {
  return entry[0];
}

function getValue(entry) {
  return entry[1];
}
</script>
