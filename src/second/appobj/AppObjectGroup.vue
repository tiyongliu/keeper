<template>
  <div class="group">
  <span class="expand-icon">

  </span>

  </div>
</template>

<script lang="ts">
  import {defineComponent, computed, unref, ref, onMounted, PropType} from 'vue';
  export default defineComponent({
    name: "AppObjectGroup"
  })
</script>

<style scoped>
  .group {
    user-select: none;
    padding: 5px;
    cursor: pointer;
    white-space: nowrap;
    font-weight: bold;
  }

  .group:hover {
    background-color: var(--theme-bg-hover);
  }
  .expand-icon {
    margin-right: 3px;
  }
</style>
