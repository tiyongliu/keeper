<template>
  <div
    class="button"
    @click="handleClick"
    :class="[disabled && 'disabled', fillHorizontal && 'fillHorizontal']">
    <div class="icon">
      <FontIcon :icon="icon"/>
    </div>
    <div class="inner">
      <slot></slot>
    </div>
  </div>
</template>

<script lang="ts">
import {defineComponent, PropType, toRef, toRefs} from 'vue'
import FontIcon from '/@/second/icons/FontIcon.vue'

export default defineComponent({
  name: "LargeButton",
  components: {
    FontIcon
  },
  props: {
    icon: {
      type: String as PropType<string>,
    },
    disabled: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    fillHorizontal: {
      type: Boolean as PropType<boolean>,
      default: false
    },
  },
  emits: ['visible'],
  setup(props, {emit}) {
    const disabled = toRef(props, 'disabled')
    const handleClick = () => {
      if (!disabled.value) {
        emit('visible')
      }
    }

    return {
      handleClick,
      ...toRefs(props),
    }
  }
})
</script>

<!--<script lang="ts" setup>
import {defineEmits, defineProps, toRefs, withDefaults} from 'vue'
import FontIcon from '/@/second/icons/FontIcon.vue'
//<{ icon?: string, disabled: boolean, fillHorizontal: boolean }>
//<{icon?: string}, {disabled: boolean, fillHorizontal: boolean}>
// const defaultProps = withDefaults()
const defaultProps = withDefaults<any, any>(defineProps({
  icon: String
}), {
  fillHorizontal: Boolean,
  disabled: Boolean,
})

const {disabled, fillHorizontal} = toRefs(defaultProps)
const emit = defineEmits(['visible'])
const handleClick = () => {
  if (!disabled.value) {
    emit('visible')
  }
}
</script>-->

<style scoped>
.button {
  padding: 0 15px;
  color: var(--theme-font-1);
  border: 1px solid var(--theme-border);
  width: 120px;
  background-color: var(--theme-bg-1);
  cursor: pointer;
}

.button.fillHorizontal {
  width: auto;
  margin: 0 10px;
}

.button:not(.disabled):hover {
  background-color: var(--theme-bg-2);
}

.button:not(.disabled):active {
  background-color: var(--theme-bg-3);
}

.button.disabled {
  color: var(--theme-font-3);
}

.icon {
  font-size: 30px;
  text-align: center;
}

.inner {
  text-align: center;
}
</style>
