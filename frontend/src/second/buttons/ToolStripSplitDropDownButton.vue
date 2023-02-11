<template>
  <ToolStripButton :title="title" :icon="icon" :disabled="disabled" @splitclick="handleClick"/>
</template>

<script lang="ts" setup>
import {Component, defineProps, PropType, toRefs} from 'vue'
import ToolStripButton from './ToolStripButton.vue'
import {ContextMenuItem} from "/@/second/modals/typing";
import {useContextMenu} from "/@/hooks/web/useContextMenu";

const props = defineProps({
  menu: {
    type: Array as PropType<ContextMenuItem[]>,
  },
  disabled: {
    type: Boolean as PropType<boolean>,
    default: false
  },
  label: {
    type: String as PropType<Nullable<string>>,
  },
  icon: {
    type: String as PropType<Nullable<string>>,
  },
  component: {
    type: [String, Object] as PropType<string | Component>,
    default: ToolStripButton
  },
  title: {
    type: String as PropType<Nullable<string>>,
    default: undefined
  },
  splitIcon: {
    type: String as PropType<string>,
    default: 'icon chevron-down'
  }
})

const {menu, title, icon, disabled} = toRefs(props)
const [createContextMenu] = useContextMenu()

function handleClick(e) {
  e.preventDefault()
  e.stopPropagation()
  // @ts-ignore
  createContextMenu({event: e, items: menu})
}
</script>

<style scoped>

</style>
