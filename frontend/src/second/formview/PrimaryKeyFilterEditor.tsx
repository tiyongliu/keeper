import {computed, defineComponent, PropType, ref, toRefs, watch} from 'vue'
import {TableFormViewDisplay} from '/@/second/keeper-datalib'
import FontIcon from '/@/second/icons/FontIcon.vue'
import ColumnLabel from '/@/second/elements/ColumnLabel.vue'
import InlineButton from '/@/second/buttons/InlineButton.vue'
import keycodes from '/@/second/utility/keycodes'
import {ColumnReference} from '/@/second/keeper-types'

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
      (formDisplay.value && column.value) && formDisplay.value.requestKeyValue(column.value.columnName, domEditor.value)
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
        <div class="space-between">
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
            <InlineButton square onClick={cancelFilter}>
              <FontIcon icon="icon delete"/>
            </InlineButton>
          }
        </div>
        <div class="flex">
          <a-input
            vModel:value={domEditor.value}
            size="small"
            onBlur={applyFilter}
            onKeydown={handleKeyDown}
          />
        </div>
      </div>
    )
  }
})
