import {defineComponent, computed, unref} from 'vue';
import WidgetsInnerContainer from './WidgetsInnerContainer.vue'
import AppObjectList from '../appobj/AppObjectList'
import { dataBaseStore } from "/@/store/modules/dataBase";
import PinnedAppObject from '../appobj/PinnedAppObject'
export default defineComponent({
  setup() {
    const dataBase = dataBaseStore()
    const filteredTables = computed(() => [])

    return () => (
      <WidgetsInnerContainer>
        <AppObjectList
          list={[...dataBase.$state.pinnedDatabases, ...unref(filteredTables)]}
          module={PinnedAppObject}
        />
      </WidgetsInnerContainer>
    )
  }
})
