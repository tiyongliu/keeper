<template>
  <JSONNested
    :name="name"
    :expanded="expanded"
    :isParentExpanded="isParentExpanded"
    :isParentArray="isParentArray"
    :isArray="true"
    bracketOpen="["
    bracketClose="]"
    :keys="keys"
    :previewKeys="previewKeys"
    :getValue="getValue"
    :label="`Array(${value?.length})`"
  />
</template>

<script lang="ts" setup>
import JSONNested from './JSONNested.vue'
import {computed, defineProps, PropType, ref, toRefs} from 'vue'

const props = defineProps({
  name: {
    type: String as PropType<string>,
  },
  value: {
    type: [String, Number, Boolean, Object, Array] as PropType<string | boolean | number | object | string[]>,
  },
  isParentExpanded: {
    type: Boolean as PropType<boolean>,
  },
  isParentArray: {
    type: Boolean as PropType<boolean>,
  },
  expanded: {
    type: Boolean as PropType<boolean>,
    default: false
  },
})

const {
  name,
  value,
  isParentExpanded,
  isParentArray,
  expanded,
} = toRefs(props)
const filteredKey = ref(new Set(['length']))
const keys = computed(() => Object.getOwnPropertyNames(value?.value))
const previewKeys = computed(() => keys?.value ? keys?.value.filter(name => !filteredKey.value.has(name)) : [])

function getValue(key) {
  return value?.value![key]
}
</script>
