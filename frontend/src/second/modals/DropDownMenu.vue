<template>
  <ul class="dropDownMenuMarker">
    <template v-for="(item, index) in items" :key="index">
      <li v-if="item.divider" class="divider" />
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
import {defineExpose, defineProps, onBeforeUnmount, onMounted, PropType, ref} from 'vue'
import {formatKeyText} from '/@/second/utility/common'
import {fixPopupPlacement} from './DropDownMenu_'
import FontIcon from '/@/second/icons/FontIcon.vue'

defineProps({
  top: {
    type: [String, Number] as PropType<string | number>,
  },
  left: {
    type: [String, Number] as PropType<string | number>,
  },
  items: {
    type: Array
  }
})

defineExpose({
  formatKeyText
})

const handleClick = (e, item) => {
  if (item.disabled) return
  if (item.submenu) {

    return
  }
}

const element = ref<null | HTMLElement>(null)

onMounted(() => {
  fixPopupPlacement(element.value!)
  document.addEventListener('mousedown', handleClickOutside, true)
})

onBeforeUnmount(() => document.removeEventListener('mousedown', handleClickOutside, true))

const changeActiveSubmenu = throttle(() => {

}, 500)

const handleClickOutside = event => {
  if (event.target.closest('ul.dropDownMenuMarker')) return

  // dispatch('close')
}
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
