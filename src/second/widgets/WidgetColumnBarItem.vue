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
    computed,
    inject,
    ref,
    onMounted,
    watch,
    toRefs,
    unref
  } from 'vue';
  import {isString} from 'lodash-es'
  import WidgetTitle from './WidgetTitle.vue'
  import {setLocalStorage, getLocalStorage} from '/@/second/utility/storageCache'

  import { cssVariableStore } from "/@/store/modules/cssVariable";

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
      const cssVariable = cssVariableStore()
      const show = ref(true)
      const size = ref(0)
      const visible = ref(false)

      const {
        skip,
        height,
        collapsed,
        storageName,
      } = toRefs(props)

      const dynamicProps = computed(() => cssVariable.getDynamicProps)

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
