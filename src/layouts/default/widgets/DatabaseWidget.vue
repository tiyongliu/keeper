<template>
    <WidgetColumnBar :hidden="hidden">
      <WidgetColumnBarItem title="Connections" name="connections" height="35%" storageName="connectionsWidget">
        DatabaseWidget isHidden {{hidden}}
      </WidgetColumnBarItem>

    </WidgetColumnBar>
</template>

<script lang="ts">
  import {defineComponent, computed, unref, ref, onMounted, PropType} from 'vue';
  import WidgetColumnBar from './WidgetColumnBar.vue'
  import WidgetColumnBarItem from './WidgetColumnBarItem.vue'
  export default defineComponent({
    name: "DatabaseWidget",
    props: {
      hidden: {
        type: [Boolean] as PropType<boolean>,
        default: false,
      }
    },
    components: {
      WidgetColumnBar,
      WidgetColumnBarItem
    },
    setup() {

    }
  })
</script>

<style scoped>

</style>
