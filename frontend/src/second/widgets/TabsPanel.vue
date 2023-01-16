<template>
  <div class="root">
    <div class="tabs" ref="domTabs" @wheel.prevent="handleTabsWheel">
      <div class="db-wrapper" v-for="tabGroup in groupedTabs">
        <div class="db-name"
             :class="{selected: draggingDbGroup ? tabGroup.grpid == draggingDbGroupTarget?.grpid : tabGroup.tabDbKey == currentDbKey}"
             @mouseup="e => {
               if (e.button == 1) {
                  closeMultipleTabs(tab => tabGroup.tabs.find(x => x.tabid == tab.tabid));
               } else {
                  handleSetDb(tabGroup.tabs[0].props);
               }
             }"
             :draggable="true"
             @dragstart="(e) => draggingDbGroup = tabGroup"
             @dragenter="(e) => draggingDbGroupTarget = tabGroup"
             @drop="dragDropTabs(draggingDbGroup.tabs, tabGroup.tabs)"
             @dragend="(e) => {
               draggingDbGroup = null
               draggingDbGroupTarget = null
             }"
        >
          <div class="db-name-inner">
            <FontIcon :icon="getDbIcon(tabGroup.tabDbKey)"/>
            {{ tabGroup.tabDbName }}
            <FontIcon
              v-if="connectionList && connectionList?.find(x => x._id == tabGroup.tabs[0]?.props?.conid)?.isReadOnly"
              icon="icon lock"/>
          </div>
          <div class="close-button-right tabCloseButton"
               @click="closeMultipleTabs(tab => tabGroup.tabs.find(x => x.tabid == tab.tabid))">
            <FontIcon icon="icon close"/>
          </div>
        </div>
        <div class="db-group">
          <div class="file-tab-item"
               v-for="tab in tabGroup.tabs"
               :id="`file-tab-item-${tab.tabid}`"
               :class="{selected: draggingTab || draggingDbGroup ? tab.tabid == draggingTabTarget?.tabid : tab.selected}"
               @click="handleTabClick($event, tab.tabid)"
               @mouseup="handleMouseUp($event, tab.tabid)"
               :draggable="true"
               @dragstart="handleDragstart(tab)"
               @dragenter="(e) => draggingTabTarget = tab"
               @drop="(e) => {
                 if (draggingTab) {
                   dragDropTabs([draggingTab], [tab])
                 }
                 if (draggingDbGroup) {
                   dragDropTabs(draggingDbGroup.tabs, [tab])
                 }
               }"
               @dragend="(e) => {
                 draggingTab = null
                 draggingTabTarget = null
               }"
          >
            <FontIcon :icon="tab.busy ? 'icon loading' : tab.icon"/>
            <span class="file-name">{{ tab.title }}</span>
            <span class="close-button tabCloseButton" @click="closeTab(tab.tabid)"> <FontIcon
              icon="icon close"/></span>
          </div>
        </div>
      </div>
    </div>
    <div class="add-icon" title="New query">
      <FontIcon icon="icon add"/>
    </div>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, nextTick, ref, unref} from 'vue'
import {findIndex, max, min} from 'lodash-es'
import {storeToRefs} from 'pinia'
import FontIcon from '/@/second/icons/FontIcon.vue'
import {useLocaleStore} from '/@/store/modules/locale'
import {useBootstrapStore} from '/@/store/modules/bootstrap'
import {getConnectionInfo} from '/@/api/bridge'
import {getTabDbKey, groupTabs, sortTabs} from '/@/second/utility/openNewTab'
import {setSelectedTab} from '/@/second/utility/common'
import {
  closeMultipleTabs,
  closeTab,
  closeWithOtherDb,
  closeWithSameDb,
  getDbIcon,
  getTabDbName
} from './TabsPanel_'
import {useClusterApiStore} from '/@/store/modules/clusterApi'
import {IPinnedDatabasesItem} from "/@/second/typings/types/standard";

