import {defineComponent, ref, unref, watch,} from 'vue'
import {storeToRefs} from "pinia"
import {difference, keys, map, pickBy} from 'lodash-es'
import {useLocaleStore} from '/@/store/modules/locale'
import tabs from '/@/second/tabs'
import TabContent from './TabContent.vue'
import {TabDefinition} from "/@/store/modules/bootstrap";

export default defineComponent({
  name: 'TabRegister',
  setup() {
    const localeStore = useLocaleStore()
    const {openedTabs} = storeToRefs(localeStore)
    const selectedTab = ref()
    const mountedTabs = ref({})

    watch(() => openedTabs.value, () => {
      selectedTab.value = (openedTabs.value as TabDefinition[]).find(x => x.selected && x.closedTime == null)
    })

    watch(() => [mountedTabs.value, selectedTab.value], () => {
      if (difference(
        keys(mountedTabs.value),
        map(
          (openedTabs.value as TabDefinition[]).find(x => x.closedTime == null) as any,
          'tabid'
        )
      ).length
      ) {
        mountedTabs.value = pickBy(mountedTabs.value, (_, k) => (openedTabs.value as TabDefinition[]).find(x => x.tabid == k && x.closedTime == null))
      }
    })

    watch(() => selectedTab.value, () => {
      if (selectedTab.value) {
        const {tabid} = unref(selectedTab)!
        if (tabid && !mountedTabs.value[tabid]) {
          const newTab = createTabComponent(selectedTab.value);
          if (newTab) {
            mountedTabs.value = {
              ...mountedTabs.value,
              [tabid]: newTab
            }
          }
        }
      }
    })

    return () => (
      <>
        {keys(mountedTabs.value).map(tabid => <TabContent
          tabComponent={unref(mountedTabs.value[tabid].tabComponent)}
          {...unref(mountedTabs.value[tabid].props)}
          tabid={tabid}
          tabVisible={tabid == (selectedTab.value && selectedTab.value.tabid)}
        />)}
      </>
    )
  }
})

function createTabComponent(selectedTab) {
  const tabComponent = tabs[selectedTab.tabComponent]?.default;
  if (tabComponent) {
    return {
      tabComponent,
      props: selectedTab && selectedTab.props,
    };
  }
  return null;
}