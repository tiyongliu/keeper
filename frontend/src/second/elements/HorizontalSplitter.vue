<template>
  <div class="container" ref="clientWidth">
    <template v-if="!hideFirst">
      <div class="child1" :style="isSplitter
      ? collapsed1
         ? 'display:none'
         : collapsed2
         ? 'flex:1'
         : `width:${size}px; min-width:${size}px; max-width:${size}px}`
      : `flex:1`">
        <slot name="1"></slot>
      </div>
    </template>
    <template v-if="isSplitter">
      <div
        v-if="!hideFirst"
        class="horizontal-split-handle"
        :style="collapsed1 || collapsed2 ? 'display:none' : ''"
        v-splitterDrag="'clientX'"
        :resizeSplitter="(e) => updateSize(e.detail)">
      </div>
      <div :class="collapsed1 ? 'child1' : 'child2'" :style="collapsed2 ? 'display:none' : collapsed1 ? 'flex:1' : 'child2'">
        <slot name="2"></slot>
      </div>
    </template>

    <template v-if="allowCollapseChild1 && !collapsed2 && isSplitter">
      <div v-if="collapsed1" class="collapse" style="left: 0px" @click="() => collapsed1 = false">
        <FontIcon icon="icon chevron-double-right" />
      </div>
      <div v-else class="collapse" :style="`left: ${size - 16}px`" @click="() => collapsed1 = true">
        <FontIcon icon="icon chevron-double-left" />
      </div>
    </template>

    <template v-if="allowCollapseChild2 && !collapsed1 && isSplitter">
      <div v-if="collapsed2" class="collapse" style="right: 0px"  @click="() => collapsed2 = false">
        <FontIcon icon="icon chevron-double-left" />
      </div>
      <div v-else class="collapse" :style="`left: ${size}px`" @click="() => collapsed2 = true">
        <FontIcon icon="icon chevron-double-left" />
      </div>
    </template>
  </div>
</template>

<script lang="ts">
import {defineComponent, nextTick, PropType, ref, toRef, toRefs, watch} from 'vue'
import {isString, omit} from 'lodash-es'
import FontIcon from '/@/second/icons/FontIcon.vue'

export function computeSplitterSize(initialValue, clientSize) {
  if (isString(initialValue) && initialValue.startsWith('~') && initialValue.endsWith('px'))
    return clientSize - parseInt(initialValue.slice(1, -2));
  if (isString(initialValue) && initialValue.endsWith('px')) return parseInt(initialValue.slice(0, -2));
  if (isString(initialValue) && initialValue.endsWith('%'))
    return (clientSize * parseFloat(initialValue.slice(0, -1))) / 100;
  return clientSize / 2;
}

export default defineComponent({
  name: "HorizontalSplitter",
  components: {
    FontIcon
  },
  props: {
    isSplitter: {
      type: Boolean as PropType<boolean>,
      default: true
    },
    initialValue: {
      type: String as PropType<string>,
    },
    hideFirst: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    allowCollapseChild1: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    allowCollapseChild2: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    size: {
      type: Number as PropType<number>,
      default: 0
    }
  },
  setup(props) {
    const collapsed1 = ref(false)
    const collapsed2 = ref(false)

    const clientWidth = ref<Nullable<HTMLElement>>(null)
    const initialValue = toRef(props, 'initialValue')
    const size = toRef(props, 'size')

    watch(() => [initialValue.value, clientWidth.value], async () => {
      await nextTick()
      size.value = computeSplitterSize(initialValue.value, clientWidth.value)
    })

    function updateSize(e: number) {
      size.value += e
    }

    return {
      clientWidth,
      ...toRefs(omit(props, ['size'])),
      size,
      collapsed1,
      collapsed2,
      updateSize
    }
  }
})
</script>

<style scoped>
.container {
  flex: 1;
  display: flex;
  position: absolute;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
}

.child1 {
  display: flex;
  position: relative;
  overflow: hidden;
}

.child2 {
  flex: 1;
  display: flex;
  position: relative;
  overflow: hidden;
}

.collapse {
  position: absolute;
  bottom: 16px;
  height: 40px;
  width: 16px;
  background: var(--theme-bg-2);
  display: flex;
  flex-direction: column;
  justify-content: center;
  z-index: 100;
}

.collapse:hover {
  color: var(--theme-font-hover);
  background: var(--theme-bg-3);
  cursor: pointer;
}
</style>
