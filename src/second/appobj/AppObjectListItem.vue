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
      },
      isHidden: {
        type: Boolean as unknown as PropType<boolean>,
      },
      isExpandedBySearch: {
        type: Boolean as unknown as PropType<boolean>,
        default: false
      },
      expandOnClick: {
        type: Boolean as PropType<boolean>,
      },
    },
    components: {},
    setup(props) {
      // const dynamicList = computed(() => unref(props.list))
      console.log(props)

      async function handleExpand() {
        // if (subItemsComponent && expandOnClick) {
        //   await tick();
        //   isExpanded = !isExpanded;
        // }
      }

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
