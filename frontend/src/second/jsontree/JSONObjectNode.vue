<template>
  <JSONNested
    :name="name"
    :expanded="expanded"
    :isParentExpanded="isParentExpanded"
    :isParentArray="isParentArray"
    :keys="keys"
    :previewKeys="keys"
    :getValue="getValue"
    :label="labelOverride || `${nodeType}`"
    bracketOpen="{"
    bracketClose="}"
    :elementValue="value"
  />
</template>

<script lang="ts" setup>
import {computed, defineProps, PropType, toRefs, unref} from 'vue'
import JSONNested from './JSONNested.vue'

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
  nodeType: {
    type: String as PropType<string>,
  },
  expanded: {
    type: Boolean as PropType<boolean>,
    default: false
  },
  labelOverride: {
    type: String as PropType<String>,
    default: null
  },
})
const {
  name,
  value,
  isParentExpanded,
  isParentArray,
  nodeType,
  expanded,
  labelOverride,
} = toRefs(props)

const keys = computed(() => value?.value ? Object.getOwnPropertyNames(unref(value)) : [])

function getValue(key: string | number) {
  return value?.value![key]
}
</script>
