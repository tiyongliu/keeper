<template>
  <div class="main" :class="isBold && 'isBold'">

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

    <span class="pin">
      <FontIcon icon="mdi mdi-pin"/>
    </span>

  </div>
  <slot />
</template>

<script lang="ts">
import {defineComponent, PropType} from 'vue'
import FontIcon from '/@/second/icons/FontIcon.vue'
export default defineComponent({
  name: "AppObjectCore",
  props: {
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
      default: undefined
    },
    statusTitle: {
      type: String as PropType<string>,
      default: undefined
    },
    extInfo: {
      type: String as PropType<string>,
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

    }
  },
  components: {
    FontIcon
  },
  setup(props) {
    const handleExpand = () => {
      //todo dispatch('expand');
    }

    return {
      ...props,
      handleExpand,
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
