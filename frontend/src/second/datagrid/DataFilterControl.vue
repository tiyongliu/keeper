<template>
  <div class="flex"></div>
</template>

<script lang="ts">
import {defineComponent, PropType} from 'vue'
import InlineButton from '/@/second/buttons/InlineButton.vue'
import FontIcon from '/@/second/icons/FontIcon.vue'
export default defineComponent({
  name: 'DataFilterControl',
  components: {
    InlineButton,
    FontIcon,
  },
  props: {
    isReadOnly: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    conid: {
      type: String as PropType<string>
    },
    database: {
      type: String as PropType<string>
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
  },
  setup() {

  }
})
</script>

<style scoped>
input {
  flex: 1;
  min-width: 10px;
  width: 1px;
}

input.isError {
  background-color: var(--theme-bg-red);
}

input.isOk {
  background-color: var(--theme-bg-green);
}
</style>
