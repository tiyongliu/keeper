<template>
  <template v-if="!skip && show">
    <WidgetTitle @click="visible = !visible">{{ title }}</WidgetTitle>
    <template v-if="visible">
      <div class="wrapper" :style="dynamicProps.splitterVisible ? `height:${size}px` : 'flex: 1 1 0'">
        <slot></slot>
      </div>
      <div
        v-if="dynamicProps.splitterVisible"
        class="vertical-split-handle"
        v-splitterDrag="'clientY'"
        :resizeSplitter="(e) => (size += e.detail)">
      </div>
    </template>
  </template>
</template>

<script lang="ts">
import {
  defineComponent,
  inject,
  ref,
  onMounted,
  watch,
  toRefs,
  reactive,
  unref,
  Ref
} from 'vue'
import {isString} from 'lodash-es'
import {setLocalStorage, getLocalStorage} from '/@/second/utility/storageCache'
import WidgetTitle from './WidgetTitle.vue'

export default defineComponent({
  name: "WidgetColumnBarItem",
  components: {
    WidgetTitle
  },
  props: {
    title: {
      type: String as PropType<string>,
    },
    name: {
      type: String as PropType<string>,
    },
    skip: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    show: {
      type: Boolean as PropType<boolean>,
      default: true
    },
    height: {
      type: String as PropType<string>,
    },
    collapsed: {
      type: Boolean as PropType<boolean>,
    },
    storageName: {
      type: String as PropType<string>,
    }
  },
  setup(props) {
    const {
      skip,
      height,
      collapsed,
      storageName,
      show,
    } = toRefs(props)

    const size = ref(0)
    const visible = ref(false)

    const dynamicProps = reactive({splitterVisible: false})
    const pushWidgetItemDefinition = inject('pushWidgetItemDefinition') as any
    const updateWidgetItemDefinition = inject('updateWidgetItemDefinition') as any
    const widgetColumnBarHeight = inject('widgetColumnBarHeight') as Ref<number>
    const widgetItemIndex = pushWidgetItemDefinition({
      collapsed: collapsed.value,
      height: height.value,
      skip: skip.value || !show.value,
    }, dynamicProps)

    watch(
      () => [updateWidgetItemDefinition, visible.value, skip.value, show.value],
      () => {
        updateWidgetItemDefinition(widgetItemIndex, {
          collapsed: !visible.value,
          height: height.value,
          skip: skip.value || !show.value
        })
      }
    )

    watch(
      () => unref(widgetColumnBarHeight),
      () => setInitialSize(height.value, unref(widgetColumnBarHeight)),
    )

    watch(
      () => [visible.value, size.value],
      () => {
        if (storageName.value && unref(widgetColumnBarHeight) > 0) {
          setLocalStorage(storageName.value, {
            relativeHeight: size.value / widgetColumnBarHeight.value,
            visible: visible.value
          })
        }
      }
    )

    function setInitialSize(initialSize, parentHeight) {
      if (storageName.value) {
        const storage = getLocalStorage(storageName.value)
        if (storage) {
          size.value = parentHeight * storage.relativeHeight
          return;
        }
      }
      if (isString(initialSize) && initialSize.endsWith('px')) size.value = parseInt(initialSize.slice(0, -2))
      else if (isString(initialSize) && initialSize.endsWith('%'))
        size.value = (parentHeight * parseFloat(initialSize.slice(0, -1))) / 100
      else size.value = parentHeight / 3
    }

    onMounted(() => {
      if (storageName.value && getLocalStorage(storageName.value) && getLocalStorage(storageName.value).visible != null) {
        visible.value = getLocalStorage(storageName.value).visible
      } else {
        visible.value = !collapsed.value
      }
    })

    return {
      title: props.title,
      name: props.name,
      skip,
      height,
      collapsed,
      storageName,
      show,
      size,
      visible,
      dynamicProps,
    }
  }
})
</script>

<style scoped>
.wrapper {
  overflow: hidden;
  position: relative;
  flex-direction: column;
  display: flex;
}
</style>
