<template>
  <li
    ref="domElement"
    :class="[isParentExpanded && 'indent', !!elementValue && 'jsonValueHolder']">
    <label>
      <JSONArrow
        v-if="expandable && isParentExpanded"
        @click="toggleExpand"
        :expanded="expandedRW"/>
      <JSONKey
        :name="name"
        :colon="context['colon']"
        :isParentExpanded="isParentExpanded"
        :isParentArray="isParentArray"
        @click="toggleExpand"/>
      <span @click="toggleExpand"><span>{{ label }}</span>{{ bracketOpen }}</span>
    </label>
    <ul v-if="isParentExpanded" :class="!expandedRW && 'collapse'" @click="expand">
      <template :key="getKey(slicedKey)" v-for="(slicedKey, index) in slicedKeys">
        <JSONNode
          :name="getKey(slicedKey)"
          :isParentExpanded="expandedRW"
          :isParentArray="isArray"
          :value="expandedRW ? getValue(slicedKey) : getPreviewValueReplica(slicedKey)"/>
        <span v-if="!expandedRW && index < previewKeysReplica.length - 1" class="comma">,</span>
      </template>
      <span v-if="slicedKeys.length < previewKeysReplica.length">…</span>
    </ul>
    <span v-else>…</span>
    <span>{{ bracketClose }}</span>
  </li>
</template>

<script lang="ts" setup>
import {
  computed,
  defineProps,
  inject,
  PropType,
  provide,
  ref,
  ToRef,
  toRefs,
  watchEffect,
} from 'vue'
import {isUndefined} from 'lodash-es'
import JSONNode from './JSONNode.vue'
import JSONArrow from './JSONArrow.vue'
import JSONKey from './JSONKey.vue'

const props = defineProps({
  name: {
    type: String as PropType<string>,
  },
  keys: {
    type: Array as PropType<string[]>,
  },
  colon: {
    type: String as PropType<string>,
    default: ':'
  },
  label: {
    type: String as PropType<string>,
    default: ''
  },
  isParentExpanded: {
    type: Boolean as PropType<boolean>,
  },
  isParentArray: {
    type: Boolean as PropType<boolean>,
  },
  isArray: {
    type: Boolean as PropType<boolean>,
    default: false
  },
  bracketOpen: {
    type: String as PropType<string>,
  },
  bracketClose: {
    type: String as PropType<string>,
  },
  previewKeys: {
    type: Array as PropType<string[]>,
  },
  getKey: {
    type: Function as PropType<<T>(key: T) => T>,
    default: key => key
  },
  getValue: {
    type: Function as PropType<<T>(key: T) => T>,
    default: key => key
  },
  getPreviewValue: {
    type: Function as PropType<<T>(key: T) => T>
  },
  expanded: {
    type: Boolean as PropType<boolean>,
    default: false
  },
  expandable: {
    type: Boolean as PropType<boolean>,
    default: true
  },
  elementValue: {
    type: [String, Number, Boolean, Object, Array] as PropType<string | boolean | number | object | string[]>,
  }
})

const {
  name,
  keys,
  colon,
  label,
  isParentExpanded,
  isParentArray,
  isArray,
  bracketOpen,
  bracketClose,
  previewKeys,
  getKey,
  getValue,
  getPreviewValue,
  expanded,
  expandable,
  elementValue,
} = toRefs(props)

const previewKeysReplica = ref(previewKeys?.value)
const getPreviewValueReplica = ref(getPreviewValue?.value)
const expandedRW = ref(expanded?.value)

if (isUndefined(previewKeysReplica.value)) previewKeysReplica.value = keys?.value

if (isUndefined(getPreviewValueReplica.value)) getPreviewValueReplica.value = getValue?.value

const domElement = ref<Nullable<HTMLElement>>(null)

const context = inject('json-tree-context-key') as { [key in string]: any }
provide('json-tree-context-key', {...context, colon: colon.value})
const elementData = inject('json-tree-element-data') as WeakMap<object, any>
const slicedKeyCount = inject('json-tree-sliced-key-count') as ToRef<number>

const slicedKeys = computed(() => expandedRW.value ? keys?.value : previewKeysReplica.value!.slice(0, slicedKeyCount.value || 5))

watchEffect(() => {
  if (!isParentExpanded?.value) {
    //todo expanded = !expanded;
    expandedRW.value = false
  }
})

watchEffect(() => {
  if (domElement.value && elementData && elementValue && elementValue.value) {
    elementData.set(domElement.value, elementValue.value)
  }
})

function toggleExpand() {
  expandedRW.value = !expandedRW?.value
}

function expand() {
  expandedRW.value = true
}
</script>
<style scoped>
label {
  display: inline-block;
}

.indent {
  padding-left: var(--li-identation);
}

.collapse {
  --li-display: inline;
  display: inline;
  font-style: italic;
}

.comma {
  margin-left: -0.5em;
  margin-right: 0.5em;
}

label {
  /* display: contents; */
  position: relative;
}
</style>
