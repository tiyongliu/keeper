<template>
  <div ref="node" class="main" @scroll="handleScroll">
    <div :style="`width: ${contentSize}px`">&nbsp;</div>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, PropType, ref, toRefs, unref, watch} from 'vue'

export default defineComponent({
  name: 'HorizontalScrollBar',
  props: {
    viewportRatio: {
      type: Number as PropType<number>,
      default: 0.5
    },
    minimum: {
      type: Number as PropType<number>
    },
    maximum: {
      type: Number as PropType<number>
    }
  },
  emits: ['scroll'],
  setup(props, {emit}) {
    const {viewportRatio, minimum, maximum} = toRefs(props)
    const node = ref<Nullable<HTMLElement>>(null)
    const width = ref(0)
    const contentSize = computed(() => Math.round(width.value / viewportRatio.value))

    function handleScroll() {
      const position = node.value!.scrollLeft
      const ratio = position / (contentSize.value - width.value);
      if (ratio < 0) return 0
      const res = ratio * (maximum.value! - minimum.value! + 1) + minimum.value!
      emit('scroll', Math.floor(res + 0.3))
    }

    watch(() => [viewportRatio.value, minimum.value, maximum.value], () => {
      width.value = node.value ? node.value.clientWidth : 0
    })

    function scroll(value) {
      const position01 = (unref(value) - minimum.value!) / (maximum.value! - minimum.value! + 1);
      const position = position01 * (contentSize.value - width.value);
      if (node.value) node.value.scrollLeft = Math.floor(position)
    }

    return {
      node,
      viewportRatio,
      minimum,
      maximum,
      handleScroll,
      contentSize,
      scroll
    }
  }
})
</script>

<style scoped>
.main {
  overflow-x: scroll;
  height: 16px;
  position: absolute;
  bottom: 0;
  right: 0;
  left: 0;
}
</style>
