<template>
  <HorizontalSplitter
    :initialValue="getInitialManagerSize()"
    v-model:size="managerSize"
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
              v-bind="pick(Object.assign({}, $props, $attrs), ['managerSize', 'display', 'isJsonView', 'isDynamicStructure', 'conid', 'database'])"
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

          <WidgetColumnBarItem
            title="Filters"
            name="tableFilters"
            height="15%"
            :skip="!display?.filterable || isDynamicStructure || display.filterCount == 0 || isFormView"
            :collapsed="isDetailView">
            <JsonViewFilters
              v-bind="Object.assign({}, $props, $attrs)"
              :managerSize="managerSize"
              :isDynamicStructure="isDynamicStructure"
              :useEvalFilters="useEvalFilters"/>
          </WidgetColumnBarItem>

          <WidgetColumnBarItem
            title="Columns"
            name="freeColumns"
            height="40%"
            :show="freeTableColumn && !isDynamicStructure">
          </WidgetColumnBarItem>

          <WidgetColumnBarItem title="Filters" name="filters" height="30%" :show="isFormView">
            <FormViewFilters
              v-bind="Object.assign({}, $props, $attrs)"
              :managerSize="managerSize"
              :driver="formDisplay?.driver"/>
          </WidgetColumnBarItem>

          <WidgetColumnBarItem
            title="References"
            name="references"
            height="30%"
            :collapsed="isDetailView"
            :show="showReferences && display?.hasReferences">
            <ReferenceManager
              v-bind="Object.assign({}, $props, $attrs)"
              :managerSize="managerSize"/>
          </WidgetColumnBarItem>

          <WidgetColumnBarItem
            title="Macros"
            name="macros"
            :skip="!showMacros"
            :collapsed="!expandMacros">
            <MacroManager v-bind="Object.assign({}, $props, $attrs)" :managerSize="managerSize"/>
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
            :collapsedLeftColumnStore="collapsedLeftColumnStore"
            :formViewAvailable="!!formViewComponent && !!formDisplay"
            :macroValues="extractMacroValuesForMacro(macroValues, selectedMacro)"
            :macroPreview="selectedMacro"
            :changeSelectedColumns="handleChangeSelectedColumns"
          />
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
import {fromPairs, isNumber, mapKeys, pick} from 'lodash-es'
import HorizontalSplitter from '/@/second/elements/HorizontalSplitter.vue'
import WidgetColumnBar from '/@/second/widgets/WidgetColumnBar.vue'
import WidgetColumnBarItem from '/@/second/widgets/WidgetColumnBarItem.vue'
import VerticalSplitter from '/@/second/elements/VerticalSplitter.vue'
import MacroDetail from '/@/second/freetable/MacroDetail.vue'
import ColumnManager from '/@/second/datagrid/ColumnManager.vue'
import JsonViewFilters from '/@/second/jsonview/JsonViewFilters'
import FormViewFilters from '/@/second/formview/FormViewFilters.vue'
import ReferenceManager from '/@/second/datagrid/ReferenceManager.vue'
import MacroManager from '/@/second/freetable/MacroManager.vue'
import {getLocalStorage, setLocalStorage} from '/@/second/utility/storageCache'
import {
  GridConfig,
  GridDisplay,
  MacroDefinition,
  MacroSelectedCell,
  TableFormViewDisplay
} from '/@/second/keeper-datalib'
import {Nullable} from '/@/utils/types'

function extractMacroValuesForMacro(vObject, mObject) {
  // return {};
  const macroValues = unref(vObject)
  const macro = unref(mObject)
  if (!macro) return {}
  return {
    ...fromPairs((macro.args || []).filter(x => x.default != null).map(x => [x.name, x.default])),
    ...mapKeys(macroValues, (_, k) => k.replace(/^.*#/, '')),
  };
}

export default defineComponent({
  name: "DataGrid",
  components: {
    HorizontalSplitter,
    VerticalSplitter,
    MacroDetail,
    WidgetColumnBar,
    WidgetColumnBarItem,
    ColumnManager,
    JsonViewFilters,
    FormViewFilters,
    ReferenceManager,
    MacroManager
  },
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
    setConfig: {
      type: Function as PropType<(target: any) => void>
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
    runMacro: {
      type: Function as PropType<(macro: MacroDefinition, params: {}, cells: MacroSelectedCell[]) => void>
    }
  },
  setup(props) {
    const {
      config,
      formDisplay,
      display,
      freeTableColumn,
      useEvalFilters,
      isDynamicStructure,
      runMacro,
      setConfig,
    } = toRefs(props)

    const domColumnManager = ref<Nullable<{ setSelectedColumns: (value: unknown[]) => void }>>(null)
    const gridCoreComponent = toRaw(props.gridCoreComponent)
    const formViewComponent = toRaw(props.formViewComponent)
    const jsonViewComponent = toRaw(props.jsonViewComponent)

    const managerSize = ref(0)

    const selectedMacro = ref<Nullable<MacroDefinition>>(null)
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
    const isJsonView = computed(() => !!(config.value && config.value['isJsonView']))
    const columnsShow = computed(() => !freeTableColumn.value || isDynamicStructure.value && isFormView.value)
    const jsonFiltersSkip = computed(() => !isDynamicStructure.value || !(display.value && display.value?.filterable))
    const selectedCellsPublished = ref(() => [])

    const handleExecuteMacro = () => {
      runMacro.value && runMacro.value(selectedMacro.value!, extractMacroValuesForMacro(macroValues.value, selectedMacro.value), selectedCellsPublished.value())
      selectedMacro.value = null
    }

    function switchViewEnabled(view) {
      if (view == 'form') return !!formViewComponent && !!formDisplay.value && !isFormView.value && (display.value && display.value?.baseTable?.primaryKey);
      if (view == 'table') return !!(isFormView || isJsonView);
      if (view == 'json') return !!jsonViewComponent && !isJsonView;
    }

    function switchToView(view) {
      if (view == 'form') {
        display.value && display.value.switchToFormView(selectedCellsPublished.value()[0]['rowData'])
      }
      if (view == 'table') {
        setConfig.value && setConfig.value(cfg => ({
          ...cfg,
          isFormView: false,
          isJsonView: false,
          formViewKey: null,
        }));
      }
      if (view == 'json') {
        display.value && display.value.switchToJsonView();
      }
    }

    function handleChangeSelectedColumns(cols) {
      domColumnManager.value && domColumnManager.value.setSelectedColumns(cols)
    }

    watch(managerSize, () => {
      if (managerSize.value) setLocalStorage('dataGridManagerWidth', managerSize.value)
    })

    return {
      domColumnManager,
      gridCoreComponent,
      formViewComponent,
      jsonViewComponent,
      managerSize,
      getInitialManagerSize,
      handleChangeSelectedColumns,
      collapsedLeftColumnStore,
      extractMacroValuesForMacro,
      handleExecuteMacro,
      unref,
      columnsShow,
      jsonFiltersSkip,
      freeTableColumn,
      display,
      macroValues,
      selectedMacro,
      formDisplay,
      isFormView,
      isJsonView,
      isDynamicStructure,
      useEvalFilters,
      pick
    }
  }
})
</script>

<style scoped>
.left {
  display: flex;
  flex: 1;
  background-color: var(--theme-bg-0);
}
</style>
