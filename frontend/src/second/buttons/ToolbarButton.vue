<template>
  <div class="button" @click="handleClick" :class="disabled && 'disabled'" :title="title">
    <div class="inner">
      <img v-if="externalImage" :src="externalImage" />
      <span>
        <FontIcon :class="disabled && 'disabled'" :icon="icon" />
        <slot></slot>
      </span>
    </div>
  </div>
</template>

<script lang="ts">
import {defineComponent, toRef, toRefs} from 'vue'
import FontIcon from '/@/second/icons/FontIcon.vue'

export default defineComponent({
  name: "ToolbarButton",
  components: {
    FontIcon
  },
  props: {
    disabled: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    externalImage: {
      type: String as PropType<string>
    },
    icon: {
      type: String as PropType<string>
    },
    title: {
      type: String as PropType<string>
    }
  },
  emits: ['click'],
  setup(props, {emit}) {
    const disabled = toRef(props, 'disabled')
    function handleClick() {
      if (disabled.value) return
      emit('click')
    }
    return {
      ...toRefs(props),
      handleClick
    }
  }
})
</script>

<style scoped>
.button {
  /* padding: 5px 15px; */
  padding-left: 15px;
  padding-right: 15px;
  color: var(--theme-font-1);
  border: 0;
  border-right: 1px solid var(--theme-border);
  align-self: stretch;
  display: flex;
  user-select: none;
}
.button.disabled {
  color: var(--theme-font-3);
}
.button:hover:not(.disabled) {
  background: var(--theme-bg-2);
}
.button:active:hover:not(.disabled) {
  background: var(--theme-bg-3);
}
.icon {
  margin-right: 5px;
  color: var(--theme-font-link);
}
.icon.disabled {
  color: var(--theme-font-3);
}
.inner {
  /* position: relative;
  top: 2px; */
  white-space: nowrap;
  align-self: center;
}
img {
  width: 20px;
  height: 20px;
}
</style>
