<template>
  <AppObjectList
    :list="sortBy"
  />
</template>

<script lang="ts">
import {computed, defineComponent, PropType, unref} from 'vue'
  // import AppObjectList from './AppObjectList.vue'
  import AppObjectList from './AppObjectList'
import {sortBy} from "lodash-es";
import {filterName} from "/@/packages/tools/src";
  export default defineComponent({
    name: "SubDatabaseList",
    props: {
      passProps: {
        type: Boolean as PropType<boolean>,
        default: false
      },
      data: {
        type: Object as PropType<{}>
      },
      filter: {
        type: String as PropType<string>,
        default: ''
      }
    },
    components: {
      AppObjectList
    },
    setup(props) {
      const {data, filter} = props
      const databases = computed((): {name: string, sortOrder?: string}[] => {
        return [{"name":"information_schema"},{"name":"crmeb_java_beta"},{"name":"mysql"},{"name":"performance_schema"},{"name":"sys"}]
      })

      return {
        databases,
        data,
        filter,
        sortBy: sortBy(
          (unref(databases) || []).filter(x => filterName(unref(filter!), x.name)),
          x => x.sortOrder ?? x.name
        ).map(db => ({ ...db, connection: unref(data) }))
      }
    }
  })
</script>
