<template>
  <a-input
    size="small"
    :placeholder="placeholder"
    v-model:value="searchValue"
    @input="handleInput"
    @keydown="handleKeyDown"
    allow-clear/>
</template>

<script lang="ts">
  import {defineComponent, ref, unref, watch, onBeforeUnmount} from 'vue';
  import {debounce} from 'lodash-es'
  import keycodes from '/@/second/utility/keycodes'

  export default defineComponent({
    name: "SearchInput",
    props: {
      placeholder: {
        type: String as PropType<string>,
      },
      isDebounced: {
        type: Boolean as PropType<false>,
      },
      searchValue: {
        type: String as PropType<string>,
      },
      value: {
        type: String as PropType<string>,
      }
    },
    emits: ['update:searchValue'],
    setup(props, {emit}) {
      const searchValue = ref<string>('');
      const value = ref<string>('');
      const {
        isDebounced,
      } = unref(props) as unknown as {
        isDebounced: boolean,
      }


      const debouncedSet = debounce(x => (value.value = x), 500)

      function handleKeyDown(e) {
        if (e.keyCode == keycodes.escape) {
          searchValue.value = ''
          value.value = ''
        }
      }

      function handleInput() {
        if (unref(isDebounced)) {
          debouncedSet(searchValue.value)
        } else {
          value.value = searchValue.value
        }
      }

      watch(() => unref(value), () => {
        emit('update:searchValue', unref(value))
      })

      watch(() => unref(props.searchValue), () => {
        if (props.searchValue === '') {
          searchValue.value = ''
        }
      })

      onBeforeUnmount(() => {
        searchValue.value = ''
        value.value = ''
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

<style scoped>
  input {
    flex: 1;
    min-width: 10px;
    width: 10px;
    border: none;
  }
</style>
