import {useBootstrapStore} from '/@/store/modules/bootstrap'
import {useLocaleStore} from '/@/store/modules/locale'
import {callWhenAppLoaded} from '/@/second/utility/appLoadManager'
import {storeToRefs} from 'pinia'
import {watch} from 'vue'
import {getConnectionInfo} from './bridge'
import {IPinnedDatabasesItem} from "/@/second/typings/types/standard";

let lastCurrentTab: Nullable<any> = null;

export function subscribeCurrentDbByTab() {
  const bootstrap = useBootstrapStore()
  const locale = useLocaleStore()
  const {openedTabs} = storeToRefs(locale)
  watch(() => openedTabs.value, () => {
    if (openedTabs.value && openedTabs.value.length > 0) {
      const newCurrentTab = (openedTabs.value || []).find(x => x.selected);
      if (newCurrentTab == lastCurrentTab) return;

      const lastTab = lastCurrentTab;
      lastCurrentTab = newCurrentTab;

      if (newCurrentTab) {
        const {conid, database} = newCurrentTab.props || {};
        if (conid && database && (conid != lastTab?.props?.conid || database != lastTab?.props?.database)) {
          const doWork = async () => {
            const connection = await getConnectionInfo({conid});
            bootstrap.setCurrentDatabase({
              connection,
              name: database,
            } as IPinnedDatabasesItem)
          };
          callWhenAppLoaded(doWork);
        }
      }
    }
  }, {immediate: true})
}
