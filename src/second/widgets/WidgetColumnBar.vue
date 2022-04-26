<template>
  <div ref="r" class="main-container" :class="hidden ? 'hidden' : ''">
    <slot/>
  </div>
</template>

<script lang="ts">
  import {defineComponent, unref, ref, onMounted, provide, PropType, reactive } from 'vue';
  export default defineComponent({
    name: "WidgetColumnBar",
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

      let definitions = reactive<{ collapsed: boolean, height: number, skip: boolean }[]>([])
      let dynamicPropsCollection = reactive<{ splitterVisible: boolean }[]>([])

      provide('widgetColumnBarHeight', widgetColumnBarHeight)
      provide('pushWidgetItemDefinition', (item, dynamicProps) => {
        dynamicPropsCollection.push(dynamicProps)
        definitions = [...unref(definitions), item];
        return definitions.length - 1
      })

      provide('updateWidgetItemDefinition', (index, item) => {
        definitions[index] = item
        computeDynamicProps(definitions)
      })

      function computeDynamicProps(defs: any[]) {
        for (let index = 0; index < defs.length; index++) {
          dynamicPropsCollection[index].splitterVisible = !!defs.slice(index + 1).
          find(x => unref(x) && !unref(x.collapsed) && !unref(x.skip));
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
