<template>
  <HorizontalSplitter :initialValue="getInitialManagerSize()" :size="managerSize">
<!--<div class="left">
      <WidgetColumnBar>
        <WidgetColumnBarItem title="Columns" name="columns" height="45%">
          <ColumnManager />
        </WidgetColumnBarItem>

        <WidgetColumnBarItem title="Filters" name="jsonFilters" height="30%">

        </WidgetColumnBarItem>
      </WidgetColumnBar>

    </div>-->

    <template #2>
      <VerticalSplitter initialValue="70%" :isSplitter="false">
        <template #1>
          <component v-if="isFormView" :is="formViewComponent" v-bind="fullProps"/>
          <component v-else-if="isJsonView" :is="jsonViewComponent" v-bind="fullProps"/>
          <component v-else :is="gridCoreComponent" v-bind="fullProps" :macroPreview="selectedMacro"/>
        </template>

        <template #2>
          <MacroDetail v-if="selectedMacro" :onExecute="handleExecuteMacro"/>
        </template>
      </VerticalSplitter>
    </template>
  </HorizontalSplitter>
</template>

<script lang="ts">
import {
  Component,
  defineComponent,
  PropType,
  ref,
  toRaw,
  onMounted,
  computed,
  toRefs,
  provide
} from 'vue'
import {isNumber, fromPairs, mapKeys} from 'lodash-es'
import HorizontalSplitter from '/@/second/elements/HorizontalSplitter.vue'
import WidgetColumnBar from '/@/second/widgets/WidgetColumnBar.vue'
import WidgetColumnBarItem from '/@/second/widgets/WidgetColumnBarItem.vue'
import VerticalSplitter from '/@/second/elements/VerticalSplitter.vue'
import MacroDetail from '/@/second/freetable/MacroDetail.vue'
import ColumnManager from '/@/second/datagrid/ColumnManager.vue'
import {getLocalStorage} from '/@/second/utility/storageCache'

import {
  TableGridDisplay,
  TableFormViewDisplay, GridConfig
} from '/@/second/keeper-datalib'


function extractMacroValuesForMacro(macroValues, macro) {
  // return {};
  if (!macro) return {};
  return {
    ...fromPairs((macro.args || []).filter(x => x.default != null).map(x => [x.name, x.default])),
    ...mapKeys(macroValues, (v, k) => k.replace(/^.*#/, '')),
  };
}

export default defineComponent({
  name: "DataGrid",
  props: {
    gridCoreComponent: {
      type: [String, Object] as PropType<string | Component | any>,
    },
    formViewComponent: {
      type: [String, Object] as PropType<string | Component>,
    },
    jsonViewComponent: {
      type: [String, Object] as PropType<string | Component>,
    },
    config: {
      type: Object as PropType<GridConfig>,
    },
    display: {
      type: Object as PropType<TableGridDisplay>
    },
    formDisplay: {
      type: Object as PropType<TableFormViewDisplay>
    }
  },
  components: {
    HorizontalSplitter,
    VerticalSplitter,
    MacroDetail,
    WidgetColumnBar,
    WidgetColumnBarItem,
    ColumnManager
  },
  emits: ['runMacro'],
  setup(props, {attrs, emit}) {
    const gridCoreComponent = toRaw(props.gridCoreComponent)
    const formViewComponent = toRaw(props.formViewComponent)
    const jsonViewComponent = toRaw(props.jsonViewComponent)

    const managerSize = ref(0)

    const selectedMacro = ref(null)
    provide('selectedMacro', selectedMacro)
    const macroValues = ref({})
    provide('macroValues', macroValues)

    function getInitialManagerSize() {
      const width = getLocalStorage('dataGridManagerWidth')
      if (isNumber(width) && width > 30 && width < 500) {
        return `${width}px`;
      }
      return '300px';
    }

    const {config, formDisplay, selectedCellsPublished} = toRefs(props)
    const isFormView = computed(() => !!(formDisplay.value && formDisplay.value.config && formDisplay.value.config.isFormView))
    const isJsonView = computed(() => !!(config.value?.isJsonView))

    const handleExecuteMacro = () => {
      emit('runMacro', () => (selectedMacro.value, extractMacroValuesForMacro(macroValues.value, selectedMacro.value), selectedCellsPublished.value()))
      selectedMacro.value = null
    }

    return {
      gridCoreComponent,
      formViewComponent,
      jsonViewComponent,
      managerSize,
      getInitialManagerSize,
      handleExecuteMacro,
      selectedMacro,
      isFormView,
      isJsonView,
      fullProps: {
        ...Object.assign(props, attrs)
      },
    }
  }
})
</script>

<style scoped>

</style>
