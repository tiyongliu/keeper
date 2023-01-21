<template>
    <li :class="isParentExpanded && 'indent'">
    <JSONKey
      :name="name"
      :colon="context['colon']"
      :isParentExpanded="isParentExpanded"
      :isParentArray="isParentArray"/>
    <span :class="nodeType">{{ showText }}</span>
  </li>
</template>

<script lang="ts" setup>
import {computed, defineProps, inject, PropType, toRefs} from 'vue'
import JSONKey from './JSONKey.vue'

const props = defineProps({
  name: {
    type: String as PropType<string>,
  },
  value: {
    type: [String, Number, Boolean, Object, Array] as PropType<string | boolean | number | object | string[]>,
  },
  valueGetter: {
    type: Function as PropType<(p: any) => string>,
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
  valueGetter,
  isParentExpanded,
  isParentArray,
  nodeType,
} = toRefs(props)

const context = inject('json-tree-context-key') as { [key in string]: any }
const showText = computed(() => valueGetter?.value ? valueGetter.value(value?.value) : value?.value)
</script>

<style scoped>
li {
  user-select: text;
  word-wrap: break-word;
  word-break: break-all;
}

.indent {
  padding-left: var(--li-identation);
}

.String {
  color: var(--string-color);
}

.Date {
  color: var(--date-color);
}

.Number {
  color: var(--number-color);
}

.Boolean {
  color: var(--boolean-color);
}

.Null {
  color: var(--null-color);
}

.Undefined {
  color: var(--undefined-color);
}

.Function {
  color: var(--function-color);
  font-style: italic;
}

.Symbol {
  color: var(--symbol-color);
}
</style>
