<template>
  <div
    class="main"
    :class="isBold && 'isBold'"
    @click="handleClick"
    @mouseup="handleMouseUp"
    @contextmenu="handleContext">
<!--
@contextmenu="$event => handleContextMenu($event, disableContextMenu ? null : menu)">
-->
    <span v-if="expandIcon" class="expand-icon" @click.stop="handleExpand">
      <FontIcon :icon="expandIcon"/>
    </span>

    <span v-if="indentLevel" :style="{marginRight: `${indentLevel * 16}px`}"></span>

    <FontIcon v-if="isBusy" icon="icon loading" />
    <FontIcon v-else :icon="icon" />

    {{title}}

    <span v-if="statusIconBefore" class="status">
      <FontIcon :icon={statusIconBefore} />
    </span>

    <span v-if="statusIcon" class="status">
       <FontIcon :icon="statusIcon" :title="statusTitle" />
    </span>

    <span v-if="extInfo" class="ext-info">
      {{extInfo}}
    </span>
    <template v-if="onPin">
    <span class="pin" @click.stop.prevent="onPin">
      <FontIcon icon="mdi mdi-pin"/>
    </span>
    </template>
    <template v-if="onUnpin">
      <span class="pin-active" v-if="showPinnedInsteadOfUnpin">
        <FontIcon icon="icon pin" />
      </span>
      <template v-else>
        <span class="unpin" @click.stop.prevent="onUnpin">
          <FontIcon icon="icon close"/>
        </span>
      </template>

    </template>
  </div>
  <slot />
</template>

<script lang="ts">
import {defineComponent, PropType, toRefs} from 'vue'
import FontIcon from '/@/second/icons/FontIcon.vue'

//todo
import {useContextMenu} from '/@/hooks/web/useContextMenu'
import {handleContextMenu} from '/@/second/utility/contextMenu'
export default defineComponent({
  name: "AppObjectCore",
  props: {
    data: {
      type: Object as PropType<unknown>
    },
    icon: {
      type: String as PropType<string | null>,
    },
    isBold: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    isBusy: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    statusIcon: {
      type: String as PropType<string>,
      default: undefined
    },
    statusIconBefore: {
      type: String as PropType<string>,
    },
    statusTitle: {
      type: String as PropType<string>,
    },
    extInfo: {
      type: String as PropType<string | null>,
      default: undefined
    },
    filter: {
      type: String as PropType<string>,
    },
    expandIcon: {
      type: String as PropType<string>,
    },
    indentLevel: {
      type: Number as PropType<number>,
      default: 0
    },
    title: {
      type: String as PropType<string>,
    },
    colorMark: {

    },
    onPin: {
      type: Function as PropType<null | Function>
    },
    onUnpin: {
      type: Function as PropType<null | Function>
    },
    showPinnedInsteadOfUnpin: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    disableContextMenu: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    menu: {
      type: Function as PropType<null | Function>
    }
  },
  components: {
    FontIcon
  },
  emits: ['click', 'expand', 'dblclick'],
  setup(props, {emit}) {
    //todo
    let checkedObjectsStore = null

    const handleExpand = () => {
      //todo dispatch('expand');
      emit('expand')
    }

    const handleClick = () => {
      if (checkedObjectsStore) {

      } else {
        emit('click')
        emit('dblclick')
      }
    }

    const handleMouseUp = (e) => {
      e.preventDefault()
      e.stopPropagation()
    }

    const [createContextMenu] = useContextMenu()

    function handleContext(e: MouseEvent) {
      if (props.menu) {
        createContextMenu({
          event: e,
          items: props.menu(),
        });
      }

    }

    return {
      ...toRefs(props),
      handleExpand,
      handleClick,
      handleMouseUp,

      handleContext,
      handleContextMenu
    }
  }
})
</script>

<style scoped>
.main {
  padding: 5px;
  cursor: pointer;
  white-space: nowrap;
  font-weight: normal;
}
.main:hover {
  background-color: var(--theme-bg-hover);
}
.isBold {
  font-weight: bold;
}
.status {
  margin-left: 5px;
}
.ext-info {
  font-weight: normal;
  margin-left: 5px;
  color: var(--theme-font-3);
}
.expand-icon {
  margin-right: 3px;
}

.pin {
  float: right;
  color: var(--theme-font-2);
}
.pin:hover {
  color: var(--theme-font-hover);
}
.main .pin {
  visibility: hidden;
}
.main:hover .pin {
  visibility: visible;
}

.unpin {
  float: right;
  color: var(--theme-font-2);
}
.unpin:hover {
  color: var(--theme-font-hover);
}

.pin-active {
  float: right;
  color: var(--theme-font-2);
}
</style>
