<template>
  <ul :class="[isDeleted && 'isDeleted', isInserted && 'isInserted', isModified && 'isModified']">
    <JSONNode
      :name="name"
      :value="value"
      :isParentExpanded="true"
      :isParentArray="false"
      :expanded="expandedReplica"
      :labelOverride="labelOverride"/>
  </ul>
</template>

<script lang="ts" setup>
import {defineProps, PropType, provide, ref, toRefs} from 'vue'
import {isUndefined} from 'lodash-es'
import JSONNode from './JSONNode.vue'

const props = defineProps({
  name: {
    type: String as PropType<string>,
    default: ''
  },
  //todo menu
  value: {
    type: [String, Number, Boolean, Object, Array] as PropType<string | boolean | number | object | string[]>,
  },
  expandAll: {
    type: Boolean as PropType<boolean>,
    default: false
  },
  expanded: {
    type: Boolean as PropType<boolean>,
  },
  labelOverride: {
    type: String as PropType<String>,
    default: null
  },
  slicedKeyCount: {
    type: Number as PropType<number>,
  },
  //todo disableContextMenu
  isDeleted: {
    type: Boolean as PropType<boolean>,
    default: false
  },
  isInserted: {
    type: Boolean as PropType<boolean>,
    default: false
  },
  isModified: {
    type: Boolean as PropType<boolean>,
    default: false
  },
})

const {
  name,
  value,
  expandAll,
  expanded,
  labelOverride,
  slicedKeyCount,
  isDeleted,
  isInserted,
  isModified,
} = toRefs(props)

provide('json-tree-context-key', {})
provide('json-tree-default-expanded', expandAll)
if (slicedKeyCount) provide('json-tree-sliced-key-count', slicedKeyCount);
const elementData = new WeakMap()
if (elementData) {
  provide('json-tree-element-data', elementData)
}
const expandedReplica = ref(expanded?.value)

if (isUndefined(expandedReplica?.value)) expandedReplica.value = expandAll.value
</script>

<style scoped>
:global(.theme-type-dark) ul {
  --json-tree-string-color: #ffc5c5;
  --json-tree-symbol-color: #ffc5c5;
  --json-tree-boolean-color: #b6c3ff;
  --json-tree-function-color: #b6c3ff;
  --json-tree-number-color: #bfbdff;
  --json-tree-label-color: #e9aaed;
  --json-tree-arrow-color: #d4d4d4;
  --json-tree-null-color: #dcdcdc;
  --json-tree-undefined-color: #dcdcdc;
  --json-tree-date-color: #dcdcdc;
}

ul {
  --string-color: var(--json-tree-string-color, #cb3f41);
  --symbol-color: var(--json-tree-symbol-color, #cb3f41);
  --boolean-color: var(--json-tree-boolean-color, #112aa7);
  --function-color: var(--json-tree-function-color, #112aa7);
  --number-color: var(--json-tree-number-color, #3029cf);
  --label-color: var(--json-tree-label-color, #871d8f);
  --arrow-color: var(--json-tree-arrow-color, #727272);
  --null-color: var(--json-tree-null-color, #8d8d8d);
  --undefined-color: var(--json-tree-undefined-color, #8d8d8d);
  --date-color: var(--json-tree-date-color, #8d8d8d);
  --li-identation: var(--json-tree-li-indentation, 1em);
  --li-line-height: var(--json-tree-li-line-height, 1.3);
  --li-colon-space: 0.3em;
  font-size: var(--json-tree-font-size, 12px);
  /* font-family: var(--json-tree-font-family, 'Courier New', Courier, monospace); */
  font-family: var(--json-tree-font-family, monospace);
}

ul :global(li) {
  line-height: var(--li-line-height);
  display: var(--li-display, list-item);
  list-style: none;
  white-space: nowrap;
}

ul,
ul :global(ul) {
  padding: 0;
  margin: 0;
}

ul.isDeleted {
  background: var(--theme-bg-volcano);
  background-image: url('data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAEElEQVQImWNgIAX8x4KJBAD+agT8INXz9wAAAABJRU5ErkJggg==');
  background-repeat: repeat-x;
  background-position: 50% 50%;
}

ul.isModified {
  background: var(--theme-bg-gold);
}

ul.isInserted {
  background: var(--theme-bg-green);
}
</style>
