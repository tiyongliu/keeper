import {defineComponent, PropType, toRefs, unref, computed} from 'vue'
import {storeToRefs} from 'pinia'
import {get, isEqual, uniqWith} from 'lodash-es'
import AppObjectCore from './AppObjectCore.vue'
import {dataBaseStore} from "/@/store/modules/dataBase"
import {useLocaleStore} from '/@/store/modules/locale'
import {IPinnedDatabasesItem} from "/@/second/typings/types/standard.d"
import {getDatabaseMenuItems} from './PinnedAppObject_'

export default defineComponent({
  name: 'DatabaseAppObject',
  props: {
    data: {
      type: Object as PropType<IPinnedDatabasesItem>,
    },
    passProps: {
      type: Object as PropType<{
        showPinnedInsteadOfUnpin: boolean
      }>,
    },
  },
  setup(props, {attrs}) {
    const dataBase = dataBaseStore()
    const {getCurrentDatabase: currentDatabase} = storeToRefs(dataBase)
    const localeStore = useLocaleStore()
    const {pinnedDatabases} = storeToRefs(localeStore)
    const {data, passProps} = toRefs(props)
    const isPinned = computed(() => unref(pinnedDatabases).find(x => x.name == unref(data)!.name && x.connection?._id == unref(data)!.connection?._id))

    return () => (
      <AppObjectCore
        {...attrs}
        title={unref(data)!.name}
        extInfo={unref(data)!.extInfo}
        icon="img database"
        isBold={get(unref(currentDatabase), 'connection._id') == get(unref(data)!.connection, '_id') &&
          get(unref(currentDatabase), 'name') == unref(data)!.name
        }
        onClick={() => dataBase.subscribeCurrentDatabase(unref(data)!)}
        menu={createMenu}
        showPinnedInsteadOfUnpin={unref(passProps)?.showPinnedInsteadOfUnpin}
        onPin={unref(isPinned) ? null : () => localeStore.subscribePinnedDatabases(uniqWith([
          ...unref(pinnedDatabases),
          unref(data!)
        ], isEqual))}
        onUnpin={unref(isPinned) ? () => {
          localeStore.subscribePinnedDatabases(
            unref(pinnedDatabases).filter(x => x.name != unref(data)!.name || x.connection?._id != unref(data)!.connection?._id) as []
          )
        } : null}
      />
    )
  }
})

function createMenu() {
  return getDatabaseMenuItems()
}
