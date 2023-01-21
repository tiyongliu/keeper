<template>
  <component
    :is="componentType"
    :name="name"
    :value="value"
    :isParentExpanded="isParentExpanded"
    :isParentArray="isParentArray"
    :nodeType="nodeType"
    :valueGetter="valueGetter"
    :expanded="expandedReplica"
    :labelOverride="labelOverride"
  />
</template>

<script lang="ts" setup>
import {computed, defineProps, inject, PropType, ref, ToRef, toRefs, unref} from 'vue'
import JSONObjectNode from './JSONObjectNode.vue'
import ErrorNode from './ErrorNode.vue'
import JSONArrayNode from './JSONArrayNode.vue'
import JSONIterableMapNode from './JSONIterableMapNode.vue'
import JSONIterableArrayNode from './JSONIterableArrayNode.vue'
import JSONMapEntryNode from './JSONMapEntryNode.vue'
import JSONValueNode from './JSONValueNode.vue'
import objType from './objType'
import {isUndefined} from 'lodash-es'

const props = defineProps({
  name: {
    type: [String, Number] as PropType<string | number>,
  },
  value: {
    type: [String, Object, Number, Boolean] as PropType<string | object>,
  },
  isParentExpanded: {
    type: Boolean as PropType<boolean>,
  },
  isParentArray: {
    type: Boolean as PropType<boolean>,
  },
  expanded: {
    type: Boolean as PropType<boolean>,
  },
  labelOverride: {
    type: String as PropType<string>,
    default: null
  }
})

const {name, value, isParentExpanded, isParentArray, expanded, labelOverride} = toRefs(props)
const getContext = inject('json-tree-default-expanded') as ToRef<boolean>
const expandedReplica = ref(expanded?.value)
if (isUndefined(expanded?.value)) {
  expandedReplica.value = !!unref(getContext)
}

const nodeType = computed(() => objType(value?.value))
const componentType = computed(() => getComponent(nodeType?.value))
const valueGetter = computed(() => getValueGetter(nodeType?.value))

function getComponent(nodeType) {
  switch (nodeType) {
    case 'Object':
      return JSONObjectNode;
    case 'Error':
      return ErrorNode;
    case 'Array':
      return JSONArrayNode;
    case 'Iterable':
    case 'Map':
    case 'Set':
      return typeof (value?.value as any).set === 'function' ? JSONIterableMapNode : JSONIterableArrayNode;
    case 'MapEntry':
      return JSONMapEntryNode
    default:
      return JSONValueNode
  }
}

function getValueGetter(nodeType) {
  switch (nodeType) {
    case 'Object':
    case 'Error':
    case 'Array':
    case 'Iterable':
    case 'Map':
    case 'Set':
    case 'MapEntry':
    case 'Number':
      return undefined;
    case 'String':
      return raw => `"${raw}"`;
    case 'Boolean':
      return raw => (raw ? 'true' : 'false');
    case 'Date':
      return raw => raw.toISOString();
    case 'Null':
      return () => 'null';
    case 'Undefined':
      return () => 'undefined';
    case 'Function':
    case 'Symbol':
      return raw => raw.toString();
    default:
      return () => `<${nodeType}>`;
  }
}
</script>
