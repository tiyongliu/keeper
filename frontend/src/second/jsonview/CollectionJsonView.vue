<template>
  <div class="wrapper">
    <div class="toolbar">
      <Pager :skip="skip" :limit="limit"/>
    </div>
    <div class="json"></div>
  </div>
</template>

<script lang="ts">
import {defineComponent, ref} from 'vue'
import Pager from '/@/second/elements/Pager.vue'
export default defineComponent({
  name: "CollectionJsonView",
  components: {
    Pager
  },
  setup() {
    const skip = ref(0)
    const limit = ref(50)
    return {
      skip,
      limit
    }
  }
})
</script>

<style scoped>
.toolbar {
  background: var(--theme-bg-1);
  display: flex;
  border-bottom: 1px solid var(--theme-border);
  border-top: 2px solid var(--theme-border);
  margin-bottom: 3px;
}

.json {
  overflow: auto;
  flex: 1;
  /* position: relative; */
}

.wrapper {
  display: flex;
  flex-direction: column;
  position: absolute;
  left: 0;
  top: 0;
  right: 0;
  bottom: 0;
}
</style>
