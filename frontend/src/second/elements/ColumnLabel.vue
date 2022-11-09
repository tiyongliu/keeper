<template>
  <span class="label" :class="notNull && 'notNull'">
   <FontIcon v-if="icon" :icon="icon"/>
    {{headerText || columnName}}

    <span v-if="extInfo" class="extinfo">{{ extInfo }}</span>

    <template v-if="showDataType">

      <span class="extinfo" v-if="foreignKey">
        <FontIcon icon="icon arrow-right"/>
         <Link v-if="conid && database" @click.stop="handler($event)">{{
             foreignKey.refTableName
           }}</Link>
        <template v-else>{{ foreignKey.refTableName }}</template>
      </span>

      <span v-else-if="dataType" class="extinfo">{{dataType.toLowerCase()}}</span>
    </template>
  </span>
</template>

<script lang="ts">
import {computed, defineComponent, PropType, toRefs, unref,} from 'vue'
import FontIcon from '/@/second/icons/FontIcon.vue'
import Link from '/@/second/elements/Link.vue'
import {openDatabaseObjectDetail} from '/@/second/appobj/DatabaseObjectAppObject'

export function getColumnIcon(column, forceIcon = false) {
  if (unref(column).autoIncrement) return 'img autoincrement';
  if (unref(column).foreignKey) return 'img foreign-key';
  if (forceIcon) return 'img column';
  return null;
}

export default defineComponent( {
  name: "ColumnLabel",
  components: {
    FontIcon,
    Link
  },
  props: {
    notNull: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    forceIcon: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    headerText: {
      type: String as PropType<string>,
      default: ''
    },
    columnName: {
      type: [String, Object] as PropType<string | object>,
      default: ''
    },
    extInfo: {
      type: String as PropType<string>
    },
    dataType: {
      type: String as PropType<string>
    },
    showDataType: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    foreignKey: {
      type: Object as PropType<{
        refSchemaName: string
        refTableName: string
      }>
    },
    conid: {
      type: String as PropType<string>,
    },
    database: {
      type: String as PropType<string>,
    },
    iconOverride: {
      type: String as PropType<string>,
    }
  },
  setup(props, {attrs}) {
    const {
      notNull,
      forceIcon,
      headerText,
      columnName,
      extInfo,
      dataType,
      showDataType,
      foreignKey,
      conid,
      database,
      iconOverride
    } = toRefs(props)

    const icon = computed(() => iconOverride.value || getColumnIcon(Object.assign({}, props, attrs), forceIcon.value))

    function handler(e: Event) {
      e.stopPropagation()
      openDatabaseObjectDetail('TableDataTab', null, {
        schemaName: foreignKey.value!.refSchemaName,
        pureName: foreignKey.value!.refTableName,
        conid,
        database,
        objectTypeField: 'tables',
      })
    }

    return {
      notNull,
      forceIcon,
      headerText,
      columnName,
      extInfo,
      dataType,
      showDataType,
      foreignKey,
      conid,
      database,
      icon,
      handler
    }
  }
})
</script>

<style scoped>
.label {
  white-space: nowrap;
}

.label.notNull {
  font-weight: bold;
}

.extinfo {
  font-weight: normal;
  margin-left: 5px;
  color: var(--theme-font-3);
}
</style>
