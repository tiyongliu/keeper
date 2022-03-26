<template>
    <SearchBoxWrapper>
      ConnectionList
    </SearchBoxWrapper>
</template>

<script lang="ts">
  import SearchBoxWrapper from './SearchBoxWrapper.vue'
  import { defineComponent, computed, unref, ref, onMounted } from 'vue';
  export default defineComponent({
    name: "ConnectionList",
    components: {
      SearchBoxWrapper
    },
    setup() {
      const hidden = ref(false)

      return {
        hidden
      }
    }
  })
</script>

<style scoped>

</style>