export default defineComponent({
  name: "TabsPanel",
  components: {FontIcon},
  setup() {
    let domTabs = ref<Nullable<HTMLElement>>(null)
    const localeStore = useLocaleStore()
    const {openedTabs} = storeToRefs(localeStore)
    const bootstrap = useBootstrapStore()
    const {currentDatabase} = storeToRefs(bootstrap)

    const clusterApi = useClusterApiStore()
    const {connectionList} = storeToRefs(clusterApi)

    const draggingTab = ref(null)
    const draggingTabTarget = ref(null)
    const draggingDbGroup = ref(null)
    const draggingDbGroupTarget = ref(null)

    const handleTabsWheel = async (e) => {
      if (!e.shiftKey) {
        domTabs.value!.scrollBy({top: 0, left: e.deltaY < 0 ? -150 : 150, behavior: 'smooth'})
      }
    }

    const currentDbKey = computed(() => {
      return currentDatabase.value && currentDatabase.value.name && currentDatabase.value.connection
        ? `database://${currentDatabase.value.name}-${currentDatabase.value.connection._id}`
        : currentDatabase.value && currentDatabase.value.connection
          ? `server://${currentDatabase.value.connection._id}`
          : '_no'
    })

    const tabsWithDb = computed(() => {
      return openedTabs.value
        .filter(x => !x.closedTime)
        .map(tab => ({
          ...tab,
          tabDbName: getTabDbName(tab, connectionList.value),
          tabDbKey: getTabDbKey(tab),
        }))
    })

    const groupedTabs = computed(() => groupTabs(tabsWithDb.value))

    const handleTabClick = (e, tabid) => {
      if (e.target.closest('.tabCloseButton')) {
        return;
      }
      setSelectedTab(tabid);
    }

    const handleMouseUp = (e, tabid) => {
      if (e.button == 1) {
        e.preventDefault();
        closeTab(tabid)
      }
    }

    const handleSetDb = async props => {
      const {conid, database} = props || {};
      if (conid) {
        const connection = await getConnectionInfo({conid, database})
        if (connection) {
          bootstrap.setCurrentDatabase({
            connection: unref(connection),
            name: unref(database)
          } as Nullable<IPinnedDatabasesItem>)
          // $currentDatabase = { connection, name: database };
          return;
        }
      }
      bootstrap.setCurrentDatabase(null)
    }

    function getDatabaseContextMenu(tabs) {
      const {tabid, props} = tabs[0];
      const {conid, database} = props || {};

      return [
        conid &&
        database && [
          {
            text: `Close tabs with DB ${database}`,
            onClick: () => closeWithSameDb(tabid),
          },
          {
            text: `Close tabs with other DB than ${database}`,
            onClick: () => closeWithOtherDb(tabid),
          },
        ],
      ];
    }

    async function handleDragstart(tab) {
      draggingTab.value = tab
      await nextTick()
      setSelectedTab(tab.tabid)
    }

    function dragDropTabs(draggingTabs, targetTabs) {
      if (unref(draggingTabs).find(x => unref(targetTabs).find(y => x.tabid == y.tabid))) return;

      const items = sortTabs(unref(openedTabs).filter(x => x.closedTime == null))
      const dstIndexes = unref(targetTabs).map(targetTab => findIndex(items, x => x.tabid == targetTab.tabid))
      const dstIndexFirst = min(dstIndexes) as number;
      const dstIndexLast = max(dstIndexes) as number;
      const srcIndex = findIndex(items, x => x.tabid == unref(draggingTabs)[0].tabid);
      if (srcIndex < 0 || dstIndexFirst < 0 || dstIndexLast < 0) {
        console.warn('Drag tab index not found');
        return;
      }

      const newItems =
        dstIndexFirst < srcIndex
          ? [
            ...items.slice(0, dstIndexFirst),
            ...unref(draggingTabs),
            ...items.slice(dstIndexFirst).filter(x => !unref(draggingTabs).find(y => y.tabid == x.tabid)),
          ]
          : [
            ...items.slice(0, dstIndexLast + 1).filter(x => !unref(draggingTabs).find(y => y.tabid == x.tabid)),
            ...unref(draggingTabs),
            ...items.slice(dstIndexLast + 1),
          ]

      localeStore.updateOpenedTabs(tabs => {
        return tabs.map(x => {
          const index = findIndex(newItems, y => y.tabid == x.tabid)
          if (index >= 0) {
            return {
              ...x,
              tabOrder: index + 1,
            };
          }
          return x;
        })
      })
    }

    function handleDrop(tab) {
      if (draggingTab.value) {
        dragDropTabs([draggingTab], [tab]);
      }
      if (draggingDbGroup.value) {
        dragDropTabs((draggingDbGroup.value as any).tabs, [tab])
      }
    }

    return {
      domTabs,
      connectionList,
      currentDbKey,
      groupedTabs,
      draggingTab,
      draggingTabTarget,
      draggingDbGroup,
      draggingDbGroupTarget,

      handleTabsWheel,
      handleTabClick,
      handleMouseUp,
      getDbIcon,
      closeMultipleTabs,
      handleSetDb,
      dragDropTabs,
      handleDragstart,
      handleDrop,
      getDatabaseContextMenu,
      closeTab,
    }
  },
})
</script>

<style scoped>
.root {
  position: absolute;
  left: 0;
  top: 0;
  right: 0;
  bottom: 0;
}

.add-icon {
  position: absolute;
  right: 5px;
  font-size: 20pt;
  top: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  color: var(--theme-font-2);
  cursor: pointer;
}

.add-icon:hover {
  color: var(--theme-font-1);
}

.tabs {
  height: var(--dim-tabs-panel-height);
  display: flex;
  overflow-x: auto;
  position: absolute;
  left: 0;
  top: 0;
  right: 35px;
  bottom: 0;
}

.tabs::-webkit-scrollbar {
  height: 7px;
}

.db-group {
  display: flex;
  flex: 1;
  align-content: stretch;
}

.db-wrapper {
  display: flex;
  flex-direction: column;
  align-items: stretch;
}

.db-name {
  display: flex;
  text-align: center;
  font-size: 8pt;
  border-bottom: 1px solid var(--theme-border);
  border-right: 1px solid var(--theme-border);
  cursor: pointer;
  user-select: none;
  padding: 1px;
  position: relative;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.db-name-inner {
  justify-content: center;
  flex-grow: 1;
}

/* .db-name:hover {
  background-color: var(--theme-bg-3);
} */
.db-name.selected {
  background-color: var(--theme-bg-0);
}

.file-tab-item {
  border-right: 1px solid var(--theme-border);
  border-bottom: 2px solid var(--theme-border);
  padding-left: 15px;
  padding-right: 15px;
  flex-shrink: 1;
  flex-grow: 1;
  min-width: 10px;
  display: flex;
  align-items: center;
  cursor: pointer;
  user-select: none;
}

.file-tab-item.selected {
  background-color: var(--theme-bg-0);
}

.file-name {
  margin-left: 5px;
  white-space: nowrap;
  flex-grow: 1;
}

.close-button {
  margin-left: 5px;
  color: var(--theme-font-3);
}

.close-button-right {
  margin-left: 5px;
  margin-right: 5px;
  color: var(--theme-font-3);
}

.close-button:hover {
  color: var(--theme-font-1);
}

.close-button-right:hover {
  color: var(--theme-font-1);
}
</style>
