<template>
  <div :style="`max-width: ${width}px`" :class="isFlex && 'isFlex'">
    <slot></slot>
  </div>
</template>

<script lang="ts">
import {defineComponent, PropType} from 'vue'
export default defineComponent({
  name: "ManagerInnerContainer",
  props: {
    width: {
      type: Number as PropType<number>,
    },
    isFlex: {
      type: Boolean as PropType<boolean>,
      default: false
    },
  },
  setup(props) {
    return {
      props
    }
  }
})
</script>

<style scoped>
div {
  flex: 1 1;
  overflow-y: auto;
  overflow-x: auto;
}

div.isFlex {
  display: flex;
}
</style>
