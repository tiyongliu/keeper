<template>
  <div class="main" :class="isBold && 'isBold'" @click="handleClick">

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

    <span class="pin" v-if="onPin">
      <FontIcon icon="mdi mdi-pin" @click="handleOnPin"/>
    </span>
    <template v-if="onUnpin">
      <span class="pin-active" v-if="showPinnedInsteadOfUnpin">
        <FontIcon icon="icon pin" />
      </span>
      <span class="unpin" v-else>
        <FontIcon icon="icon close" @click="handleOnUnpin"/>
      </span>
    </template>
  </div>
  <slot />
</template>

<script lang="ts">
import {defineComponent, PropType, Ref, unref, onMounted} from 'vue'
import FontIcon from '/@/second/icons/FontIcon.vue'
import {IPinnedDatabasesItem} from '/@/second/types/standard.d'
export default defineComponent({
  name: "AppObjectCore",
  props: {
    data: {
      type: Object as PropType<IPinnedDatabasesItem>
    },
    icon: {
      type: String as PropType<string>,
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
      type: [String, Object] as PropType<string>,
    },
    extInfo: {
      type: [String, Object] as PropType<string | Ref<string>>,
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
  },
  components: {
    FontIcon
  },
  emits: ['click', 'expand'],
  setup(props, {emit}) {
    const {onPin, onUnpin} = props
    const handleExpand = () => {
      //todo dispatch('expand');
      emit('expand')
    }

    const handleOnPin = (e) => {
      e.preventDefault()
      e.stopPropagation()
      onPin && onPin()
    }

    const handleOnUnpin = (e) => {
      e.preventDefault()
      e.stopPropagation()
      onUnpin && onUnpin()
    }

    const handleClick = () => {
      //todo if (checkedObjectsStore) {
      emit('click')
    }

    onMounted(() => {
    })

    return {
      ...props,
      handleExpand,
      handleOnPin,
      handleOnUnpin,
      handleClick
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
