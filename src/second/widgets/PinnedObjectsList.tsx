import {defineComponent, computed, unref, onMounted, watch, ref} from 'vue';
import WidgetsInnerContainer from './WidgetsInnerContainer.vue'
import AppObjectList from '../appobj/AppObjectList'
import { dataBaseStore } from "/@/store/modules/dataBase";
import PinnedAppObject from '../appobj/PinnedAppObject'
import {IPinnedDatabasesItem} from "/@/second/types/IStore";
export default defineComponent({
  setup() {
    const dataBase = dataBaseStore()
    const filteredTables = computed(() => [])

    // const getPinnedDatabases = ref<IPinnedDatabasesItem[]>([])
    // watch(() => dataBase.getPinnedDatabases, () => {
    //   console.log(dataBase.getPinnedDatabases, '*-')
    //   getPinnedDatabases.value = unref(dataBase.getPinnedDatabases)
    // })

    const getPinnedDatabases = computed(() => {
      return dataBase.$state.pinnedDatabases
    })

    return () => (
      <WidgetsInnerContainer>
        <AppObjectList
          list={[...unref(getPinnedDatabases), ...unref(filteredTables)]}
          module={PinnedAppObject}
        />
      </WidgetsInnerContainer>
    )
  }
})
