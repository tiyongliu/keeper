<template>
  <a-input
    size="small"
    autoComplete="off"
    :placeholder="placeholder"
    v-model:value="searchValue"
    @input="handleInput"
    @keydown="handleKeyDown"
    allow-clear/>
</template>

<script lang="ts">
  import {defineComponent, ref, unref, watch, toRef} from 'vue'
  import {debounce} from 'lodash-es'
  import keycodes from '/@/second/utility/keycodes'

  export default defineComponent({
    name: 'SearchInput',
    props: {
      placeholder: {
        type: String as PropType<string>,
      },
      isDebounced: {
        type: Boolean as PropType<boolean>,
      },
      value: {
        type: String as PropType<string>,
      }
    },
    emits: ['update:value'],
    setup(props, {emit}) {
      const isDebounced = toRef(props, 'isDebounced')
      const value = toRef(props, 'value')

      const searchValue = ref<string>('')
      const debouncedSet = debounce(x => {emit('update:value', unref(x))}, 500)

      function handleKeyDown(e) {
        if (e.keyCode == keycodes.escape) {
          searchValue.value = ''
          emit('update:value', '')
        }
      }

      function handleInput() {
        if (unref(isDebounced)) {
          debouncedSet(searchValue.value)
        } else {
          emit('update:value', searchValue.value)
        }
      }

      watch(() => value.value, () => {
        if (value.value === '') searchValue.value = ''
      })

      return {
        value,
        searchValue,
        placeholder: props.placeholder,
        debouncedSet,
        handleKeyDown,
        handleInput,
      }
    }
  })
</script>
