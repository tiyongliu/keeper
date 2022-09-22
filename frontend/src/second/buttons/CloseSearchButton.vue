<template>
  <InlineButton v-if="filter || showDisabled" @click="clear">
    <close-outlined/>
  </InlineButton>
</template>

<script lang="ts">
  import {defineComponent, toRef} from 'vue'
  import {CloseOutlined} from '@ant-design/icons-vue'
  import InlineButton from './InlineButton.vue'

  export default defineComponent({
    name: "CloseSearchButton",
    props: {
      filter: {
        type: String as PropType<string>,
      },
      showDisabled: {
        type: Boolean as PropType<false>,
        default: false
      }
    },
    components: {
      InlineButton,
      CloseOutlined
    },
    emits: ['close'],
    setup(props, {emit}) {
      const filter = toRef(props, 'filter')

      function clear() {
        emit('close')
      }

      return {
        filter,
        clear
      }
    }
  })
</script>
