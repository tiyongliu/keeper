<template>
  <div class="button" :class="disabled && 'disabled'" :title="title">
    <div class="inner" :class="disabled && 'disabled'">
      <div class="main" :class="disabled && 'disabled'" @click="handleClick">
        <span class="icon" :class="disabled && 'disabled'"><FontIcon :icon="icon"/></span>
        <slot/>
      </div>
      <span class="split-icon" :class="disabled && 'disabled'" @click="handleSplitClick">
        <FontIcon :icon="splitIcon"/></span>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {defineEmits, defineProps, PropType, toRefs} from 'vue'
import FontIcon from '/@/second/icons/FontIcon.vue'

const props = defineProps({
  disabled: {
    type: Boolean as PropType<boolean>,
    default: false
  },
  icon: {
    type: String as PropType<Nullable<string>>,
  },
  title: {
    type: String as PropType<Nullable<string>>,
  },
  splitIcon: {
    type: String as PropType<string>,
    default: 'icon chevron-down'
  }
})

const {disabled, icon, title, splitIcon} = toRefs(props)

const emit = defineEmits(['childClick', 'splitclick'])

function handleClick(e) {
  if (disabled?.value) return;
  emit('childClick', {target: e.target})
}

function handleSplitClick(e) {
  if (disabled?.value) return;
  emit('splitclick', {target: e.target});
}
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

.main {
  background: var(--theme-bg-2);
  padding: 3px 0px 3px 8px;
  border-radius: 4px 0px 0px 4px;
}

.main:hover:not(.disabled) {
  background: var(--theme-bg-3);
}

.main:active:hover:not(.disabled) {
  background: var(--theme-bg-4);
}

.split-icon:hover:not(.disabled) {
  background: var(--theme-bg-3);
}

.split-icon:active:hover:not(.disabled) {
  background: var(--theme-bg-4);
}

.split-icon {
  background: var(--theme-bg-2);
  padding: 3px 8px 3px 0px;
  border-radius: 0px 4px 4px 0px;
}

.icon {
  margin-right: 5px;
  color: var(--theme-font-link);
}

.icon.disabled {
  color: var(--theme-font-3);
}

.inner {
  white-space: nowrap;
  align-self: center;
  cursor: pointer;
  display: flex;
}

.main {
  display: flex;
  padding-right: 5px;
}

.split-icon {
  padding-left: 5px;
  color: var(--theme-font-link);
  border-left: 1px solid var(--theme-bg-4);
}
</style>
