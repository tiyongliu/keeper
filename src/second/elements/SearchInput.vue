<template>
  <Input
    size="small"
    :placeholder="placeholder"
    v-model:value="searchValue"
    @input="handleInput"
    @keydown="handleKeyDown"
    allow-clear/>
  <!--
  on:input={e => {
    if (isDebounced) debouncedSet(domInput.value);
    else value = domInput.value;
  }}
  -->
</template>

<script lang="ts">
  import {defineComponent, ref, unref, watch, onBeforeUnmount} from 'vue';
  import {Input} from 'ant-design-vue'
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
      }
    },
    components: {
      Input
    },
    emits: ['update:value'],
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
          emit('update:value', searchValue)
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
        emit('update:value', unref(value))
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
