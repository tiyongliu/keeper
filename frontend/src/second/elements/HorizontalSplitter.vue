<template>
  <div class="container" ref="container">
    <template v-if="!hideFirst">
      <div class="child1" :style="isSplitter
      ? collapsed1
         ? 'display:none'
         : collapsed2
         ? 'flex:1'
         : `width:${sizeRw}px; min-width:${sizeRw}px; max-width:${sizeRw}px}`
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
      <div :class="collapsed1 ? 'child1' : 'child2'"
           :style="collapsed2 ? 'display:none' : collapsed1 ? 'flex:1' : 'child2'">
        <slot name="2"></slot>
      </div>
    </template>

    <template v-if="allowCollapseChild1 && !collapsed2 && isSplitter">
      <div v-if="collapsed1" class="collapse" style="left: 0px" @click="() => collapsed1 = false">
        <FontIcon icon="icon chevron-double-right"/>
      </div>
      <div v-else class="collapse" :style="`left: ${sizeRw - 16}px`"
           @click="() => collapsed1 = true">
        <FontIcon icon="icon chevron-double-left"/>
      </div>
    </template>

    <template v-if="allowCollapseChild2 && !collapsed1 && isSplitter">
      <div v-if="collapsed2" class="collapse" style="right: 0px" @click="() => collapsed2 = false">
        <FontIcon icon="icon chevron-double-left"/>
      </div>
      <div v-else class="collapse" :style="`left: ${sizeRw}px`" @click="() => collapsed2 = true">
        <FontIcon icon="icon chevron-double-left"/>
      </div>
    </template>
  </div>
</template>

<script lang="ts">
import {defineComponent, onMounted, PropType, ref, toRef, toRefs, watch} from 'vue'
import {isString} from 'lodash-es'
import FontIcon from '/@/second/icons/FontIcon.vue'
import bus from '/@/second/utility/bus'

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
      type: Object as PropType<object>,
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
  emits: ['update:size'],
  setup(props, {emit}) {
    const size = toRef(props, 'size')
    const initialValue = toRef(props, 'initialValue')
    const collapsed1 = ref(false)
    const collapsed2 = ref(false)
    const sizeRw = ref(size.value)

    const container = ref<Nullable<HTMLElement>>(null)

    bus.emitter.on(bus.resize, updateWidgetStyle)

    watch(() => initialValue.value, updateWidgetStyle)

    function updateWidgetStyle() {
      if (container.value) {
        sizeRw.value = computeSplitterSize(initialValue.value, container.value.clientWidth)
      }
    }

    function updateSize(e: number) {
      sizeRw.value += e
      emit('update:size', sizeRw.value)
    }

    onMounted(async () => {
      await updateWidgetStyle()
      emit('update:size', sizeRw.value)
    })

    return {
      container,
      ...toRefs(props),
      sizeRw,
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
