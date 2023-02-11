<template>
  <div
    class="main"
    :class="isBold && 'isBold'"
    @click="handleClick"
    @mouseup="handleMouseUp"
    @contextmenu="handleContext">
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
    <template v-if="pin">
      <span class="pin" @click.stop.prevent="handlePin">
        <FontIcon icon="mdi mdi-pin"/>
      </span>
    </template>
    <template v-if="unpin">
      <span class="pin-active" v-if="showPinnedInsteadOfUnpin">
        <FontIcon icon="icon pin" />
      </span>
      <template v-else>
        <span class="unpin" @click.stop.prevent="handleUnpin">
          <FontIcon icon="icon close"/>
        </span>
      </template>

    </template>
  </div>
  <slot></slot>
</template>

<script lang="ts">
import {defineComponent, PropType, toRefs} from 'vue'
import FontIcon from '/@/second/icons/FontIcon.vue'
import {useContextMenu} from '/@/hooks/web/useContextMenu'

export default defineComponent({
  name: 'AppObjectCore',
  components: {
    FontIcon
  },
  props: {
    data: {
      type: Object as PropType<unknown>
    },
    icon: {
      type: String as PropType<Nullable<string>>,
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
    },
    statusIconBefore: {
      type: String as PropType<Nullable<string>>,
    },
    statusTitle: {
      type: String as PropType<string>,
    },
    extInfo: {
      type: String as PropType<Nullable<string>>,
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
    pin: {
      type: Function as PropType<Nullable<Function>>
    },
    unpin: {
      type: Function as PropType<Nullable<Function>>
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
  emits: ['click', 'expand', 'dblclick', 'middleclick'],
  setup(props, {emit}) {
    const {pin, unpin} = toRefs(props)
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

    function handlePin(e) {
      e.preventDefault()
      e.stopPropagation()
      pin.value && pin.value()
    }

    function handleUnpin(e) {
      e.preventDefault();
      e.stopPropagation();
      unpin.value && unpin.value()
    }

    return {
      ...toRefs(props),
      handleExpand,
      handleClick,
      handleMouseUp,
      handleContext,
      handlePin,
      handleUnpin
    }
  }
})
</script>

<style scoped>
.main {
  padding: 3px 5px;
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
