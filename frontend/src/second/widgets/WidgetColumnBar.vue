<template>
  <div ref="r" class="main-container" :class="hidden ? 'hidden' : ''">
    <slot></slot>
  </div>
</template>

<script lang="ts">
import {
  defineComponent,
  nextTick,
  PropType,
  provide,
  ref,
  toRef,
  unref,
  UnwrapNestedRefs,
  watch
} from 'vue'
import {useLocaleStore} from '/@/store/modules/locale'
import {storeToRefs} from "pinia"

export default defineComponent({
  name: "WidgetColumnBar",
  props: {
    hidden: {
      type: Boolean as PropType<boolean>,
      default: false,
    }
  },
  setup(props) {
    const r = ref<Nullable<HTMLElement>>(null)
    const widgetColumnBarHeight = ref(0)
    let definitions = ref<{ collapsed: boolean, height: number, skip: boolean }[]>([])
    let dynamicPropsCollection: UnwrapNestedRefs<{ splitterVisible: boolean }>[] = []

    const localeStore = useLocaleStore()
    const {selectedWidget} = storeToRefs(localeStore)

    provide('widgetColumnBarHeight', widgetColumnBarHeight)

    provide('pushWidgetItemDefinition', (item, dynamicProps) => {
      dynamicPropsCollection.push(dynamicProps)
      definitions.value = [...unref(definitions), item];
      return definitions.value.length - 1
    })

    provide('updateWidgetItemDefinition', (index, item) => {
      definitions.value = definitions.value.map((_, i) => i === index ? item : definitions.value[i])
    })

    function computeDynamicProps(defs: any[]) {
      for (let index = 0; index < defs.length; index++) {
        Object.assign(dynamicPropsCollection[index], {
          splitterVisible: !!defs.slice(index + 1).find(x => unref(x) && !unref(x.collapsed) && !unref(x.skip))
        })
      }
    }

    watch(() => [...definitions.value], () => {
      computeDynamicProps(definitions.value)
    })

    watch(() => [r.value, selectedWidget.value], async () => {
      await nextTick()
      widgetColumnBarHeight.value = r.value ? r.value!.clientHeight : 0
    })

    return {
      hidden: toRef(props, 'hidden'),
      r,
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
