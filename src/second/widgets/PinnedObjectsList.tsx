import {defineComponent, computed, unref, onMounted} from 'vue';
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
          list={[...dataBase.getPinnedDatabases, ...unref(filteredTables)]}
          module={PinnedAppObject}
        />
      </WidgetsInnerContainer>
    )
  }
})
