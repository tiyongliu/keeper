<template>
  <!--
  hidden={$selectedWidget != 'database'}
  -->
  <DatabaseWidget :hidden="false"/>
</template>

<script lang="ts">
  import DatabaseWidget from './DatabaseWidget.vue'
  import { defineComponent, computed, unref, ref, onMounted } from 'vue';
  export default defineComponent({
    components: {
      DatabaseWidget
    }
  })
</script>
