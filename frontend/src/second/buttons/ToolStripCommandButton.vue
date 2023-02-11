<template>
  <template v-if="visible">
    <component
      :is="component"
      :title="title"
      :icon="cmd.icon"
      :disabled="!cmd.enabled"
      v-bind="$attrs">
      {{ buttonLabel || cmd.toolbarName || cmd.name }}
    </component>
  </template>
</template>

<script lang="ts">
import {Component, computed, defineComponent, PropType, toRaw, toRef, toRefs} from 'vue'
import {storeToRefs} from 'pinia'
import {useBootstrapStore} from '/@/store/modules/bootstrap'
import {formatKeyText} from '/@/second/utility/common'
import ToolStripButton from './ToolStripButton.vue'

function getCommandTitle(command) {
  let res = command.text
  if (command.keyText || command.keyTextFromGroup) {
    res += ` (${formatKeyText(command.keyText || command.keyTextFromGroup)})`
  }
  return res
}

export default defineComponent({
  name: 'ToolStripCommandButton',
  props: {
    hideDisabled: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    component: {
      type: [String, Object] as PropType<string | Component | any>,
      default: ToolStripButton
    },
    command: {
      type: String as PropType<string>
    },
    buttonLabel: {
      type: String as PropType<string>
    }
  },
  setup(props) {
    const command = toRef(props, 'command')
    const hideDisabled = toRef(props, 'hideDisabled')
    const component = toRaw(props.component)

    const bootstrap = useBootstrapStore()
    const {getCommandsCustomized} = storeToRefs(bootstrap)
    const cmd = computed<{
      enabled: boolean
      toolbarName: string
      name: string
    }>(() => Object.values(getCommandsCustomized.value).find((x: any) => x.id == command.value) as any)

    setTimeout(() => {
      console.log(Object.values(getCommandsCustomized.value))
    }, 5000)

    const visible = computed(() => cmd.value && (!hideDisabled.value || cmd.value.enabled))
    return {
      ...toRefs(props),
      cmd,
      getCommandsCustomized,
      title: computed(() => getCommandTitle(cmd.value)),
      component,
      visible
    }
  }
})
</script>
