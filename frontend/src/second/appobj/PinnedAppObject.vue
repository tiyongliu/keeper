<template>
  <template v-if="data">
    <DatabaseObjectAppObject
      v-if="data && data?.objectTypeField"
      v-bind="Object.assign({}, $props, $attrs)"/>
    <DatabaseAppObject
      v-else
      v-bind="Object.assign({}, $props, $attrs)"/>
  </template>
</template>

<script lang="ts">
import {defineComponent, PropType, toRef} from 'vue'
import DatabaseAppObject from './DatabaseAppObject'
import DatabaseObjectAppObject from './DatabaseObjectAppObject'
import {IPinnedDatabasesItem} from "/@/second/typings/types/standard.d";

export default defineComponent({
  name: 'PinnedAppObject',
  props: {
    data: {
      type: Object as PropType<IPinnedDatabasesItem>,
    },
  },
  components: {
    DatabaseAppObject,
    DatabaseObjectAppObject
  },
  setup(props) {
    return {
      data: toRef(props, 'data')
    }
  }
})
</script>
