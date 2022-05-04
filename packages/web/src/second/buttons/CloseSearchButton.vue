<template>
  <InlineButton v-if="filterValue || showDisabled" @click="clear">
    <close-outlined/>
  </InlineButton>
</template>

<script lang="ts">
  import {computed, defineComponent, unref} from 'vue'
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
      const filterValue = computed(() => unref(props.filter))

      function clear() {
        emit('close')
      }

      return {
        filterValue,
        clear
      }
    }
  })
</script>

<style scoped>

</style>
