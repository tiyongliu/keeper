<template>
  <WidgetsInnerContainer>
    PinnedObjectsList
  </WidgetsInnerContainer>
</template>

<script lang="ts">
  import {
    defineComponent,
    computed,
    inject,
    ref,
    onMounted,
    watch,
    toRefs,
    reactive,
    unref
  } from 'vue';
  import WidgetsInnerContainer from './WidgetsInnerContainer.vue'
  export default defineComponent({
    components: {
      WidgetsInnerContainer
    }
  })
</script>

<style scoped>

</style>
