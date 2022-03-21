<template>
    <WidgetColumnBar>
      <WidgetColumnBarItem title="Connections" height="35%">
        DatabaseWidget
      </WidgetColumnBarItem>

    </WidgetColumnBar>
</template>

<script lang="ts">
  import { defineComponent, computed, unref, ref, onMounted } from 'vue';
  import WidgetColumnBar from './WidgetColumnBar.vue'
  import WidgetColumnBarItem from './WidgetColumnBarItem.vue'
  export default defineComponent({
    name: "DatabaseWidget",
    components: {
      WidgetColumnBar,
      WidgetColumnBarItem
    }
  })
</script>

<style scoped>

</style>
