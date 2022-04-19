import {defineComponent, unref, computed, PropType, Ref} from 'vue'
import AppObjectCore from './AppObjectCore.vue'
import { dataBaseStore } from "/@/store/modules/dataBase";
export const extractKey = props => props.name

export default defineComponent({
  name: 'DatabaseAppObject',
  props: {
    data: {
      type: Object as PropType<{
        name: string,
        extInfo?: PropType<string | Ref<string>>,
        connection: {
          _id: string
        }
      }>,
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

    const ddddd = () => {
      console.log(`3324234`)
      console.log(unref(dataBase.$state.pinnedDatabases))
      console.log(data)
      dataBase.subscribePinnedDatabases([...unref(dataBase.$state.pinnedDatabases), data])
      // return [...unref(dataBase.$state.pinnedDatabases), data]
    }

    return () => (
      <AppObjectCore
        title={unref(data!).name}
        extInfo={unref(data!).extInfo}
        icon="img database"
        showPinnedInsteadOfUnpin={unref(passProps)?.showPinnedInsteadOfUnpin}

        onPin={unref(isPinned) ? null : ddddd}


        // onPin={unref(isPinned) ? null : () => dataBase.subscribePinnedDatabases([
        //   ...dataBase.$state.pinnedDatabases,
        //   unref(data)
        // ])}
        onUnpin={unref(isPinned) ? () => {
          dataBase.subscribePinnedDatabases(
            dataBase.$state.pinnedDatabases.filter(x => x.name != unref(data!).name || x.connection?._id != unref(data!).connection?._id) as []
          )
        } : null}
      />

    )
  }
})
