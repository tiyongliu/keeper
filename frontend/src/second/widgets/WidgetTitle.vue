<template>
    <div class="widgetTitle">
      <slot />
    </div>
</template>

<script lang="ts">
  import { defineComponent, computed, unref, ref, onMounted } from 'vue';
  export default defineComponent({
    name: "WidgetTitle",
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

<style lang="less">
  .widgetTitle {
    padding: 5px;
    font-weight: bold;
    text-transform: uppercase;
    background-color: var(--theme-bg-1);
    border: 2px solid var(--theme-border);
    &:hover {
       background-color: var(--theme-bg-2);
     }
  }
</style>

