<template>
  <WidgetTitle v-if="!skip && show">{{title}}</WidgetTitle>
  <div class="wrapper">
    <slot/>
  </div>
</template>

<script lang="ts">
  import {defineComponent, computed, unref, ref, onMounted} from 'vue';
  import WidgetTitle from './WidgetTitle.vue'

  export default defineComponent({
    name: "WidgetColumnBarItem",
    components: {
      WidgetTitle
    },
    props: {
      title: {
        type: String as PropType<string>,
      },
      name: {
        type: String as PropType<string>,
      },
      height: {
        type: String as PropType<string>,
      }
    },
    setup(props) {
      const skip = ref(false)
      const show = ref(true)

      return {
        skip,
        show,
        title: props.title,
        height: props.height,
      }
    }
  })
</script>

<style scoped>
  .wrapper {
    overflow: hidden;
    position: relative;
    flex-direction: column;
    display: flex;
  }
</style>
