<template>
  <template v-if="!skip && show">
    <WidgetTitle @click="visible = !visible">{{ title }}</WidgetTitle>
    <template v-if="visible">
      <div class="wrapper" :style="dynamicProps.splitterVisible ? `height:${size}px` : 'flex: 1 1 0'">
        <slot/>
      </div>
      <div
        v-if="dynamicProps.splitterVisible"
        class="vertical-split-handle"
        v-splitterDrag="'clientY'"
        :resizeSplitter="(e) => (size += e.detail)" />
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
  } from 'vue';
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
      const show = ref(true)
      const size = ref(0)
      const visible = ref(false)

      const {
        skip,
        height,
        collapsed,
        storageName,
      } = toRefs(props)

      const dynamicProps = reactive({splitterVisible: false})

      const pushWidgetItemDefinition = inject('pushWidgetItemDefinition') as any
      const updateWidgetItemDefinition = inject('updateWidgetItemDefinition') as any
      const widgetColumnBarHeight = inject('widgetColumnBarHeight') as any
      const widgetItemIndex = pushWidgetItemDefinition({
        collapsed,
        height,
        skip: skip || !show.value,
      }, dynamicProps)

      watch(
        () => [updateWidgetItemDefinition, unref(visible), unref(skip), unref(show)],
        ([_, watchVisible, watchSkip, watchShow]) => {
          updateWidgetItemDefinition(widgetItemIndex, {
            collapsed: !watchVisible,
            height,
            skip: watchSkip || !watchShow
          })
        },
      )

      watch(
        () => [unref(height), unref(widgetColumnBarHeight), unref(visible)],
        () => {
          setInitialSize(unref(height), unref(widgetColumnBarHeight))
        },
      )

      watch(
        () => [unref(widgetColumnBarHeight), unref(visible), unref(size)],
        ([watchHeight, watchVisible, watchSize]) => {
          if (unref(storageName) && unref(widgetColumnBarHeight) > 0) {
            setLocalStorage(unref(storageName), {
              relativeHeight: watchSize / watchHeight,
              visible: watchVisible
            })
          }
        }
      )

      function setInitialSize(initialSize, parentHeight) {
        if (unref(storageName)) {
          const storage = getLocalStorage(unref(storageName))
          if (storage) {
            size.value = parentHeight * storage.relativeHeight;
            return;
          }
        }
        if (isString(initialSize) && initialSize.endsWith('px')) size.value = parseInt(initialSize.slice(0, -2));
        else if (isString(initialSize) && initialSize.endsWith('%'))
          size.value = (parentHeight * parseFloat(initialSize.slice(0, -1))) / 100;
        else size.value = parentHeight / 3;
      }

      onMounted(() => {
        if (unref(storageName) && getLocalStorage(unref(storageName)) && getLocalStorage(unref(storageName)).visible != null) {
          visible.value = getLocalStorage(unref(storageName)).visible
        } else {
          visible.value = !collapsed.value
        }
      })

      return {
        ...toRefs(props),
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
