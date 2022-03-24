<template>
  <WidgetTitle v-if="!skip && show">{{ title }}</WidgetTitle>
  <div class="wrapper" v-if="visible">
    <slot/>
  </div>
</template>

<script lang="ts">
import {defineComponent, computed, inject, ref, onMounted, watch} from 'vue';
import WidgetTitle from './WidgetTitle.vue'
import {getLocalStorage} from '/@/utils/utility/storageCache'
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

    const pushWidgetItemDefinition = inject('pushWidgetItemDefinition')

    console.log(pushWidgetItemDefinition, 'pushWidgetItemDefinition')

    // const widgetColumnBarHeight = inject('widgetColumnBarHeight')
    //临时
    const widgetColumnBarHeight = 0


    if (props.storageName && widgetColumnBarHeight > 0) {
      // setLocalStorage(storageName, {relativeHeight: })
    }

    function setInitialSize(initialSize, parentHeight) {
      if (props.storageName) {
        const storage = getLocalStorage(props.storageName)
        if (storage) {
          size.value = parentHeight * storage.relativeHeight;
          return;
        }
      }
      if (_.isString(initialSize) && initialSize.endsWith('px')) size.value = parseInt(initialSize.slice(0, -2));
      else if (_.isString(initialSize) && initialSize.endsWith('%'))
        size.value = (parentHeight * parseFloat(initialSize.slice(0, -1))) / 100;
      else size.value = parentHeight / 3;
    }

    computed(() => setInitialSize)
    onMounted(() => {
      if (props.storageName && getLocalStorage(props.storageName) && getLocalStorage(props.storageName).visible != null) {
        visible.value = getLocalStorage(props.storageName).visible
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
