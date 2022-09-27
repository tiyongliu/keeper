import {defineComponent, onMounted, PropType, unref} from 'vue'
import ToolStripContainer from '/@/second/buttons/ToolStripContainer.vue'
import DataGrid from '/@/second/datagrid/DataGrid.vue'
import {createGridCache} from '/@/second/keeper-datalib'
import CollectionDataGridCore from '/@/second/datagrid/CollectionDataGridCore'
export default defineComponent({
  name: 'CollectionDataTab',
  props: {
    tabid: {
      type: String as PropType<string>
    },
    conid: {
      type: String as PropType<string>
    },
    database: {
      type: String as PropType<string>
    },
    schemaName: {
      type: String as PropType<string>
    },
    pureName: {
      type: String as PropType<string>
    }
  },
  setup() {
    const cache = unref(createGridCache())

    onMounted(() => {

    })

    return () => (
      <>
        <ToolStripContainer>
          <DataGrid cache={unref(cache)} focusOnVisible gridCoreComponent={CollectionDataGridCore}/>
        </ToolStripContainer>
        {/*<ToolStripCommandButton command="dataGrid.refresh" hideDisabled/>*/}
      </>
    )
  },
})


export const matchingProps = ['conid', 'database', 'schemaName', 'pureName'];
export const allowAddToFavorites = _ => true
