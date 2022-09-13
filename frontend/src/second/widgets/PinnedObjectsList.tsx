import {computed, defineComponent, unref} from 'vue'
import {storeToRefs} from 'pinia'
import {dataBaseStore} from '/@/store/modules/dataBase'
import {useLocaleStore} from '/@/store/modules/locale'
import AppObjectList from '/@/second/appobj/AppObjectList'
import PinnedAppObject from '/@/second/appobj/PinnedAppObject'
import WidgetsInnerContainer from './WidgetsInnerContainer.vue'

export default defineComponent({
  name: 'PinnedObjectsList',
  setup() {
    const dataBase = dataBaseStore()
    const localeStore = useLocaleStore()
    const {pinnedDatabases, pinnedTables} = storeToRefs(localeStore)
    const {currentDatabase} = storeToRefs(dataBase)

    const filteredTables = computed(() => (pinnedTables.value || []).filter(
      x => x.conid == currentDatabase.value?.connection._id && x.database == currentDatabase.value?.name
    ))

    return () => (
      <WidgetsInnerContainer>
        <AppObjectList
          list={[...unref(pinnedDatabases), ...unref(filteredTables)]}
          module={PinnedAppObject}
        />
      </WidgetsInnerContainer>
    )
  }
})
