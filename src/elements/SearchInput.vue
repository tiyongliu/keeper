<template>
  <input type="text" :placeholder="placeholder"/>
</template>

<script>
    export default {
        name: "SearchInput"
    }
</script>

<script lang="ts">
  import { defineComponent, computed, unref, ref, onMounted } from 'vue';
  export default defineComponent({
    name: "SearchInput",
    props: {
      placeholder: {
        type: [String] as PropType<string>,
      }
    },
    components: {

    },
    setup(props) {
      return {
        placeholder: props.placeholder
      }
    }
  })
</script>

<style scoped>
  input {
    flex: 1;
    min-width: 10px;
    width: 10px;
    border: none;
  }
</style>
