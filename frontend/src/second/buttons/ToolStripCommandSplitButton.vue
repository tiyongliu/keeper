<template>
  <ToolStripCommandButton
    :command="commands[O]"
    :component="ToolStripSplitDropDownButton"
    :menu="menu"
    :hideDisabled="hideDisabled"
    :buttonLabel="buttonLabel"
  />
</template>

<script lang="ts">
import {defineComponent, PropType, ref, unref, watchEffect, toRefs} from 'vue'
import {compact, isString} from 'lodash-es'
import {ContextMenuItem} from '/@/second/modals/typing'
import ToolStripCommandButton from './ToolStripCommandButton.vue'
import ToolStripSplitDropDownButton from './ToolStripSplitDropDownButton.vue'

export default defineComponent({
  name: "ToolStripCommandSplitButton",
  components: {
    ToolStripCommandButton,
    ToolStripSplitDropDownButton
  },
  props: {
    hideDisabled: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    buttonLabel: {
      type: String as PropType<Nullable<string>>,
    },
    commands: {
      type: Array as PropType<ContextMenuItem[]>,
      //ContextMenuItem[]
    }
  },
  setup(props) {
    const {commands, hideDisabled, buttonLabel} = toRefs(props)
    const menu = ref<any[]>([])

    watchEffect(() => {
      if (unref(commands)) {
        menu.value = compact(unref(commands)).map(command => (isString(unref(command)) ? {command: unref(command)} :  unref(command)))
        console.log(unref(commands), `unref(commands)`)
        console.log(unref(hideDisabled), `unref(hideDisabled)`)
        console.log(unref(buttonLabel), `unref(buttonLabel)`)
      }
    })

    return {
      ToolStripSplitDropDownButton,
      menu,
      commands,
      hideDisabled,
      buttonLabel
    }
  }
})
</script>
