import {defineComponent} from 'vue'
import LoadingDataGridCore from './LoadingDataGridCore'
import eb_system_config from "/@/second/tabs/eb_system_config.json";
export default defineComponent({
  name: 'CollectionDataGridCore',
  setup(_, {attrs}) {

    function dataPageAvailable() {
      return true;
      // const { display } = props;
      // const sql = display.getPageQuery(0, 1);
      // return !!sql;
    }

    async function loadRowCount() {
      return 179
    }

    return () => (
      <LoadingDataGridCore
        {...attrs}
        loadDataPage={loadCollectionDataPage}
        dataPageAvailable={dataPageAvailable}
        loadRowCount={loadRowCount}
      />
    )
  }
})

export async function loadCollectionDataPage(props, offset, limit): Promise<any[]> {
  // const { conid, database } = props
  return eb_system_config.rows
}
