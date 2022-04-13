<template>
  <div class="main isBold">
    <FontIcon icon="mdi mdi-plus-box-outline"/>
    <FontIcon icon="mdi mdi-minus-box-outline"/>

    <FontIcon icon="mdi mdi-server color-icon-blue"/>
    localhost
    <FontIcon icon="mdi mdi-check-circle color-icon-green"/>
    <span class="ext-info">mysql</span>
  </div>

  <div class="subitems">
    <div class="main" draggable="true">
      <FontIcon icon="mdi mdi-database color-icon-gold" />crmeb_java_beta
      <span class="pin">
         <FontIcon icon="mdi mdi-pin"/>
      </span>
    </div>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, PropType, unref} from 'vue';
import FontIcon from '/@/second/icons/FontIcon.vue'
export default defineComponent({
  name: "AppObjectListItem",
  props: {
    data: {
      type: Array as unknown as PropType<[]>,
      default: [],
    },
    isHidden: {
      type: Boolean as unknown as PropType<boolean>,
    },
    isExpandedBySearch: {
      type: Boolean as unknown as PropType<boolean>,
      default: false
    },
    expandOnClick: {
      type: Boolean as PropType<boolean>,
    },
  },
  components: {
    FontIcon
  },
  setup(props) {
    // const dynamicList = computed(() => unref(props.list))
    console.log(props)

    async function handleExpand() {
      // if (subItemsComponent && expandOnClick) {
      //   await tick();
      //   isExpanded = !isExpanded;
      // }
    }

    function handleExpandButton() {

    }

    return {
      ...props,
      handleExpand,
      handleExpandButton,
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

<style>
.subitems {
  margin-left: 28px;
}
</style>
