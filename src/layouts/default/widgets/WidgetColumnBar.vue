<template>
  <div class="main-container" :class="hidden ? 'hidden' : ''">
    <slot />
  </div>
</template>

<script lang="ts">
  import { defineComponent, computed, unref, ref, onMounted } from 'vue';
  export default defineComponent({
    name: "DatabaseWidget",
    components: {

    },
    setup() {
      const hidden = ref(false)

      return {
        hidden
      }
    }
  })
</script>

<style>
  .hidden {
    display: none;
  }

  .main-container {
    position: relative;
    flex: 1;
    flex-direction: column;
    user-select: none;
  }

  .main-container:not(.hidden) {
    display: flex;
  }
</style>
