<template>
  <WidgetTitle v-if="!skip && show">{{ title }}</WidgetTitle>
  <div class="widgetColumnBarItem wrapper" v-if="visible"
       :style="dynamicProps.splitterVisible ? `height:${size}px` : 'flex: 1 1 0'">
    <slot/>
  </div>
  <div class="vertical-split-handle" v-splitterDrag="'clientY'"></div>
</template>

<script lang="ts">
  import {
    defineComponent,
    computed,
    inject,
    ref,
    onMounted,
    watch,
    toRefs,
    reactive,
    unref
  } from 'vue';
  import {isString} from 'lodash-es'
  import WidgetTitle from './WidgetTitle.vue'
  import {setLocalStorage, getLocalStorage} from '/@/second/utility/storageCache'

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
        type: [Boolean] as PropType<boolean>,
        default: false
      },
      show: {
        type: [Boolean] as PropType<boolean>,
        default: true
      },
      height: {
        type: [String] as PropType<string>,
      },
      collapsed: {
        type: [Boolean] as PropType<boolean>,
      },
      storageName: {
        type: [String] as PropType<string>,
      }
    },
    setup(props) {
      const skip = ref(false)
      const show = ref(true)
      const size = ref(0)
      const visible = ref(false)

      const {
        title,
        name,
        height,
        collapsed,
        storageName,
      } = unref(props) as unknown as {
        title: string,
        name: string,
        height: string,
        collapsed: boolean,
        storageName: string,
      }

      const dynamicProps = reactive({
        splitterVisible: false,
      })

      const pushWidgetItemDefinition = inject('pushWidgetItemDefinition') as any
      const updateWidgetItemDefinition = inject('updateWidgetItemDefinition') as any
      const widgetColumnBarHeight = inject('widgetColumnBarHeight') as any
      const widgetItemIndex = pushWidgetItemDefinition({
        collapsed,
        height,
        skip: skip || !show,
      }, dynamicProps)

      watch(
        () => updateWidgetItemDefinition,
        () => {
          updateWidgetItemDefinition(widgetItemIndex, {
            collapsed: !visible,
            height,
            skip: skip || !show
          })
        },
      )

      watch(
        () => [unref(height), unref(widgetColumnBarHeight)],
        () => {
          setInitialSize(unref(height), unref(widgetColumnBarHeight))
        },
      )

      watch(
        () => [unref(storageName), unref(widgetColumnBarHeight)],
        () => {
          if (unref(storageName) && unref(widgetColumnBarHeight) > 0) {
            setLocalStorage(storageName, {
              relativeHeight: unref(size) / unref(widgetColumnBarHeight),
              visible: visible.value
            })
          }
        }
      )

      function setInitialSize(initialSize, parentHeight) {
        if (storageName) {
          const storage = getLocalStorage(storageName)
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
        if (storageName && getLocalStorage(storageName) && getLocalStorage(storageName).visible != null) {
          visible.value = getLocalStorage(storageName).visible
        } else {
          visible.value = !props.collapsed
        }
      })

      //VUE3中watch和watchEffect的用法
      //svelte $ https://blog.csdn.net/qq_33325899/article/details/103554590

      return {
        skip,
        show,
        title: props.title,
        height: props.height,
        size,
        visible,

        dynamicProps,
      }
    }
  })
</script>

<style>
  .widgetColumnBarItem.wrapper {
    overflow: hidden;
    position: relative;
    flex-direction: column;
    display: flex;
  }
</style>
