<template>
  <div ref="r" class="main-container" :class="hidden ? 'hidden' : ''">
    <slot/>
  </div>
</template>

<script lang="ts">
  import {defineComponent, computed, unref, ref, onMounted, provide, PropType, watch} from 'vue';

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
        dynamicPropsCollection.value.push(dynamicProps)
        definitions.value = [...definitions.value, item];
        return definitions.value.length - 1
      })

      provide('updateWidgetItemDefinition', (index, item) => {
        console.log(`更新了`, index, item)
        definitions.value[index] = item
      })

      watch(
        () => unref(definitions),
        (defs) => computeDynamicProps(defs),
      );

      function computeDynamicProps(defs: any[]) {
        for (let index = 0; index < defs.length; index++) {
          // const definition = defs[index];
          dynamicPropsCollection.value[index].splitterVisible = !!defs.slice(index + 1).find(x => x && !x.collapsed && !x.skip);
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
