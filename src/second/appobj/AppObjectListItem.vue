<template>
  <div class="subitems"></div>
</template>

<script lang="ts">
  import {computed, defineComponent, PropType, unref} from 'vue';
  export default defineComponent({
    name: "AppObjectListItem",
    props: {
      data: {
        type: Array as unknown as PropType<[]>,
        default: [],
      }
    },
    components: {},
    setup(props) {
      // const dynamicList = computed(() => unref(props.list))

      return {

      }
    }
  })
</script>

<style lang="less">
  .subitems {
    margin-left: 28px;
  }
</style>
