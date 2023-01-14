<template>
  <TabContent
    :key="tabid"
    v-for="tabid in tabids"
    :tabComponent="mountedTabs[tabid].tabComponent"
    v-bind="mountedTabs[tabid].props"
    :tabid="tabid"
    :tabVisible="tabid == (selectedTab && selectedTab.tabid)"
  />
</template>

<script lang="ts">
import {computed, defineComponent, ref, unref, watchEffect} from 'vue'
import {storeToRefs} from "pinia"
import {difference, keys, map, pickBy} from 'lodash-es'
import {useLocaleStore} from '/@/store/modules/locale'
import tabs from '/@/second/tabs'
import TabContent from './TabContent.vue'
import {TabDefinition} from "/@/store/modules/bootstrap"

function createTabComponent(selectedTab) {
  const tabComponent = tabs[selectedTab.tabComponent]?.default;
  if (tabComponent) {
    return {
      tabComponent,
      props: selectedTab && selectedTab.props,
    };
  }
  return null
}

export default defineComponent({
  name: "TabRegister",
  components: {
    TabContent
  },
  setup() {
    const localeStore = useLocaleStore()
    const {openedTabs} = storeToRefs(localeStore)
    const mountedTabs = ref({})
    const selectedTab = computed(() => (openedTabs.value as TabDefinition[]).find(x => x.selected && x.closedTime == null))

    watchEffect(() => {
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

    const tabids = computed(() => keys(mountedTabs.value))

    return {
      tabids,
      mountedTabs,
      selectedTab,
    }
  }
})
</script>
