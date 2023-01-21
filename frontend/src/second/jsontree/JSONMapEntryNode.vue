<template>
  <JSONNested
    :expanded="expanded"
    :isParentExpanded="isParentExpanded"
    :isParentArray="isParentArray"
    :name="isParentExpanded ? String(name) : (value.name || value.key)"
    :keys="keys"
    :getValue="getValue"
    :label="isParentExpanded ? 'Entry ' : '=> '"
    bracketOpen="{"
    bracketClose="}"
  />
</template>

<script lang="ts" setup>
import {defineProps, PropType, ref, toRefs} from 'vue'
import JSONNested from './JSONNested.vue'

const props = defineProps({
  name: {
    type: String as PropType<string>,
  },
  value: {
    type: Object as PropType<{key: string; name?: string}>,
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
  }
})

const {
  name,
  value,
  isParentExpanded,
  isParentArray,
  expanded,
} = toRefs(props)
const keys = ref(['key', 'value'])

function getValue(key) {
  return value?.value![key]
}
</script>
