<template>
  <div class="button" :class="disabled && 'disabled'" :title="title">
    <div class="inner" :class="disabled && 'disabled'" @click="handleClick($event)">
      <span class="icon" :class="disabled && 'disabled'"><FontIcon :icon="icon"/></span>
      <slot></slot>
    </div>
  </div>
</template>

<script lang="ts">
import {defineComponent, PropType, toRef, toRefs} from 'vue'
import FontIcon from '/@/second/icons/FontIcon.vue'

export default defineComponent({
  name: "ToolStripButton",
  components: {
    FontIcon
  },
  props: {
    title: {
      type: String as PropType<string>
    },
    disabled: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    icon: {
      type: String as PropType<string>
    }
  },
  emits: ['click'],
  setup(props, {emit}) {
    const disabled = toRef(props, 'disabled')

    function handleClick(e: Event) {
      if (disabled.value) return
      emit('click', {target: e.target})
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
  padding-left: 5px;
  padding-right: 5px;
  color: var(--theme-font-1);
  border: 0;
  align-self: stretch;
  display: flex;
  user-select: none;
}

.button.disabled {
  color: var(--theme-font-3);
}

.inner:hover:not(.disabled) {
  background: var(--theme-bg-3);
}

.inner:active:hover:not(.disabled) {
  background: var(--theme-bg-4);
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
  background: var(--theme-bg-2);
  padding: 3px 8px;
  border-radius: 4px;
  cursor: pointer;
}
</style>
