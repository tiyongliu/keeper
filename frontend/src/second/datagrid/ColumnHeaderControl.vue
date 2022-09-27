<template>
  <div class="header">
    <div class="label">
      <span v-if="grouping" class="grouping">
        {{ grouping == 'COUNT DISTINCT' ? 'distinct' : grouping.toLowerCase() }}
      </span>

      <ColumnLabel v-bind="{...column}"/>
      <span v-if="isString(column.dataType) && !order" class="data-type"
            :title="`${column.dataType}`">{{ column.dataType.toLowerCase() }}</span>
    </div>

    <span v-if="order == 'ASC'" class="icon">
      <FontIcon icon="img sort-asc"/>
      <span v-if="orderIndex >= 0" class="color-icon-green order-index">{{ orderIndex + 1 }}}</span>
    </span>

    <span class="icon">
      <FontIcon icon="img sort-desc"/>
      <span v-if="orderIndex >= 0" class="color-icon-green order-index">{{ orderIndex + 1 }}</span>
    </span>
    <!--    <DropDownButton narrow />-->
    <!--    <div-->
    <!--      class="horizontal-split-handle resizeHandleControl"-->
    <!--      v-splitterDrag="'clientX'"-->
    <!--      :resizeSplitter="(e) => handleResizeSplitter(e)"></div>-->
  </div>
</template>

<script lang="ts">
import {defineComponent, PropType, readonly} from 'vue'
import {isString} from 'lodash-es'
import ColumnLabel from '/@/second/elements/ColumnLabel.vue'
import FontIcon from '/@/second/icons/FontIcon.vue'
import DropDownButton from '/@/second/buttons/DropDownButton'

export default defineComponent({
  name: "ColumnHeaderControl",
  components: {
    ColumnLabel,
    FontIcon,
    DropDownButton,
  },
  props: {
    column: {
      type: Object as PropType<{
        dataType: string
      }>
    },
    grouping: {
      type: String as PropType<string>
    },
    order: {
      type: [Object, String] as PropType<object | string>
    },
    orderIndex: {
      type: Number as PropType<number>
    }
  },
  emits: ['resizeSplitter'],
  setup(props, {emit}) {
    const column = readonly(props.column!)

    function getMenu() {
      return []
    }

    function handleResizeSplitter(e: Event) {
      emit('resizeSplitter', e)
    }

    return {
      ...props,
      column,
      isString,
      getMenu,
      handleResizeSplitter
    }
  }
})
</script>

<style scoped>

</style>
