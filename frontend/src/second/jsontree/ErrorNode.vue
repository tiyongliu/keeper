<template>
  <li :class="isParentExpanded && 'indent'">
    <JSONArrow v-if="isParentExpanded" @click="toggleExpand" :expanded="expanded"/>
    <JSONKey
      :colon="context['colon']"
      :isParentExpanded="isParentExpanded"
      :isParentArray="isParentArray"/>
    <span @click="toggleExpand">Error: {{ expanded ? '' : value.message }}</span>
    <ul :class="!expanded && 'collapse'">
      <template v-if="expanded">
        <JSONNode name="message" :value="value.message"/>
        <li>
          <JSONKey name="stack" colon=":" :isParentExpanded="isParentExpanded"/>
          <span v-for="(line, index) in stack">
            <span :class="index > 0 && 'indent'">{{ line }}</span><br/>
          </span>
        </li>
      </template>
    </ul>
  </li>
</template>

<script lang="ts" setup>
import {computed, defineProps, inject, PropType, provide, toRefs} from 'vue'
import JSONArrow from './JSONArrow.vue'
import JSONKey from './JSONKey.vue'
import JSONNode from './JSONKey.vue'

const props = defineProps({
  name: {
    type: String as PropType<string>,
  },
  value: {
    type: Object as PropType<{ stack: string; message: string }>,
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

const {value, isParentExpanded, isParentArray} = toRefs(props)

const stack = computed(() => value?.value ? value?.value.stack.split('\n') : [])

const context = inject('json-tree-context-key') as { [key in string]: any }
provide('json-tree-context-key', {...context, colon: ':'})

function toggleExpand() {

}

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

.collapse {
  --li-display: inline;
  display: inline;
  font-style: italic;
}
</style>
