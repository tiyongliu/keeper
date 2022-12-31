import {computed, defineComponent, PropType, toRefs, unref} from 'vue'
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
    const {data, passProps} = toRefs(props)

    const bootstrap = useBootstrapStore()
    const {getCurrentDatabase: currentDatabase} = storeToRefs(bootstrap)
    const localeStore = useLocaleStore()
    const {pinnedDatabases} = storeToRefs(localeStore)

    const isPinned = computed(() =>
      !!unref(pinnedDatabases).find(x => unref(x).name == unref(data)!.name && unref(x).connection?._id == unref(data)!.connection?._id))

    return () => (
      <AppObjectCore
        {...attrs}
        title={unref(data)!.name}
        extInfo={unref(data)!.extInfo}
        icon="img database"
        isBold={get(unref(currentDatabase), 'connection._id') == get(data.value!.connection, '_id') &&
          get(unref(currentDatabase), 'name') == data.value!.name
        }
        onClick={() => bootstrap.setCurrentDatabase(data.value!)}
        menu={createMenu}
        showPinnedInsteadOfUnpin={unref(passProps)?.showPinnedInsteadOfUnpin}
        pin={isPinned.value ? null : () => localeStore.updatePinnedDatabases(list => uniqWith([...list, data.value], isEqual))}
        unpin={isPinned.value ?
          () => localeStore.updatePinnedDatabases(list => list.filter(x => x.name != data.value!.name || x.connection?._id != data.value!.connection?._id))
          : null}
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
    {divider: true},
  ]
}
