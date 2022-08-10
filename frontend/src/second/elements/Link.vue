<template>
  <a href="">
    <slot></slot>
  </a>
</template>

<script lang="ts">
  import {defineComponent, toRefs, PropType} from 'vue'
  // import FontIcon from '/@/second/icons/FontIcon.vue'
  export default defineComponent({
    name: "Link",
    props: {
      href: {
        type: String as PropType<string>,
      },
    },
    // components: {FontIcon},
    setup(props) {
      return {
        ...toRefs(props)
      }
    }
  })
</script>

<style scoped>
  a {
    text-decoration: none;
    cursor: pointer;
    color: var(--theme-font-link);
  }
  a:hover {
    text-decoration: underline;
  }
</style>

