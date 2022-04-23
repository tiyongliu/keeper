<template>
  <WidgetColumnBarItem title="Database content" name="dbObjects" storageName="dbObjectsWidget">
    <WidgetsInnerContainer>
      sdfsdf
    </WidgetsInnerContainer>
  </WidgetColumnBarItem>
</template>

<script lang="ts">
  import {defineComponent, PropType, computed, watch, unref} from 'vue';
  import WidgetColumnBarItem from './WidgetColumnBarItem.vue'
  import WidgetsInnerContainer from './WidgetsInnerContainer.vue'
  export default defineComponent({
    name: "SqlObjectList",
    props: {
      conid: {
        type: String as PropType<string>
      },
      database: {
        type: String as PropType<string>
      }
    },
    components: {
      WidgetColumnBarItem,
      WidgetsInnerContainer
    },
    setup() {
      //todo
      // $: status = useDatabaseStatus({ conid, database });

      const status = computed(() => {})

      return {
        status
      }
    }
  })
</script>
