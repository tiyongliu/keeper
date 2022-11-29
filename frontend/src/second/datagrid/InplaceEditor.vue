<template>
  <input
    type="text"
    ref="domEditor"
    @change="handleChange"
    @keydown="handleKeyDown"
    @blur="handleBlur"
    v-model="text"
    :style="realWidth ? `width:${realWidth}px;min-width:${realWidth}px;max-width:${realWidth}px` : undefined"
    :class="[fillParent && 'fillParent', showEditorButton && 'showEditorButton']"
  />
  <ShowFormButton v-if="showEditorButton" icon="icon edit"/>
</template>

<script lang="ts">
import {computed, defineComponent, onBeforeUnmount, onMounted, PropType, ref, toRefs} from 'vue'
import {Input} from 'ant-design-vue'
import ShowFormButton from '/@/second/formview/ShowFormButton.vue'
import keycodes from '/@/second/utility/keycodes'
import {parseCellValue, stringifyCellValue} from '/@/second/keeper-tools'
import {isCtrlOrCommandKey} from '/@/second/utility/common';
import createRef from '/@/second/utility/createRef'

export default defineComponent({
  name: "InplaceEditor",
  components: {
    ShowFormButton,
    [Input.name]: Input
  },
  props: {
    width: {
      type: Number as PropType<number>,
    },
    fillParent: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    inplaceEditorState: {
      type: Object as PropType<{ [key in string]: unknown }>,
    },
    dispatchInsplaceEditor: {
      type: Function as PropType<(action: any) => void>
    },
    cellValue: {
      type: [String, Boolean, Number] as PropType<boolean | string | number>
    }
  },
  emits: ['setValue'],
  setup(props, {emit}) {
    const {width, dispatchInsplaceEditor, cellValue, inplaceEditorState} = toRefs(props)

    const showEditorButton = ref(true)
    const text = ref<Nullable<unknown>>(null)
    const domEditor = ref<Nullable<HTMLElement>>(null)

    const realWidth = computed(() => width.value ? width.value - (showEditorButton.value ? 16 : 0) : undefined)

    const isChangedRef = createRef(!!(inplaceEditorState.value!.text))

    function handleBlur() {
      if (isChangedRef.get()) {
        emit('setValue', parseCellValue(text.value))
        isChangedRef.set(false);
      }
      dispatchInsplaceEditor.value!({type: 'close'});
    }

    onMounted(() => {
      text.value = (inplaceEditorState.value!.text) || stringifyCellValue(cellValue.value)
      domEditor.value!.focus()

      // domEditor.value?.value = (inplaceEditorState.value!.text) || stringifyCellValue(cellValue.value)
      // domEditor.value!.focus()
      // if (inplaceEditorState.value!.selectAll) {
      //   domEditor.value!.select();
      // }
    })

    onBeforeUnmount(() => {
      text.value = null
    })

    function handleChange() {
      isChangedRef.set(true);
      showEditorButton.value = false
    }

    function handleKeyDown(event) {
      showEditorButton.value = false;

      switch (event.keyCode) {
        case keycodes.escape:
          isChangedRef.set(false);
          dispatchInsplaceEditor.value!({type: 'close'});
          break;
        case keycodes.enter:
          if (isChangedRef.get()) {
            emit('setValue', parseCellValue(text.value))
            isChangedRef.set(false);
          }
          domEditor.value!.blur();
          event.preventDefault();
          dispatchInsplaceEditor.value!({type: 'close', mode: 'enter'});
          break;
        case keycodes.tab:
          if (isChangedRef.get()) {
            emit('setValue', parseCellValue(text.value))
            isChangedRef.set(false);
          }
          domEditor.value!.blur();
          event.preventDefault();
          dispatchInsplaceEditor.value!({type: 'close', mode: event.shiftKey ? 'shiftTab' : 'tab'});
          break;
        case keycodes.s:
          if (isCtrlOrCommandKey(event)) {
            if (isChangedRef.get()) {
              emit('setValue', parseCellValue(text.value))
              isChangedRef.set(false);
            }
            event.preventDefault();
            dispatchInsplaceEditor.value!({type: 'close', mode: 'save'});
          }
          break;
      }
    }

    return {
      showEditorButton,
      domEditor,
      text,
      realWidth,
      handleChange,
      handleKeyDown,
      handleBlur,
    }
  }
})
</script>

<style scoped>
input {
  border: 0px solid;
  outline: none;
  margin: 0px;
  padding: 0px;
}

input.fillParent {
  position: absolute;
  left: 0;
  top: 0;
  right: 0;
  bottom: 0;
  margin: auto;
}

input.showEditorButton {
  margin-right: 16px;
}
</style>
