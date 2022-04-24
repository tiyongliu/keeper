<template>
  <div ref="r" class="main-container" :class="hidden ? 'hidden' : ''">
    <slot/>
  </div>
</template>

<script lang="ts">
  import {defineComponent, computed, unref, ref, onMounted, provide, PropType, watch, toRefs } from 'vue';
  export default defineComponent({
    name: "DatabaseWidget",
    props: {
      hidden: {
        type: [Boolean] as PropType<boolean>,
        default: false,
      }
    },
    components: {},
    setup() {
      const r = ref()
      const hidden = ref(false)
      const widgetColumnBarHeight = ref(0)

      const definitions = ref<{ collapsed: boolean, height: number, skip: boolean }[]>([])
      const dynamicPropsCollection = ref<{ splitterVisible: boolean }[]>([])

      provide('widgetColumnBarHeight', widgetColumnBarHeight)
      provide('pushWidgetItemDefinition', (item, dynamicProps) => {
        console.log(`update pushWidgetItemDefinition line 28 前`)

        dynamicPropsCollection.value.push(dynamicProps)
        definitions.value = [...unref(definitions), item];
        return definitions.value.length - 1
      })

      provide('updateWidgetItemDefinition', (index, item) => {
        console.log(`update updateWidgetItemDefinition line 36 后`, index, item)
        definitions.value[index] = item
      })

      watch(
        () => unref(definitions),
        (defs) => computeDynamicProps(defs),
      );

      function computeDynamicProps(defs: any[]) {
        for (let index = 0; index < defs.length; index++) {
          const definition = defs[index];
          const splitterVisible = !!defs.slice(index + 1).find(x => unref(x) && !unref(x.collapsed) && !unref(x.skip));
          console.log(splitterVisible)
          dynamicPropsCollection.value[index].splitterVisible = splitterVisible
        }
      }

      onMounted(() => {
        widgetColumnBarHeight.value = r.value.clientHeight
      })

      return {
        hidden,
        clientHeight: widgetColumnBarHeight,
        r
      }
    }
  })
</script>

<style lang="less">
  .main-container {
    position: relative;
    flex: 1;
    flex-direction: column;
    user-select: none;
    &.hidden {
      display: none;
    }

    &:not(.hidden) {
      display: flex;
    }
  }
</style>
