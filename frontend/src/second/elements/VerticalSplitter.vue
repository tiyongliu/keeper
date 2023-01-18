<template>
  <div class="container" ref="container">
    <div
      class="child1"
      :style="isSplitter
      ? collapsed1
        ? 'display:none'
        : collapsed2
        ? 'flex:1'
        : `height:${size}px; min-height:${size}px; max-height:${size}px}`
      : `flex:1`"
    >
      <slot name="1"></slot>
    </div>
    <div
      v-if="isSplitter"
      class="vertical-split-handle"
      :style="collapsed1 || collapsed2 ? 'display:none' : ''"
      v-splitterDrag="'clientY'"
      :resizeSplitter="(e) => updateSize(e.detail)">
      <div :class="collapsed1 ? 'child1' : 'child2'"
           :style="collapsed2 ? 'display:none' : collapsed1 ? 'flex:1' : 'child2'">
        <slot name="2"></slot>
      </div>

      <template v-if="allowCollapseChild1 && !collapsed2 && isSplitter">
        <div v-if="collapsed1" class="collapse" style="top: 0px" @click="() => collapsed1 = false">
          <FontIcon icon="icon chevron-double-down"/>
        </div>
        <div v-else class="collapse" :style="`top: ${size - 16}px`"
             @click="() => collapsed1 = true">
          <FontIcon icon="icon chevron-double-up"/>
        </div>
      </template>

      <template v-if="allowCollapseChild2 && !collapsed1 && isSplitter">
        <div v-if="collapsed2" class="collapse" style="bottom: 0px"
             @click="() => collapsed2 = false">
          <FontIcon icon="icon chevron-double-up"/>
        </div>

        <div v-else class="collapse" :style="`top: ${size}px`" @click="() => collapsed2 = true">
          <FontIcon icon="icon chevron-double-down"/>
        </div>
      </template>
    </div>
  </div>
</template>

<script lang="ts">
import {defineComponent, onMounted, ref, toRef, toRefs, watch} from 'vue'
import FontIcon from '/@/second/icons/FontIcon.vue'
import {computeSplitterSize} from '/@/second/elements/HorizontalSplitter.vue'
import bus from "/@/second/utility/bus";

export default defineComponent({
  name: "VerticalSplitter",
  components: {
    FontIcon
  },
  props: {
    isSplitter: {
      type: Boolean as PropType<boolean>,
      default: true
    },
    allowCollapseChild1: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    allowCollapseChild2: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    initialValue: {
      type: String as PropType<string>
    }
  },
  setup(props) {
    const container = ref<Nullable<HTMLElement>>(null)
    const collapsed1 = ref(false)
    const collapsed2 = ref(false)
    const initialValue = toRef(props, 'initialValue')
    const size = ref(0)

    bus.emitter.on(bus.resize, updateWidgetStyle)

    watch(() => initialValue.value, updateWidgetStyle)

    function updateWidgetStyle() {
      if (container.value) {
        size.value = computeSplitterSize(initialValue.value, container.value.clientHeight)
      }
    }

    function updateSize(e: number) {
      size.value += e
    }

    onMounted(async () => {
      await updateWidgetStyle()
    })

    return {
      container,
      ...toRefs(props),
      collapsed1,
      collapsed2,
      size,
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
  flex-direction: column;
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
  right: 16px;
  width: 40px;
  height: 16px;
  background: var(--theme-bg-2);
  display: flex;
  justify-content: center;
  z-index: 100;
}

.collapse:hover {
  color: var(--theme-font-hover);
  background: var(--theme-bg-3);
  cursor: pointer;
}
</style>
