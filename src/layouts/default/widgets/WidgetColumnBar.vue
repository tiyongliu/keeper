<template>
  <div class="main-container" :class="hidden ? 'hidden' : ''">
    <slot />
  </div>
</template>

<script lang="ts">
  import {defineComponent, computed, unref, ref, onMounted, provide, PropType} from 'vue';
  export default defineComponent({
    name: "DatabaseWidget",
    props: {
      hidden: {
        type: [Boolean] as PropType<boolean>,
        default: false,
      }
    },
    components: {

    },
    setup() {
      const hidden = ref(false)

      const definitions = ref([])
      const dynamicPropsCollection = []


      provide('pushWidgetItemDefinition', (item, dynamicProps) => {
        console.log('pushWidgetItemDefinition')

        dynamicPropsCollection.push(dynamicProps);
        definitions.value = [...definitions, item];
        return definitions.value.length - 1;
      })

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
