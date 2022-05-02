<template>
  <div class="group" @click="isExpanded = !isExpanded">
    <span class="expand-icon">
      <FontIcon :icon="plusExpandIcon(isExpanded)"/>
    </span>
    {{group}}
    {{items && `(${countText})`}}
  </div>

  <template v-if="isExpanded">
    <!--<div class="ml-2">
      <Link>Check all</Link>
      |
      <Link>Uncheck all</Link>
    </div>-->
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
  import {Component, defineComponent, PropType, ref, toRefs, computed, unref} from 'vue';
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
        type: Array  as PropType<{isMatched?: string}[]>
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
      }
    },
    components: {
      AppObjectListItem,
      FontIcon,
      Link
    },
    setup(props) {
      const isExpanded = ref(true)

      const {items} = toRefs(props)

      const filtered = computed(() => items.value!.filter(x => x.isMatched))

      const countText = computed(() => unref(filtered).length < unref(items)!.length ? `${unref(filtered).length}/${unref(items)!.length}` : `${unref(items)!.length}`)

      return {
        ...toRefs(props),
        isExpanded,
        plusExpandIcon,
        countText,
      }
    }
  })
</script>

<style scoped>
  .group {
    user-select: none;
    padding: 5px;
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
