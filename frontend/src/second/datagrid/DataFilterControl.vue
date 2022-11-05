<template>
  <div class="flex">
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
  </div>
</template>

<script lang="ts">
import {defineComponent, PropType, toRefs} from 'vue'
import InlineButton from '/@/second/buttons/InlineButton.vue'
import FontIcon from '/@/second/icons/FontIcon.vue'
import {EngineDriver} from '/@/second/keeper-types'
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
    }
  },
  setup(props) {
    return {
      ...toRefs(props),
    }
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
