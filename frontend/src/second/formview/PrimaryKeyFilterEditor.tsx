import {computed, defineComponent, PropType, ref, toRefs, watch} from 'vue'
import {isNumber} from 'lodash-es'
import {TableFormViewDisplay} from '/@/second/keeper-datalib'
import FontIcon from '/@/second/icons/FontIcon.vue'
import ColumnLabel from '/@/second/elements/ColumnLabel.vue'
import keycodes from '/@/second/utility/keycodes'
import {ColumnReference} from '/@/second/keeper-types'
import {DeleteOutlined} from '@ant-design/icons-vue'

export default defineComponent({
  name: 'PrimaryKeyFilterEditor',
  props: {
    column: {
      type: Object as PropType<ColumnReference>,
    },
    baseTable: {
      type: Object as PropType<{
        columns: ColumnReference[]
      }>,
    },
    formDisplay: {
      type: Object as PropType<TableFormViewDisplay>
    },
  },
  setup(props) {
    const {column, baseTable, formDisplay} = toRefs(props)
    const domEditor = ref<Nullable<string>>(null)

    const value = computed(() => formDisplay.value && column.value
      ? formDisplay.value.getKeyValue(column.value.columnName) : null)

    const applyFilter = () => {
      (formDisplay.value && column.value) && formDisplay.value.requestKeyValue(column.value.columnName, isNumber(value.value)
        ? Number(domEditor.value) : domEditor.value)
    }

    const cancelFilter = () => {
      formDisplay.value && formDisplay.value.cancelRequestKey(null)
      formDisplay.value && formDisplay.value.reload()
    }

    const handleKeyDown = ev => {
      if (ev.keyCode == keycodes.enter) {
        applyFilter();
      }
      if (ev.keyCode == keycodes.escape) {
        cancelFilter();
      }
    }

    watch(() => value.value, () => {
      domEditor.value = value.value
    }, {immediate: true})

    return () => (
      <div class="m-1">
        <div class="space-between align-items-center">
          <div>
            <FontIcon icon="img primary-key"/>
            {(baseTable.value && column.value)
              ?
              <ColumnLabel {...baseTable.value.columns.find(x => x.columnName == column.value!.columnName)} />
              : null
            }
          </div>
          {
            formDisplay.value && formDisplay.value.config.formViewKeyRequested &&
            <DeleteOutlined style={{"cursor": "pointer"}} onClick={cancelFilter}/>
          }
        </div>
        <div class="flex">
          <a-input
            size="small"
            vModel:value={domEditor.value}
            onBlur={applyFilter}
            onKeydown={handleKeyDown}
          />
        </div>
      </div>
    )
  }
})
