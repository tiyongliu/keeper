import {defineComponent, computed, unref} from 'vue';
import WidgetsInnerContainer from './WidgetsInnerContainer.vue'
import AppObjectList from '/@/second/appobj/AppObjectList'
import { dataBaseStore } from "/@/store/modules/dataBase";
import PinnedAppObject from '../appobj/PinnedAppObject'
export default defineComponent({
  setup() {
    const dataBase = dataBaseStore()
    const filteredTables = computed(() => [])

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
