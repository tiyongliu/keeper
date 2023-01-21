<template>
  <div class="group" @click="isExpanded = !isExpanded">
    <span class="expand-icon">
      <FontIcon :icon="plusExpandIcon(isExpanded)"/>
    </span>
    {{group}}
    {{items && `(${countText})`}}
  </div>

  <template v-if="isExpanded">
    <div class="ml-2" v-if="checkedObjectsStore">
      <Link @click="handleCheckAll(true)">Check all</Link>
      |
      <Link @click="handleCheckAll(false)">Uncheck all</Link>
    </div>
    <AppObjectListItem
      v-for="item in items"
      :isHidden="!item.isMatched"
      v-bind="$attrs"
      :module="module"
      :data="item.data"
      :disableContextMenu="disableContextMenu"
      :passProps="passProps"
    />
  </template>
</template>

<script lang="ts">
  import {Component, defineComponent, PropType, ref, toRef, toRefs, computed, unref} from 'vue'
  import AppObjectListItem from '/@/second/appobj/AppObjectListItem.vue'
  import FontIcon from '/@/second/icons/FontIcon.vue'
  import Link from '/@/second/elements/Link.vue'
  import {plusExpandIcon} from '/@/second/icons/expandIcons'
  export default defineComponent({
    name: 'AppObjectGroup',
    props: {
      group: {
        type: String as PropType<string>
      },
      groupFunc: {
        type: Function as PropType<Function>,
      },
      items: {
        type: Array  as PropType<any>
      },
      passProps: {
        type: Object as PropType<{ showPinnedInsteadOfUnpin: boolean }>,
      },
      module: {
        type: [String, Object] as PropType<string | Component>,
      },
      disableContextMenu: {
        type: Boolean as PropType<boolean>,
        default: false
      },
      checkedObjectsStore: { //todo 它是传递一个数据过来，我们只能传一个标识符过来。
        type: Object as PropType<unknown>,
      }
    },
    components: {
      AppObjectListItem,
      FontIcon,
      Link
    },
    setup(props) {
      const items = toRef(props, 'items')

      const isExpanded = ref(true)

      const filtered = computed(() => items.value!.filter(x => x.isMatched))

      const countText = computed(() => unref(filtered).length < unref(items)!.length ? `${unref(filtered).length}/${unref(items)!.length}` : `${unref(items)!.length}`)

      const handleCheckAll = (isChecked: boolean) => {
        console.log(isChecked)
      }

      return {
        ...toRefs(props),
        isExpanded,
        plusExpandIcon,
        countText,
        handleCheckAll,
      }
    }
  })
</script>

<style scoped>
  .group {
    user-select: none;
    padding: 3px 5px;
    cursor: pointer;
    white-space: nowrap;
    font-weight: bold;
  }

  .group:hover {
    background-color: var(--theme-bg-hover);
  }

  .expand-icon {
    margin-right: 3px;
  }
</style>
