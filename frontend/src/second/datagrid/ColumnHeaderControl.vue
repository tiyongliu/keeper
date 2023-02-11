<template>
  <div class="columnHeader">
    <div class="label">
      <span v-if="grouping" class="grouping">
        {{ grouping == 'COUNT DISTINCT' ? 'distinct' : grouping.toLowerCase() }}
      </span>

      <ColumnLabel v-bind="{...column}"/>
      <span
        v-if="isStringDataType && !order"
        class="data-type"
        :title="`${column.dataType}`">{{ column.dataType.toLowerCase() }}</span>
    </div>
    <span v-if="order == 'ASC'" class="icon">
      <FontIcon icon="img sort-asc"/>
      <span v-if="orderIndex >= 0" class="color-icon-green order-index">{{ orderIndex + 1 }}}</span>
    </span>
    <span v-if="order == 'DESC'" class="icon">
      <FontIcon icon="img sort-desc"/>
      <span v-if="orderIndex >= 0" class="color-icon-green order-index">{{ orderIndex + 1 }}</span>
    </span>
    <DropDownButton :menu="getMenu" narrow/>
    <div
      class="horizontal-split-handle resizeHandleControl"
      v-splitterDrag="'clientX'"
      :resizeSplitter="(e) => handleResizeSplitter(e)">
    </div>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, PropType, toRefs} from 'vue'
import {isString} from 'lodash-es'
import {message} from 'ant-design-vue'
import ColumnLabel from '/@/second/elements/ColumnLabel.vue'
import FontIcon from '/@/second/icons/FontIcon.vue'
import DropDownButton from '/@/second/buttons/DropDownButton'
import {GroupFunc} from '/@/second/keeper-datalib'
import {openDatabaseObjectDetail} from '/@/second/appobj/DatabaseObjectAppObject'
import {copyTextToClipboard} from '/@/second/utility/clipboard'
import {isTypeDateTime} from '/@/second/keeper-tools'

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
        foreignKey?: any
        columnName: string
      }>
    },
    conid: {
      type: String as PropType<string>
    },
    database: {
      type: String as PropType<string>
    },
    grouping: {
      type: String as PropType<GroupFunc>
    },
    setSort: {
      type: Function as PropType<(order: string) => void>
    },
    addToSort: {
      type: Function as PropType<(order: any) => void>
    },
    order: {
      type: String as PropType<Partial<'ASC' | 'DESC'>>
    },
    orderIndex: {
      type: Number as PropType<number>,
      default: -1
    },
    isSortDefined: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    clearSort: {
      type: Function as PropType<() => void>
    },
    setGrouping: {
      type: Function as PropType<(groupFunc: any) => void>
    },
    allowDefineVirtualReferences: {
      type: Boolean as PropType<boolean>,
      default: false
    }
  },
  emits: ['resizeSplitter'],
  setup(props, {emit}) {
    const {
      column,
      setSort,
      setGrouping,
      isSortDefined,
      addToSort,
      order,
      clearSort,
      allowDefineVirtualReferences,
      conid,
      database
    } = toRefs(props)

    const openReferencedTable = () => {
      openDatabaseObjectDetail('TableDataTab', null, {
        schemaName: column.value!.foreignKey.refSchemaName,
        pureName: column.value!.foreignKey.refTableName,
        conid: conid.value,
        database: database.value,
        objectTypeField: 'tables',
      })
    }

    const handleDefineVirtualForeignKey = () => {
      message.warning('handleDefineVirtualForeignKey')
    }

    function getMenu() {
      return [
        setSort.value && {onClick: () => setSort.value!('ASC'), text: 'Sort ascending'},
        setSort.value && {onClick: () => setSort.value!('DESC'), text: 'Sort descending'},
        isSortDefined.value && addToSort.value && !order.value && {
          onClick: () => addToSort.value!('ASC'),
          text: 'Add to sort - ascending'
        },
        isSortDefined.value && addToSort.value && !order.value && {
          onClick: () => addToSort.value!('DESC'),
          text: 'Add to sort - descending'
        },
        order.value && clearSort.value && {
          onClick: () => clearSort.value!(),
          text: 'Clear sort criteria'
        },
        {onClick: () => copyTextToClipboard(column.value!.columnName), text: 'Copy column name'},

        (column.value && column.value?.foreignKey) && [{divider: true}, {
          onClick: openReferencedTable,
          text: column.value.foreignKey.refTableName
        }],

        setGrouping.value && {divider: true},
        setGrouping.value && {onClick: () => setGrouping.value!('GROUP'), text: 'Group by'},
        setGrouping.value && {onClick: () => setGrouping.value!('MAX'), text: 'MAX'},
        setGrouping.value && {onClick: () => setGrouping.value!('MIN'), text: 'MIN'},
        setGrouping.value && {onClick: () => setGrouping.value!('SUM'), text: 'SUM'},
        setGrouping.value && {onClick: () => setGrouping.value!('AVG'), text: 'AVG'},
        setGrouping.value && {onClick: () => setGrouping.value!('COUNT'), text: 'COUNT'},
        setGrouping.value && {
          onClick: () => setGrouping.value!('COUNT DISTINCT'),
          text: 'COUNT DISTINCT'
        },

        isTypeDateTime(column.value!.dataType) && [
          {divider: true},
          {onClick: () => setGrouping.value!('GROUP:YEAR'), text: 'Group by YEAR'},
          {onClick: () => setGrouping.value!('GROUP:MONTH'), text: 'Group by MONTH'},
          {onClick: () => setGrouping.value!('GROUP:DAY'), text: 'Group by DAY'},
        ],

        allowDefineVirtualReferences.value && [
          {divider: true},
          {onClick: handleDefineVirtualForeignKey, text: 'Define virtual foreign key'},
        ],
      ]
    }

    function handleResizeSplitter(e: Event) {
      emit('resizeSplitter', e)
    }

    const isStringDataType = computed(() => column.value && isString(column.value.dataType))

    return {
      ...toRefs(props),
      isStringDataType,
      getMenu,
      handleResizeSplitter
    }
  }
})
</script>

<style lang="less" scoped>
.columnHeader {
  display: flex;
  flex-wrap: nowrap;

.order-index {
  font-size: 10pt;
  margin-left: -3px;
  margin-right: 2px;
  top: -1px;
  position: relative;
}

.label {
  flex: 1;
  min-width: 10px;
  padding: 0 2px;
  margin: auto;
  white-space: nowrap;
}

.icon {
  margin-left: 3px;
  align-self: center;
  font-size: 18px;
}

/* .resizer {
  background-color: var(--theme-border);
  width: 2px;
  cursor: col-resize;
  z-index: 1;
} */

.grouping {
  color: var(--theme-font-alt);
  white-space: nowrap;
}

.data-type {
  color: var(--theme-font-3);
}

}
</style>
