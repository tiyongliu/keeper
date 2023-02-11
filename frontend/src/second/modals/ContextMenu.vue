<template>
  <ul
    class="dropDownMenuMarker" ref="wrapRef" :style="getStyle">
    <template v-for="item in preparedItems">
      <li v-if="item.divider" class="divider"/>
      <li v-else @mouseenter="e => handleMouseenter(e, item)">
        <a @click="handleClick($event, item)">
          {{ item.text || item.label }}
          <span v-if="item.keyText" class="keyText">{{ formatKeyText(item.keyText) }}</span>
          <div v-if="item.submenu" class="menu-right">
            <FontIcon icon="icon menu-right"/>
          </div>
        </a>
      </li>
    </template>
  </ul>
  <ContextMenu
    v-if="submenuItem && submenuItem?.submenu"
    :items="submenuItem?.submenu"
    v-bind="{...submenuOffset}"
    :closeParent="handleCloseParent"
  />
</template>
<script lang="ts">
import type {CSSProperties, PropType} from 'vue';
import {computed, defineComponent, nextTick, onMounted, onUnmounted, ref, toRefs, unref} from 'vue'
import {throttle} from "lodash-es";
import type {ContextMenuItem} from './typing';
import {prepareMenuItems} from '/@/second/utility/contextMenu'
import {formatKeyText} from '/@/second/utility/common'
import Icon from '/@/components/Icon';
import FontIcon from '/@/second/icons/FontIcon.vue'
import {useBootstrapStoreWithOut} from "/@/store/modules/bootstrap";

const props = {
  styles: {type: Object as PropType<CSSProperties>},
  left: {
    type: Number as PropType<number>,
    default: 0,
  },
  top: {
    type: Number as PropType<number>
  },
  items: {
    // The most important list, if not, will not be displayed
    type: [Array, Function] as PropType<ContextMenuItem[]>,
    default() {
      return [];
    },
  },
  closeParent: {
    type: Function as PropType<() => void>
  },
  targetElement: {
    type: [String, Object, Array]
  }
};

function getElementOffset(element, side: Nullable<string> = null) {
  var de = document.documentElement;
  var box = element.getBoundingClientRect();
  var top = box.top + window.pageYOffset - de.clientTop;
  var left = box.left + window.pageXOffset - de.clientLeft;
  if (side == 'right') return {top: top, left: left + box.width};
  return {top: top, left: left};
}

function fixPopupPlacement(element) {
  const {width, height} = element.getBoundingClientRect();
  let offset = getElementOffset(element);

  let newLeft: Nullable<number> = null;
  let newTop: Nullable<number> = null;

  if (offset.left + width > window.innerWidth) {
    newLeft = offset.left - width;

    if (newLeft < 0) newLeft = 0;
  }

  if (offset.top + height > window.innerHeight) {
    newTop = offset.top - height;

    if (newTop < 0) newTop = 0;
    if (newTop + height > window.innerHeight) {
      element.style.height = `${window.innerHeight - newTop}px`;
    }
  }

  if (newLeft != null) element.style.left = `${newLeft}px`;
  if (newTop != null) element.style.top = `${newTop}px`;
}

export default defineComponent({
  name: 'ContextMenu',
  components: {
    Icon,
    FontIcon
  },
  props,
  emits: ['close'],
  setup(props, {emit}) {
    const {targetElement, styles, closeParent, left, top} = props
    const {items} = toRefs(props)

    const wrapRef = ref<Nullable<HTMLElement>>(null)

    const hoverItem = ref<Nullable<ContextMenuItem>>(null)
    const hoverOffset = ref<Nullable<{ top: number, left: number }>>(null)

    const submenuItem = ref<Nullable<ContextMenuItem>>(null)
    const submenuOffset = ref<Nullable<{ top: number, left: number }>>(null)

    let closeHandlers: Function[] = []

    const showRef = ref(false);

    const getStyle = computed((): CSSProperties => {
      return {
        ...styles,
        left: `${left}px`,
        top: `${top}px`,
      }
    })

    onMounted(() => {
      nextTick(() => (showRef.value = true));
      fixPopupPlacement(wrapRef.value)
    });

    onUnmounted(() => {
      const el = unref(wrapRef);
      el && document.body.removeChild(el);
    });

    function dispatchClose() {
      emit('close')
      for (const handler of closeHandlers) {
        handler()
      }
      closeHandlers = []
    }

    function registerCloseHandler(handler) {
      closeHandlers.push(handler);
    }

    function handleClick(e, item) {
      if (item.disabled) return
      if (item.submenu) {
        hoverItem.value = item;
        hoverOffset.value = getElementOffset(e.target, 'right')

        submenuItem.value = item
        submenuOffset.value = hoverOffset.value
      }
      if (item.onClick) item.onClick()
    }

    function handleMouseenter(e, item) {
      hoverOffset.value = getElementOffset(e.target, 'right');
      hoverItem.value = item;
      changeActiveSubmenu();
    }

    function handleCloseParent() {
      if (closeParent) closeParent()
      dispatchClose()
    }

    const changeActiveSubmenu = throttle(() => {
      submenuItem.value = hoverItem.value;
      submenuOffset.value = hoverOffset.value;
    }, 500)

    const bootstrap = useBootstrapStoreWithOut()
    const preparedItems = computed<ContextMenuItem[]>(() => prepareMenuItems(items.value, {
      targetElement: targetElement,
      registerCloseHandler
    }, bootstrap))

    return {
      wrapRef,
      getStyle,
      preparedItems,
      formatKeyText,
      handleClick,
      handleMouseenter,
      submenuItem,
      submenuOffset,
      handleCloseParent,
    }
  },
});
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
