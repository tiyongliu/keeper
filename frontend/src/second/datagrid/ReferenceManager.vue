<template>
  <SearchBoxWrapper>
    <SearchInput placeholder="Search references" v-model:value="filter"/>
    <CloseSearchButton :filter="filter" @filter="filter = ''"/>
  </SearchBoxWrapper>
  <ManagerInnerContainer :width="managerSize">
    <template v-if="foreignKeys.length > 0">
      <div class="bold nowrap ml-1">References tables ({{ foreignKeys.length }}})</div>
      <div
        class="link"
        v-for="fk in foreignKeys.filter(f => filterName(filter, f.refTableName))"
        @click="() => handleRoreignKeys(fk)">
        <FontIcon icon="img link"/>
        <div class="ml-1 nowrap">
          {{ fk.refTableName }}
          ({{ fk.columns.map(x => x.columnName).join(', ') }})
        </div>
      </div>
    </template>

    <template v-if="dependencies.length > 0">
      <div class="bold nowrap ml-1">Dependend tables ({{ dependencies.length }})</div>
      <div
        class="link"
        v-for="fk in dependencies.filter(f => filterName(filter, f.pureName))"
        @click="() => handleDependencies(fk)">
        <FontIcon icon="img reference"/>
        <div class="ml-1 nowrap">
          {{ fk.refTableName }}
          ({{ fk.columns.map(x => x.columnName).join(', ') }})
        </div>
      </div>
    </template>
  </ManagerInnerContainer>
</template>

<script lang="ts">
import {computed, defineComponent, PropType, ref, toRefs} from 'vue'
import {GridDisplay} from '/@/second/keeper-datalib'
import SearchBoxWrapper from '/@/second/widgets/SearchBoxWrapper.vue'
import SearchInput from '/@/second/elements/SearchInput.vue'
import CloseSearchButton from '/@/second/buttons/CloseSearchButton'
import ManagerInnerContainer from '/@/second/elements/ManagerInnerContainer.vue'

import {filterName} from '/@/second/keeper-tools'

export default defineComponent({
  name: 'ReferenceManager',
  components: {
    SearchBoxWrapper,
    SearchInput,
    CloseSearchButton,
    ManagerInnerContainer,
  },
  props: {
    managerSize: {
      type: Number as PropType<number>,
    },
    display: {
      type: Object as PropType<GridDisplay>
    },

    referenceClick: {
      type: Function as PropType<(value: any) => void>,
      default: () => {
      }
    }
  },
  setup(props) {
    const {display, managerSize, referenceClick} = toRefs(props)
    const filter = ref('')
    const foreignKeys = computed(() => display.value?.baseTable?.foreignKeys || [])
    const dependencies = computed(() => display.value?.baseTable?.dependencies || [])

    function handleRoreignKeys(fk) {
      referenceClick.value && referenceClick.value({
        schemaName: fk.refSchemaName,
        pureName: fk.refTableName,
        columns: fk.columns.map(col => ({
          baseName: col.columnName,
          refName: col.refColumnName,
        })),
      })
    }

    function handleDependencies(fk) {
      referenceClick.value && referenceClick.value({
        schemaName: fk.schemaName,
        pureName: fk.pureName,
        columns: fk.columns.map(col => ({
          baseName: col.refColumnName,
          refName: col.columnName,
        })),
      })
    }

    return {
      managerSize,
      filter,
      foreignKeys,
      dependencies,
      filterName,
      handleRoreignKeys,
      handleDependencies,
    }
  }
})
</script>

<style scoped>
.link {
  color: var(--theme-font-link);
  margin: 5px;
  cursor: pointer;
  display: flex;
  flex-wrap: nowrap;
}

.link:hover {
  text-decoration: underline;
}
</style>
