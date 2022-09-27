<template>
  <td
    :class="[
      isSelected && 'isSelected',
      isFrameSelected && 'isFrameSelected',
      isModifiedRow && 'isModifiedRow',
      isModifiedCell && 'isModifiedCell',
      isInserted && 'isInserted',
      isDeleted && 'isDeleted',
      isAutofillSelected && 'isAutofillSelected',
      isFocusedColumn && 'isFocusedColumn'
    ]"
    :data-row="rowIndex"
    :data-col="colIndex == null ? col.colIndex : colIndex">

  </td>
</template>

<script lang="ts">
import {defineComponent, PropType, toRefs} from 'vue'
import ShowFormButton from '/@/second/formview/ShowFormButton.vue'
export default defineComponent({
  name: "DataGridCell",
  components: {
    ShowFormButton
  },
  props: {
    rowIndex: {
      type: Number as PropType<number>,
    },
    col: {
      type: Object as PropType<{colIndex: number}>
    },
    colIndex: {
      type: Number as PropType<number>,
    },
    isSelected: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    isFrameSelected: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    isModifiedRow: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    isModifiedCell: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    isInserted: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    isDeleted: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    isAutofillSelected: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    isFocusedColumn: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    //domCell
    showSlot: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    //onSetFormView
    isDynamicStructure: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    isAutoFillMarker: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    isCurrentCell: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    // onDictionaryLookup
    // onSetValue
  },
  setup(props) {
    return {
      ...toRefs(props)
    }
  }
})
</script>

<style scoped>
td {
  font-weight: normal;
  border: 1px solid var(--theme-border);
  padding: 2px;
  white-space: nowrap;
  position: relative;
  overflow: hidden;
}
td.isFrameSelected {
  outline: 3px solid var(--theme-bg-selected);
  outline-offset: -3px;
}
td.isAutofillSelected {
  outline: 3px solid var(--theme-bg-selected);
  outline-offset: -3px;
}
td.isFocusedColumn {
  background: var(--theme-bg-alt);
}
td.isModifiedRow {
  background: var(--theme-bg-gold);
}
td.isModifiedCell {
  background: var(--theme-bg-orange);
}
td.isInserted {
  background: var(--theme-bg-green);
}
td.isDeleted {
  background: var(--theme-bg-volcano);
}
td.isSelected {
  background: var(--theme-bg-selected);
}
td.isDeleted {
  background-image: url('data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAEElEQVQImWNgIAX8x4KJBAD+agT8INXz9wAAAABJRU5ErkJggg==');
  background-repeat: repeat-x;
  background-position: 50% 50%;
}

.hint {
  color: var(--theme-font-3);
  margin-left: 5px;
}

.autoFillMarker {
  width: 8px;
  height: 8px;
  background: var(--theme-bg-selected-point);
  position: absolute;
  right: 0px;
  bottom: 0px;
  overflow: visible;
  cursor: crosshair;
}
</style>
