<template>
  <div class="collapseButtonMarker">
    <FontIcon :icon="collapsed ? 'icon triple-right' : 'icon triple-left'"/>
  </div>
</template>

<script lang="ts">
import {defineComponent, PropType, toRef} from 'vue'
import FontIcon from '/@/second/icons/FontIcon.vue'

export default defineComponent({
  name: "CollapseButton",
  components: {
    FontIcon
  },
  props: {
    collapsed: {
      type: Boolean as PropType<boolean>
    }
  },
  setup(props) {
    const collapsed = toRef(props, 'collapsed')
    return {
      collapsed,
    }
  }
})
</script>

<style scoped>
div {
  color: var(--theme-font-3);
  text-align: center;
  /* position: absolute;
  left: 0px;
  top: 4px;
  background-color: var(--theme-bg-1);
  border: 1px solid var(--theme-bg-1); */
}

div:hover {
  color: var(--theme-font-hover);
  border: 1px solid var(--theme-font-1);
}
</style>
