import {defineComponent, unref, computed, PropType, watch, toRefs} from 'vue'
import { isEqual, uniqWith } from 'lodash-es'
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
    const dataBase = dataBaseStore()
    const {data, passProps} = toRefs(props)
    const isPinned = computed(() => {
      return dataBase.$state.pinnedDatabases.find(x => x.name == unref(data)!.name && x.connection?._id == unref(data)!.connection?._id)
    })

    watch(props, () => {
      console.log(unref(props), ` unref(list)     00     unref(list)`)
    })

    return () => (
      <AppObjectCore
        data={unref(data)}
        title={unref(data)!.name}
        extInfo={unref(data)!.extInfo}
        icon="img database"
        showPinnedInsteadOfUnpin={unref(passProps)?.showPinnedInsteadOfUnpin}
        onPin={unref(isPinned) ? null : () => dataBase.subscribePinnedDatabases(uniqWith([
          ...unref(dataBase.$state.pinnedDatabases),
          unref(data!)
        ], isEqual))}
        onUnpin={unref(isPinned) ? () => {
          dataBase.subscribePinnedDatabases(
            dataBase.$state.pinnedDatabases.filter(x => x.name != unref(data)!.name || x.connection?._id != unref(data)!.connection?._id) as []
          )
        } : null}
      />

    )
  }
})
