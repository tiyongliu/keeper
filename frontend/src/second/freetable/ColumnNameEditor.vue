<template>
  <input
    type="text"
    v-bind="$attrs"
    v-model="value"
    ref="domEditor"
    @keydown=""
    @blur=""
    :class="isError && 'isError'"
  />
</template>

<script lang="ts">
import {computed, defineComponent, ref, toRefs} from 'vue'
import keycodes from '/@/second/utility/keycodes'
export default defineComponent({
  name: 'ColumnNameEditor',
  props: {
    focusOnCreate: {
      type: Boolean as PropType<false>,
      default: false
    },
    blurOnEnter: {
      type: Boolean as PropType<false>,
      default: false
    },
    existingNames: {
      type: Array as PropType<string[]>,
    },
    defaultValue: {
      type: String as PropType<string>,
      default: ''
    }
  },
  setup(props) {
    const {defaultValue, existingNames} = toRefs(props)
    const domEditor = ref<Nullable<HTMLElement>>(null)
    const value = ref(defaultValue.value || '')

    const isError = computed(() => value.value && existingNames.value && existingNames.value?.includes(value.value))

    const handleKeyDown = (event) => {
      if (value.value && event.keyCode == keycodes.enter && !isError.value) {

      }
    }

    const handleBlur = () => {

    }
    return {
      domEditor,
      defaultValue,
      value,
      handleKeyDown,
      handleBlur,
      isError
    }
  }
})
</script>

<style scoped>
input {
  width: calc(100% - 10px);
}

input.isError {
  background: var(--theme-bg-red);
}
</style>
