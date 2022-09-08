<template>
  <div v-if="isSmall" class="container-small">
    <FontIcon :icon="icon" />
    {{message || 'Unknown error'}}
  </div>
  <div v-else-if="alignTop">
    <div class="container">
      <div class="icon">
        <FontIcon :icon="icon" />
      </div>
      {{message || 'Unknown error'}}
    </div>
  </div>
  <div v-else class="container">
    <div class="icon">
      <FontIcon :icon="icon" />
    </div>
    {{message || 'Unknown error'}}
  </div>
</template>

<script lang="ts">
  import {defineComponent, toRefs, PropType} from 'vue'
  import FontIcon from '/@/second/icons/FontIcon.vue'

  export default defineComponent({
    name: "ErrorInfo",
    props: {
      message: {
        type: String as PropType<string>,
      },
      icon: {
        type: String as PropType<string>,
        default: 'img error'
      },
      isSmall: {
        type: Boolean as PropType<false>,
        default: false
      },
      alignTop: {
        type: Boolean as PropType<false>,
        default: false
      },
    },
    components: {
      FontIcon
    },
    setup(props) {
      return {
        ...toRefs(props)
      }
    }
  })
</script>

<style scoped>
  .container {
    display: flex;
    padding-right: 10px;
    align-items: center;
  }
  .icon {
    font-size: 20pt;
    margin: 10px;
  }

  .container-small {
    display: flex;
    margin-right: 10px;
  }
</style>

