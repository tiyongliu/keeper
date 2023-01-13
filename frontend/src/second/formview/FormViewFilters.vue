<template>
  <div class="m-1">
    <div>Column filter</div>
    <div class="flex">
      <a-input
        size="small"
        v-model:value="formColumnFilterText"
        @keydown="handleKeydown"
        @input="handleInput"
      />
    </div>

    <ManagerInnerContainer v-if="baseTable && baseTable?.primaryKey" :width="managerSize">
      <PrimaryKeyFilterEditor
        v-for="col in baseTable.primaryKey.columns"
        :column="col"
        :baseTable="baseTable"
        :formDisplay="formDisplay"
      />
      <FormViewFilterColumn
        v-for="uniqueName in allFilterNames"
        :column="formDisplay ? formDisplay.columns.find(x => x.uniqueName == uniqueName) : null"
        :formDisplay="formDisplay"
        :filters="filters"
        :driver="driver"
        :conid="conid"
        :database="database"
        :schemaName="schemaName"
        :pureName="pureName"
      />
    </ManagerInnerContainer>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, PropType, ref, toRefs} from 'vue'
import {keys, union} from 'lodash-es'
import ManagerInnerContainer from '/@/second/elements/ManagerInnerContainer.vue'
import PrimaryKeyFilterEditor from '/@/second/formview/PrimaryKeyFilterEditor'
import FormViewFilterColumn from './FormViewFilterColumn'
import {TableFormViewDisplay} from '/@/second/keeper-datalib'
import {EngineDriver} from '/@/second/keeper-types'
import keycodes from '/@/second/utility/keycodes'

export default defineComponent({
  name: 'FormViewFilters',
  components: {
    ManagerInnerContainer,
    PrimaryKeyFilterEditor,
    FormViewFilterColumn,
  },
  props: {
    managerSize: {
      type: Number as PropType<number>,
    },
    formDisplay: {
      type: Object as PropType<TableFormViewDisplay>
    },
    setConfig: {
      type: Function as PropType<(target: any) => void>
    },
    driver: {
      type: Object as PropType<EngineDriver>
    },
    conid: {
      type: String as PropType<string>
    },
    database: {
      type: String as PropType<string>
    },
    schemaName: {
      type: String as PropType<string>
    },
    pureName: {
      type: String as PropType<string>
    }
  },
  setup(props) {
    const {formDisplay, managerSize, setConfig} = toRefs(props)
    const formColumnFilterText = ref(formDisplay.value?.config?.formColumnFilterText || '')

    const baseTable = computed(() => formDisplay.value
      ? formDisplay.value?.baseTable : null)

    const formFilterColumns = computed(() => formDisplay.value
      ? formDisplay.value?.config?.formFilterColumns : null)

    const filters = computed(() => formDisplay.value
      ? formDisplay.value?.config?.filters : null)

    const allFilterNames = computed(() => union(keys(filters.value || {}), formFilterColumns.value || []))

    function handleKeydown(e) {
      if (e.keyCode == keycodes.escape) {
        setConfig.value && setConfig.value(x => ({
          ...x,
          formColumnFilterText: '',
        }))
      }
    }

    function handleInput(e) {
      if (e.keyCode == keycodes.escape) {
        setConfig.value && setConfig.value(x => ({
          ...x,
          formColumnFilterText: formColumnFilterText.value
        }))
      }
    }

    return {
      formColumnFilterText,
      managerSize,
      baseTable,
      formFilterColumns,
      filters,
      allFilterNames,
      handleKeydown,
      handleInput
    }
  }
})
</script>
