import {defineComponent, computed, unref, onMounted, watch, ref} from 'vue';
import WidgetsInnerContainer from './WidgetsInnerContainer.vue'
import AppObjectList from '../appobj/AppObjectList'
import { dataBaseStore } from "/@/store/modules/dataBase";
import PinnedAppObject from '../appobj/PinnedAppObject'
export default defineComponent({
  setup() {
    const dataBase = dataBaseStore()
    const filteredTables = computed(() => [])

    // const getPinnedDatabases = ref([])
    // watch(() => dataBase.getPinnedDatabases, () => {
    //   console.log(dataBase.getPinnedDatabases, '*-')
    //   getPinnedDatabases.value = unref(dataBase.getPinnedDatabases)
    // })

    const getPinnedDatabases = computed(() => {
      // console.log(`        dataBase.getPinnedDatabases    `, dataBase.getPinnedDatabases)
      return dataBase.getPinnedDatabases
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
