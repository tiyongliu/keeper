<template>
  <div class="cell-data-wrapper">
    <WidgetTitle>Cell data view</WidgetTitle>
    <div class="main">
      <div class="toolbar">
        Format:<span>&nbsp;</span>
        <Select
          style="width: 200px"
          size="small"
          v-model:value="selectedFormatType"
          :options="[
              { value: 'autodetect', label: `Autodetect - ${autodetectFormat.title}` },
              ...formats.map(fmt => ({ label: fmt.title, value: fmt.type }))
          ]"/>
      </div>
      <div class="data">
        <ErrorInfo
          v-if="usedFormat.single && transformSelection.length != 1"
          message="Must be selected one cell" alignTop/>
        <ErrorInfo
          v-else-if="usedFormat == null"
          message="Format not selected" alignTop/>
        <ErrorInfo
          v-else-if="!selection || transformSelection.length == 0"
          message="No data selected" alignTop/>
        <component v-else :is="usedFormat.component" :selection="selection"/>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {computed, defineProps, ref, toRef, watch, withDefaults} from 'vue'
import {isArray, isPlainObject, isString} from "lodash-es";
import {Select} from 'ant-design-vue'
import ErrorInfo from '/@/second/elements/ErrorInfo.vue'
import HtmlCellView from '/@/second/celldata/HtmlCellView.vue';
import MapCellView from '/@/second/celldata/MapCellView.vue';
import TextCellViewWrap from '/@/second/celldata/TextCellViewWrap.vue';
import TextCellViewNoWrap from '/@/second/celldata/TextCellViewNoWrap.vue';
import JsonCellView from '/@/second/celldata/JsonCellView.vue';
import JsonRowView from '/@/second/celldata/JsonRowView.vue';
import PictureCellView from '/@/second/celldata/PictureCellView.vue';
import WidgetTitle from './WidgetTitle.vue'
import {selectionCouldBeShownOnMap} from '/@/second/elements/MapView'

const formats = [
  {
    type: 'textWrap',
    title: 'Text (wrap)',
    component: TextCellViewWrap,
    single: false,
  },
  {
    type: 'text',
    title: 'Text (no wrap)',
    component: TextCellViewNoWrap,
    single: false,
  },
  {
    type: 'json',
    title: 'Json',
    component: JsonCellView,
    single: true,
  },
  {
    type: 'jsonRow',
    title: 'Json - Row',
    component: JsonRowView,
    single: false,
  },
  {
    type: 'picture',
    title: 'Picture',
    component: PictureCellView,
    single: true,
  },
  {
    type: 'html',
    title: 'HTML',
    component: HtmlCellView,
    single: false,
  },
  {
    type: 'map',
    title: 'Map',
    component: MapCellView,
    single: false,
  },
]
const selectedFormatType = ref<string>('autodetect')

const props = withDefaults(defineProps<{ selection?: any[] | undefined }>(), {
  selection: undefined
})

const selection = toRef(props, 'selection')
const transformSelection = ref([])

const autodetectFormatType = computed(() => autodetect(transformSelection.value))
const autodetectFormat = computed(() => formats.find(x => x.type == autodetectFormatType.value))

const usedFormatType = computed(() => selectedFormatType.value == 'autodetect' ? autodetectFormatType.value : selectedFormatType.value)
const usedFormat = computed(() => formats.find(x => x.type == usedFormatType.value))

watch(() => selection.value, () => {
  //todo 后期需要把它写活，只是当前没有数据，无法看到效果，故采用默认值
  //$: selection = $selectedCellsCallback ? $selectedCellsCallback() : [];
  transformSelection.value = []
})

function autodetect(selection) {
  if (selection[0]?.engine?.databaseEngineTypes?.includes('document')) {
    return 'jsonRow';
  }

  if (selectionCouldBeShownOnMap(selection)) {
    return 'map'
  }

  const value = selection.length == 1 ? selection[0].value : null;
  if (isString(value)) {
    if (value.startsWith('[') || value.startsWith('{')) return 'json';
  }
  if (isPlainObject(value) || isArray(value)) {
    return 'json';
  }
  return 'textWrap';
}
</script>

<style lang="less" scoped>
.cell-data-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;

  .main {
    display: flex;
    flex: 1;
    flex-direction: column;
  }

  .toolbar {
    display: flex;
    background: var(--theme-bg-1);
    align-items: center;
    border-bottom: 1px solid var(--thene-border);
    margin: 2px;
  }

  .data {
    display: flex;
    flex: 1;
    position: relative;
  }
}
</style>
