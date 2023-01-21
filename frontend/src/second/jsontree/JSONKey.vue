<template>
  <label v-if="showKey && name" :class="isParentExpanded && 'spaced'">
    <span>{{ name }}{{ colon }}</span>
  </label>
</template>

<script lang="ts" setup>
import {computed, defineProps, PropType, toRefs} from "vue"

const props = defineProps({
  name: {
    type: String as PropType<string>,
  },
  isParentExpanded: {
    type: Boolean as PropType<boolean>,
  },
  isParentArray: {
    type: Boolean as PropType<boolean>,
    default: false
  },
  colon: {
    type: String as PropType<string>,
    default: ''
  },
})

const {isParentArray, colon, isParentExpanded, name} = toRefs(props)

const showKey = computed(() => (isParentExpanded?.value || !isParentArray.value || name?.value != +name.value))
</script>

<style scoped>
label {
  display: inline-block;
  color: var(--label-color);
  padding: 0;
}

.spaced {
  padding-right: var(--li-colon-space);
}
</style>
