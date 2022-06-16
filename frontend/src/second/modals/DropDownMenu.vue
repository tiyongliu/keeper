<template>
  <ul class="dropDownMenuMarker" :style="style" ref="element">
    <template v-for="(item, index) in preparedItems" :key="index">
      <li v-if="item.divider" class="divider"/>
      <li v-else>
        <a @click="$event => handleClick($event, item)">
          {{ item.text || item.label }}
          <span v-if="item.keyText" class="keyText">{{ formatKeyText(item.keyText) }}</span>
          <div v-if="item.submenu" class="menu-right">
            <FontIcon icon="icon menu-right"/>
          </div>
        </a>
      </li>
    </template>
  </ul>
</template>

<script lang="ts" setup>
import {throttle} from 'lodash-es'
import {
  defineEmits,
  defineExpose,
  defineProps,
  watch,
  onBeforeUnmount,
  onMounted,
  computed,
  PropType,
  ref,
  toRefs
} from 'vue'
import { prepareMenuItems } from '/@/second/utility/contextMenu'
import {formatKeyText} from '/@/second/utility/common'
import FontIcon from '/@/second/icons/FontIcon.vue'
import { dataBaseStore } from "/@/store/modules/dataBase"
import {fixPopupPlacement} from './DropDownMenu_'

const props = defineProps({
  top: {
    type: [String, Number] as PropType<string | number>,
  },
  left: {
    type: [String, Number] as PropType<string | number>,
  },
  items: {
    type: Array
  },
  targetElement: {
    type: Object
  }
})

const {left, top, items, targetElement} = toRefs(props)
const style = `left: ${left?.value}px; top: ${top?.value}px`

defineExpose({
  formatKeyText
})

const handleClick = (e, item) => {
  if (item.disabled) return
  if (item.submenu) {
    return
  }

  emit('close')
}

const emit = defineEmits(['close'])

const element = ref<null | HTMLElement>(null)

const dataBase = dataBaseStore()

onMounted(() => {
  fixPopupPlacement(element.value!)
  document.addEventListener('mousedown', handleClickOutside, true)
})

// const preparedItems = computed(() => prepareMenuItems(items?.value, {targetElement: targetElement?.value}, dataBase.$state.commandsCustomized))
const preparedItems = computed(() => [{"text":"Edit"},{"text":"Delete"},{"text":"Duplicate"},{"text":"Connect"},{"text":"New query","isNewQuery":true},{"text":"Restore/import SQL dump"}])

onBeforeUnmount(() => document.removeEventListener('mousedown', handleClickOutside, true))

const changeActiveSubmenu = throttle(() => {

}, 500)

const handleClickOutside = event => {
  if (event.target.closest('ul.dropDownMenuMarker')) return

  emit('close')
}

watch(() => preparedItems, () => {
  console.log(preparedItems.value, `------------------------------preparedItems`)
})
</script>


<style scoped>
ul {
  position: absolute;
  list-style: none;
  background-color: var(--theme-bg-0);
  border-radius: 4px;
  border: 1px solid var(--theme-border);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.175);
  padding: 5px 0;
  margin: 2px 0 0;
  font-size: 14px;
  text-align: left;
  min-width: 160px;
  z-index: 1050;
  cursor: default;
  white-space: nowrap;
  overflow-y: auto;
}

.keyText {
  font-style: italic;
  font-weight: bold;
  text-align: right;
  margin-left: 16px;
}

a {
  padding: 3px 20px;
  line-height: 1.42;
  white-space: nop-wrap;
  color: var(--theme-font-1);
  display: flex;
  justify-content: space-between;
}

a.disabled {
  color: var(--theme-font-3);
}

a:hover:not(.disabled) {
  background-color: var(--theme-bg-1);
  text-decoration: none;
  color: var(--theme-font-1);
}

.divider {
  margin: 9px 0px 9px 0px;
  border-top: 1px solid var(--theme-border);
  border-bottom: 1px solid var(--theme-bg-0);
}

.menu-right {
  position: relative;
  left: 15px;
}
</style>
