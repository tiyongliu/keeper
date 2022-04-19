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
    }
  },
  setup(props) {
    const {data, passProps} = props
    const dataBase = dataBaseStore()

    const isPinned = computed(() => {
      return dataBase.$state.pinnedDatabases.find(x => x.name == unref(data!).name && x.connection?._id == unref(data!).connection?._id)
    })

    console.log(props, `pppppppppppppppppppp`)
    console.log(isPinned, `pppppppppppppppppppp`)
    return () => (
      <AppObjectCore
        title={unref(data!).name}
        extInfo={unref(data!).extInfo}
        icon="img database"
        showPinnedInsteadOfUnpin={unref(passProps)?.showPinnedInsteadOfUnpin}
        onPin={unref(isPinned) ? null : () => console.log(`22`)}
        onUnpin={unref(isPinned) ? () => {
          dataBase.subscribePinnedDatabases(
            dataBase.$state.pinnedDatabases.filter(x => x.name != unref(data!).name || x.connection?._id != unref(data!).connection?._id) as []
          )
        } : null}
      />

    )
  }
})
