<template>
  <HorizontalSplitter
    :initialValue="getInitialManagerSize()"
    @dispatchSize="dispatchSize"
    :hideFirst="collapsedLeftColumnStore">
    <template #1>
      <div class="left">
        <WidgetColumnBar>
          <WidgetColumnBarItem
            title="Columns"
            name="columns"
            height="45%"
            :show="columnsShow">
            <ColumnManager
              v-bind="Object.assign({}, $props, $attrs)"
              :managerSize="managerSize"
              :isJsonView="isJsonView"
              :isDynamicStructure="isDynamicStructure"
              ref="domColumnManager"/>
          </WidgetColumnBarItem>

          <WidgetColumnBarItem
            title="Filters"
            name="jsonFilters"
            height="30%"
            :skip="jsonFiltersSkip">
            <JsonViewFilters
              v-bind="Object.assign({}, $props, $attrs)"
              :managerSize="managerSize"
              :isDynamicStructure="isDynamicStructure"
              :useEvalFilters="useEvalFilters"/>
          </WidgetColumnBarItem>
        </WidgetColumnBar>

      </div>
    </template>
    <template #2>
      <VerticalSplitter initialValue="70%" :isSplitter="false">
        <template #1>
          <component
            v-if="isFormView"
            :is="formViewComponent"
            v-bind="Object.assign({}, $props, $attrs)"/>
          <component
            v-else-if="isJsonView"
            :is="jsonViewComponent"
            v-bind="Object.assign({}, $props, $attrs)"/>
          <component
            v-else
            :is="gridCoreComponent"
            v-bind="Object.assign({}, $props, $attrs)"
            :macroPreview="selectedMacro"
            :formViewAvailable="!!formViewComponent && !!formDisplay"/>
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
  computed,
  defineComponent,
  inject,
  PropType,
  provide,
  Ref,
  ref,
  toRaw,
  toRefs,
  unref,
  watch
} from 'vue'
import {fromPairs, isNumber, mapKeys} from 'lodash-es'
import HorizontalSplitter from '/@/second/elements/HorizontalSplitter.vue'
import WidgetColumnBar from '/@/second/widgets/WidgetColumnBar.vue'
import WidgetColumnBarItem from '/@/second/widgets/WidgetColumnBarItem.vue'
import VerticalSplitter from '/@/second/elements/VerticalSplitter.vue'
import MacroDetail from '/@/second/freetable/MacroDetail.vue'
import ColumnManager from '/@/second/datagrid/ColumnManager.vue'
import JsonViewFilters from '/@/second/jsonview/JsonViewFilters'
import {getLocalStorage, setLocalStorage} from '/@/second/utility/storageCache'
import {GridConfig, GridDisplay, TableFormViewDisplay,} from '/@/second/keeper-datalib'
import {Nullable} from "/@/utils/types";

function extractMacroValuesForMacro(macroValues, macro) {
  // return {};
  if (!macro) return {};
  return {
    ...fromPairs((macro.args || []).filter(x => x.default != null).map(x => [x.name, x.default])),
    ...mapKeys(macroValues, (_, k) => k.replace(/^.*#/, '')),
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
      type: Object as PropType<GridDisplay>
    },
    formDisplay: {
      type: Object as PropType<TableFormViewDisplay>
    },
    macroCondition: {
      type: Function as PropType<(macro: any) => boolean>
    },
    useEvalFilters: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    isDetailView: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    showReferences: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    showMacros: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    expandMacros: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    freeTableColumn: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    isDynamicStructure: {
      type: Boolean as PropType<boolean>,
      default: false
    },

  },
  components: {
    HorizontalSplitter,
    VerticalSplitter,
    MacroDetail,
    WidgetColumnBar,
    WidgetColumnBarItem,
    ColumnManager,
    JsonViewFilters
  },
  emits: ['runMacro'],
  setup(props, {emit}) {
    const {
      config,
      formDisplay,
      display,
      freeTableColumn,
      useEvalFilters,
      isDynamicStructure
    } = toRefs(props)

    const domColumnManager = ref<Nullable<{ setSelectedColumns: (value: unknown[]) => void }>>(null)
    const gridCoreComponent = toRaw(props.gridCoreComponent)
    const formViewComponent = toRaw(props.formViewComponent)
    const jsonViewComponent = toRaw(props.jsonViewComponent)

    const managerSize = ref(0)

    const selectedMacro = ref(null)
    provide('selectedMacro', selectedMacro)
    const macroValues = ref({})
    provide('macroValues', macroValues)

    const watchVisible = inject<Ref<boolean>>('collapsedLeftColumnStore')
    const collapsedLeftColumnStore = computed(() => unref(watchVisible) || ref(getLocalStorage('dataGrid_collapsedLeftColumn', false)))

    function getInitialManagerSize() {
      const width = getLocalStorage('dataGridManagerWidth')
      if (isNumber(width) && width > 30 && width < 500) {
        return `${width}px`;
      }
      return '300px';
    }

    const isFormView = computed(() => !!(formDisplay.value && formDisplay.value.config && formDisplay.value.config.isFormView))
    const isJsonView = computed(() => !!(config.value && config.value?.isJsonView))
    const columnsShow = computed(() => !freeTableColumn.value || isDynamicStructure.value && isFormView.value)
    const jsonFiltersSkip = computed(() => !isDynamicStructure.value || !(display.value && display.value?.filterable))
    const selectedCellsPublished = ref(() => [])

    const handleExecuteMacro = () => {
      emit('runMacro', () => (selectedMacro.value, extractMacroValuesForMacro(macroValues.value, selectedMacro.value), selectedCellsPublished.value()))
      selectedMacro.value = null
    }

    function dispatchSize(size) {
      console.log(`bind:size={managerSize}`, size)
      managerSize.value = size
    }

    watch(managerSize, () => setLocalStorage('dataGridManagerWidth', managerSize.value))


    return {
      domColumnManager,
      gridCoreComponent,
      formViewComponent,
      jsonViewComponent,
      managerSize,
      getInitialManagerSize,
      dispatchSize,
      collapsedLeftColumnStore,
      handleExecuteMacro,
      columnsShow,
      jsonFiltersSkip,
      freeTableColumn,
      display,
      selectedMacro,
      formDisplay,
      isFormView,
      isJsonView,
      isDynamicStructure,
      useEvalFilters,
    }
  }
})
</script>

<style scoped>

</style>
