import {watch} from 'vue'
import {storeToRefs} from 'pinia'
import {compact} from 'lodash-es'
import {useLocaleStore} from '/@/store/modules/locale'
import {useBootstrapStore} from '/@/store/modules/bootstrap'

export function subscribeRecentDatabaseSwitch() {
  const bootstrap = useBootstrapStore()
  const {currentDatabase} = storeToRefs(bootstrap)
  const locale = useLocaleStore()
  watch(() => currentDatabase.value, () => {
    if (!currentDatabase.value) return
    const value = currentDatabase.value
    locale.updateRecentDatabases(list => {
      return [
        value,
        // @ts-ignore
        ...compact(list).filter(x => x.name != value.name || x.connection?._id != value.connection?._id)
      ].slice(0, 10)
    })
  }, {immediate: true})
}
