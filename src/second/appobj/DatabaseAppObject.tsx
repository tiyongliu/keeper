import {defineComponent, unref, computed, PropType, Ref, watch} from 'vue'
import AppObjectCore from './AppObjectCore.vue'
import { dataBaseStore } from "/@/store/modules/dataBase";
import {IPinnedDatabasesItem} from "/@/second/types/standard.d";
export const extractKey = props => props.name

export default defineComponent({
  name: 'DatabaseAppObject',
  props: {
    data: {
      type: Object as PropType<IPinnedDatabasesItem>,
    },
    passProps: {
      type: Object as unknown as PropType<{
        showPinnedInsteadOfUnpin: boolean
      }>,
    },

  },
  setup(props) {
    const {data, passProps} = props
    const dataBase = dataBaseStore()

    const isPinned = computed(() => {
      return dataBase.$state.pinnedDatabases.find(x => x.name == unref(data!).name && x.connection?._id == unref(data!).connection?._id)
    })


    watch(() => unref(data!), () => {
      console.log(unref(props.data), ` unref(list)          unref(list)`)
    })

    return () => (
      <AppObjectCore
        data={unref(data)}
        title={unref(data!).name}
        extInfo={unref(data!).extInfo}
        icon="img database"
        showPinnedInsteadOfUnpin={unref(passProps)?.showPinnedInsteadOfUnpin}
        onPin={unref(isPinned) ? null : () => dataBase.subscribePinnedDatabases([
          ...unref(dataBase.$state.pinnedDatabases),
          unref(data!)
        ])}
        onUnpin={unref(isPinned) ? () => {
          dataBase.subscribePinnedDatabases(
            dataBase.$state.pinnedDatabases.filter(x => x.name != unref(data!).name || x.connection?._id != unref(data!).connection?._id) as []
          )
        } : null}
      />

    )
  }
})
