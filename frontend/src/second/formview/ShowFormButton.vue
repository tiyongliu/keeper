<template>
  <div class="showFormButtonMarker">
    <FontIcon :icon="icon" />
  </div>
</template>

<script lang="ts">
import {defineComponent, PropType} from 'vue'
import FontIcon from '/@/second/icons/FontIcon.vue'
export default defineComponent({
  name: "ShowFormButton",
  components: {
    FontIcon
  },
  props: {
    icon: {
      type: String as PropType<string>,
      default: 'icon form'
    },
  },
  setup(props) {
    return {
      ...props
    }
  }
})
</script>

<style scoped>
div {
  position: absolute;
  right: 0px;
  top: 1px;
  color: var(--theme-font-3);
  background-color: var(--theme-bg-1);
  border: 1px solid var(--theme-bg-1);
}

div:hover {
  color: var(--theme-font-hover);
  border: var(--theme-border);
  top: 1px;
  right: 0px;
}
</style>
