<template>
  <div class="collection-json-view-wrapper">
    <div class="toolbar">
      <Pager v-model:skip="skip" v-model:limit="limit" @load="reload"/>
    </div>
    <div class="json">
      <template v-for="rowIndex in rangeRowCount" :key="rowIndex">
        <CollectionJsonRow :grider="grider" :rowIndex="rowIndex"/>
      </template>
    </div>
  </div>
  <LoadingInfo v-if="isLoading" wrapper message="Loading data"/>
</template>

<script lang="ts">
import {computed, defineComponent, onMounted, PropType, ref, toRefs, watchEffect} from 'vue'
import LoadingInfo from '/@/second/elements/LoadingInfo.vue'
import Pager from '/@/second/elements/Pager.vue'
import CollectionJsonRow from './CollectionJsonRow'

import {ChangeCacheFunc, GridDisplay} from "/@/second/keeper-datalib";
import {loadCollectionDataPage} from '/@/second/datagrid/CollectionDataGridCore'
import ChangeSetGrider from "/@/second/datagrid/ChangeSetGrider";
import {range} from 'lodash-es'

export default defineComponent({
  name: "CollectionJsonView",
  components: {
    Pager,
    CollectionJsonRow,
    LoadingInfo,
  },
  props: {
    conid: {
      type: String as PropType<string>
    },
    database: {
      type: String as PropType<string>
    },
    cache: {
      type: Object as PropType<{ refreshTime: number }>
    },
    display: {
      type: Object as PropType<GridDisplay>
    },
    setCache: {
      type: Function as PropType<ChangeCacheFunc>
    },
    changeSetState: {
      type: Object as PropType<any>
    },
    dispatchChangeSet: {
      type: Function as PropType<(action: any) => void>
    },
    loadedRows: {
      type: Array as PropType<any[]>,
      default: () => []
    },
  },
  setup(props, {attrs}) {
    const {cache, changeSetState, dispatchChangeSet, display, loadedRows} = toRefs(props)
    const isLoading = ref(false)
    const loadedTime = ref<Nullable<number>>(null)
    const loadedRowsRW = ref(loadedRows.value)

    const grider = ref()
    const skip = ref(0)
    const limit = ref(50)

    async function loadData() {
      isLoading.value = true
      loadedRowsRW.value = await loadCollectionDataPage(
        Object.assign({}, props, attrs),
        parseInt(skip.value) || 0, parseInt(limit.value) || 50
      )
      isLoading.value = false
      loadedTime.value = new Date().getTime()
    }

    watchEffect(() => {
      if (cache.value && loadedTime.value && cache.value?.refreshTime > loadedTime.value) {
        loadData()
      }
    })

    onMounted(() => {
      loadData()
    })

    watchEffect(() => {
      grider.value = new ChangeSetGrider(
        loadedRowsRW.value,
        changeSetState.value,
        dispatchChangeSet.value,
        display.value!,
      )
    })

    const rangeRowCount = computed(() => range(0, grider.value.rowCount))

    function reload() {
      display.value && display.value?.reload()
    }

    return {
      skip,
      limit,
      reload,
      isLoading,
      grider,
      rangeRowCount
    }
  }
})
</script>

<style scoped>
.collection-json-view-wrapper {
  display: flex;
  flex-direction: column;
  position: absolute;
  left: 0;
  top: 0;
  right: 0;
  bottom: 0;
}

.json {
  overflow: auto;
  flex: 1;
  /* position: relative; */
}

.toolbar {
  background: var(--theme-bg-1);
  display: flex;
  border-bottom: 1px solid var(--theme-border);
  border-top: 2px solid var(--theme-border);
  margin-bottom: 3px;
}
</style>
