<template>
  <ManagerInnerContainer :width="managerSize">
    <SearchBoxWrapper>
      <SearchInput placeholder="Search macros" v-model:value="filter"/>
      <CloseSearchButton :filter="filter" @filter="filter = ''"/>
    </SearchBoxWrapper>
    <AppObjectList
      :list="macrosList"
      :filter="filter"
      :module="macroAppObject"
      :groupFunc="data => data.group"/>
  </ManagerInnerContainer>
</template>

<script lang="ts">
import {defineComponent, PropType, ref, toRefs, computed} from 'vue'
import {sortBy} from 'lodash-es'
import ManagerInnerContainer from '/@/second/elements/ManagerInnerContainer.vue'
import MacroAppObject from '/@/second/appobj/MacroAppObject'
import AppObjectList from '/@/second/appobj/AppObjectList'
import SearchBoxWrapper from '/@/second/widgets/SearchBoxWrapper.vue'
import SearchInput from '/@/second/elements/SearchInput.vue'
import CloseSearchButton from '/@/second/buttons/CloseSearchButton'
import {GridDisplay} from '/@/second/keeper-datalib'
import macros from './macros'
export default defineComponent({
  name: 'MacroManager',
  components: {
    ManagerInnerContainer,
    MacroAppObject,
    AppObjectList,
    SearchBoxWrapper,
    SearchInput,
    CloseSearchButton
  },
  props: {
    managerSize: {
      type: Number as PropType<number>,
    },
    display: {
      type: Object as PropType<GridDisplay>
    },
    macroCondition: {
      type: Function as PropType<(value: any) => void>
    }
  },
  setup(props) {
    const {managerSize, display, macroCondition} = toRefs(props)

    const filter = ref('')

    const macrosList = computed(() =>
      sortBy(macros, 'title').filter(x => (macroCondition.value ? macroCondition.value(x) : true)))

    return {
      filter,
      managerSize,
      display,
      macroCondition,
      macrosList,
      macroAppObject: MacroAppObject
    }
  }
})
</script>
