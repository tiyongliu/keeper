<template>

</template>

<script lang="ts">
  import {defineComponent, PropType} from 'vue'
  import AppObjectList from './AppObjectList.vue'
  export default defineComponent({
    name: "SubDatabaseList",
    props: {
      filter: {
        type: String as unknown as PropType<string>,
      }
    },
    components: {
      AppObjectList
    },
    setup() {

    }
  })
</script>
