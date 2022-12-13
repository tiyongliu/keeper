<template>
  <div class="row" v-if="column">
    <div class="name">{{ column.columnName }}</div>
    <div class="nowrap">
      <span class="icon">
        <FontIcon icon="icon edit"/>
      </span>
      <span class="icon">
      <FontIcon icon="icon delete" />
      </span>
        <span class="icon">
        <FontIcon icon="icon arrow-up" />
      </span>
        <span class="icon">
        <FontIcon icon="icon arrow-down" />
      </span>
    </div>
  </div>
</template>

<script lang="ts">
import {defineComponent, toRef} from 'vue'
import {DisplayColumn} from '/@/second/keeper-datalib'
import FontIcon from '/@/second/icons/FontIcon.vue'

export default defineComponent({
  name: 'ColumnManagerRow',
  components: {FontIcon},
  props: {
    column: {
      type: Object as PropType<DisplayColumn>
    },
  },
  setup(props) {
    return {
      column: toRef(props, 'column')
    }
  }
})
</script>

<style scoped>
.row {
  display: flex;
  justify-content: space-between;
  cursor: pointer;
}

.row:hover {
  background-color: var(--theme-bg-selected);
}

.name {
  white-space: nowrap;
  margin: 5px;
}

.icon {
  position: relative;
  top: 5px;
  padding: 5px;
}

.icon:hover {
  background-color: var(--theme-bg-3);
}
</style>
