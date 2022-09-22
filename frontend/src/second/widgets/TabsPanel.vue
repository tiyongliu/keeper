<template>
  <div class="root">
    <div class="tabs" ref="tabs" @wheel.prevent="handleTabsWheel" >
      <div class="db-wrapper">
        <div class="db-name">
          <div class="db-name-inner">
            <FontIcon icon="icon lock" />
          </div>
          <div class="close-button-right tabCloseButton">
            <FontIcon icon="icon close" />
          </div>
        </div>
        <div class="db-group">
            <div class="file-tab-item">
              <FontIcon  />
              <span class="file-name"></span>
              <span class="lose-button tabCloseButton" @click="(e)=>{}"> <FontIcon icon="icon close" /></span>
            </div>
        </div>
      </div>
    </div>
    <div class="add-icon" title="New query"><FontIcon icon="icon add" /></div>
  </div>
</template>

<script lang="ts">
import {defineComponent, ref, isRef} from 'vue'
import {storeToRefs} from 'pinia'
import FontIcon from '/@/second/icons/FontIcon.vue'
import {useLocaleStore} from '/@/store/modules/locale'
 export default defineComponent({
  name: "TabsPanel",
  components: { FontIcon },
  setup() {
    let tabs = ref<Nullable<HTMLElement>>(null)
    const localeStore = useLocaleStore()
    const {openedTabs} = storeToRefs(localeStore)
    console.log(isRef(openedTabs),'openedTabs-----openedTabs')
    // const tabsWithDb = computed(() => {
    //   return openedTabs.filter(x=> !x.colseTime).map(tab =>({
    //     ...tab,
    //     tabDbName: '',
    //     tabDbKey: '',
    //   }))
    // });

    const handleTabsWheel = (e) =>{
      if(!e.shiftKey){
        // dbgate中是滚动tabs元素，暂未验证出效果
        // await nextTick()
        // tabs.scrollBy({ top: 0, left: e.deltaY < 0 ? -150 : 150, behavior: 'smooth' })
        window.scrollBy({ top: 0, left: e.deltaY < 0 ? -150 : 150, behavior: 'smooth' })
      }
    }

    return {
      handleTabsWheel,
      tabs
    }
  }
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
