<template>
  <WidgetTitle v-if="!skip && show">{{ title }}</WidgetTitle>
  <div class="wrapper"
       v-if="visible"
       :style="dynamicProps.splitterVisible ? `height:${size}px` : 'flex: 1 1 0'">
    <slot/>
  </div>
  <div
    class="vertical-split-handle"
    v-splitterDrag="'clientY'"

  />
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
  import {setLocalStorage, getLocalStorage} from '/@/utils/utility/storageCache'

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
        title: Ref<string>,
        name: Ref<string>,
        height: Ref<string>,
        collapsed: Ref<boolean>,
        storageName: Ref<string>,
      }

      const dynamicProps = reactive({
        splitterVisible: false,
      })

      const pushWidgetItemDefinition = inject('pushWidgetItemDefinition')
      //todo 需要动态
      const updateWidgetItemDefinition = inject('updateWidgetItemDefinition')
      const widgetColumnBarHeight = inject('widgetColumnBarHeight')
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

      if (storageName && unref(widgetColumnBarHeight) > 0) {
        setLocalStorage(storageName, {
          relativeHeight: size.value / widgetColumnBarHeight,
          visible: visible.value
        })
      }

      function setInitialSize(initialSize, parentHeight) {
        console.log(initialSize, parentHeight, '00')
        // size.value = 105
        console.log(storageName)
        if (storageName) {
          const storage = getLocalStorage(storageName)
          if (storage) {
            size.value = parentHeight * storage.relativeHeight;
            console.log(size.value, '1')
            return;
          }
        }
        if (isString(initialSize) && initialSize.endsWith('px')) {
          size.value = parseInt(initialSize.slice(0, -2));
          console.log(size.value, '2')
        }
        else if (isString(initialSize) && initialSize.endsWith('%')) {
          size.value = (parentHeight * parseFloat(initialSize.slice(0, -1))) / 100;
          console.log(size.value, '3')
        }
        else {
          size.value = parentHeight / 3;
          console.log(size.value, '4')
        }
      }

      onMounted(() => {
        if (storageName && getLocalStorage(storageName) && getLocalStorage(storageName).visible != null) {
          visible.value = getLocalStorage(storageName).visible
        } else {
          visible.value = !props.collapsed
        }

        setInitialSize(unref(height), unref(widgetColumnBarHeight))
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

<style scoped>
  .wrapper {
    overflow: hidden;
    position: relative;
    flex-direction: column;
    display: flex;
  }
</style>
