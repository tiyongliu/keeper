<template>
  <div class="flex">
    <a-input
      class="dataFilter-input"
      ref="domInput"
      size="small"
      autocomplete="off"
      :readOnly="isReadOnly"
      v-model:value="value"
      :class="[isError && 'isError', isOk && 'isOk']"
      @keydown="handleKeyDown"
      @blur="applyFilter"
      @paste="handlePaste"
      placeholder="Filter"
    />
    <InlineButton v-if="customCommandIcon" :title="customCommandTooltip">
      <FontIcon :icon="customCommandIcon"/>
    </InlineButton>
    <template v-if="conid && database && driver">
      <InlineButton
        v-if="driver?.databaseEngineTypes?.includes('sql') && foreignKey"
        narrow square>
        <FontIcon icon="icon dots-horizontal" />
      </InlineButton>
      <InlineButton
        v-else-if="(pureName && columnName) ||
        (pureName && uniqueName && driver?.databaseEngineTypes?.includes('document'))"
        narrow square>
        <FontIcon icon="icon dots-vertical" />
      </InlineButton>
    </template>
    <template v-else-if="jslid">
      <InlineButton narrow square>
        <FontIcon icon="icon dots-vertical" />
      </InlineButton>
    </template>
    <DropDownButton icon="icon filter" narrow />
    <div
      v-if="showResizeSplitter"
      class="horizontal-split-handle resizeHandleControl"
      v-splitterDrag="'clientX'"
      :resizeSplitter="(e) => dispatchResizeSplitter(e)"/>
  </div>
</template>

<script lang="ts">
import {defineComponent, PropType, ref, toRef, toRefs, watch} from 'vue'
import {Input} from 'ant-design-vue'
import InlineButton from '/@/second/buttons/InlineButton.vue'
import FontIcon from '/@/second/icons/FontIcon.vue'
import {EngineDriver} from '/@/second/keeper-types'
import {FilterType} from '/@/second/keeper-filterparser'
import {parseFilter, createMultiLineFilter} from '/@/second/keeper-filterparser'
import keycodes from '/@/second/utility/keycodes'
import DropDownButton from '/@/second/buttons/DropDownButton'
export default defineComponent({
  name: 'DataFilterControl',
  components: {
    [Input.name]: Input,
    InlineButton,
    FontIcon,
    DropDownButton,
  },
  props: {
    isReadOnly: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    filterType: {
      type: String as PropType<FilterType>,
    },
    filter: {
      type: String as PropType<string>,
    },
    setFilter: {
      type: Function as PropType<(value: any) => void>
    },
    foreignKey: {
      type: Object as PropType<{ refTableName: string }>
    },
    conid: {
      type: String as PropType<string>
    },
    database: {
      type: String as PropType<string>
    },
    driver: {
      type: Object as PropType<EngineDriver>
    },
    jslid: {
      type: [String, Number] as PropType<string | number>
    },
    pureName: {
      type: String as PropType<string>
    },
    schemaName: {
      type: String as PropType<string>
    },
    columnName: {
      type: String as PropType<string>
    },
    uniqueName: {
      type: String as PropType<string>
    },
    customCommandIcon: {
      type: String as PropType<string>
    },
    customCommandTooltip: {
      type: String as PropType<string>
    },
    showResizeSplitter: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    FocusGrid: {
      type: Function as PropType<() => void>
    }
  },
  emits: ['dispatchResizeSplitter', 'update:filter'],
  setup(props, {emit}) {
    const filterType = toRef(props, 'filterType')
    const isReadOnly = toRef(props, 'isReadOnly')
    const setFilter = toRef(props, 'setFilter')
    const filter = toRef(props, 'filter')
    const FocusGrid = toRef(props, 'FocusGrid')

    const domInput = ref<Nullable<HTMLElement>>(null)
    const value = ref<Nullable<string>>(null)
    const isOk = ref<boolean>(false)
    const isError = ref<boolean>(false)

    function dispatchResizeSplitter(e) {
      emit('dispatchResizeSplitter', e)
    }

    const handleKeyDown = ev => {
      if (isReadOnly.value) return
      if (ev.keyCode == keycodes.enter) {
        applyFilter()
      }
      if (ev.keyCode == keycodes.escape) {
        setFilter.value!('')
      }
      if (ev.keyCode == keycodes.downArrow) {
        if (FocusGrid.value) FocusGrid.value()
        // ev.stopPropagation();
        ev.preventDefault();
      }
      // if (ev.keyCode == KeyCodes.DownArrow || ev.keyCode == KeyCodes.UpArrow) {
      //     if (this.props.onControlKey) this.props.onControlKey(ev.keyCode);
      // }
    }

    function handlePaste(event) {
      var pastedText: string | undefined = undefined
      // @ts-ignore
      if (window.clipboardData && window.clipboardData.getData) {
        // IE
        // @ts-ignore
        pastedText = window.clipboardData.getData('Text')
      } else if (event.clipboardData && event.clipboardData.getData) {
        pastedText = event.clipboardData.getData('text/plain');
      }
      if (pastedText && pastedText.includes('\n')) {
        event.preventDefault();
        setFilter.value!(createMultiLineFilter('is', pastedText));
      }
    }

    watch(() => filter.value, () => {
      value.value = <string>filter.value
    }, {immediate: true})

    watch(() => value.value, () => {
      try {
        isOk.value = false
        isError.value = false
        if (value.value) {
          parseFilter(value.value, filterType.value!)
          isOk.value = true
        }
      } catch (e) {
        // console.error(err)
        isError.value = true
      }
    }, {immediate: true})

    function applyFilter() {
      if ((filter.value || '') == (value.value || '')) return
      setFilter.value!(value.value)
    }

    return {
      domInput,
      value,
      isOk,
      isError,
      ...toRefs(props),
      dispatchResizeSplitter,
      handleKeyDown,
      applyFilter,
      handlePaste,
    }
  }
})
</script>

<style scoped>
.dataFilter-input {
  flex: 1;
  min-width: 10px;
  width: 1px;
}

.dataFilter-input.isError {
  background-color: var(--theme-bg-red);
}

.dataFilter-input.isOk {
  background-color: var(--theme-bg-green);
}
</style>
