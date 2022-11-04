import {defineComponent, PropType, toRefs, unref, computed} from 'vue'
import {storeToRefs} from 'pinia'
import {get, isEqual, uniqWith} from 'lodash-es'
import AppObjectCore from './AppObjectCore.vue'
import {useBootstrapStore} from "/@/store/modules/bootstrap"
import {useLocaleStore} from '/@/store/modules/locale'
import {IPinnedDatabasesItem} from "/@/second/typings/types/standard.d"

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
    const bootstrap = useBootstrapStore()
    const {getCurrentDatabase: currentDatabase} = storeToRefs(bootstrap)
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
        onClick={() => bootstrap.subscribeCurrentDatabase(unref(data)!)}
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

export function getDatabaseMenuItems() {
  const handleNewQuery = () => {
  }

  const handleNewTable = () => {
  }

  const handleNewCollection = () => {
  }

  const handleImport = () => {

  }
  const handleExport = () => {

  }

  const handleSqlRestore = () => {
  }

  const handleSqlDump = () => {

  }

  return [
    {onClick: handleNewQuery, text: 'New query', isNewQuery: true},
    {onClick: handleNewTable, text: 'New table'},
    {onClick: handleNewCollection, text: 'New collection'},
    {onClick: handleImport, text: 'Import wizard'},
    {onClick: handleExport, text: 'Export wizard'},
    {onClick: handleSqlRestore, text: 'Restore/import SQL dump'},
    {onClick: handleSqlDump, text: 'Backup/export SQL dump'},
    { divider: true },
  ]
}
