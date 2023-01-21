<template>
  <div v-if="wrapper" class="wrapper">
    <div class="box">
      <div class="container">
        <div class="spinner">
          <FontIcon icon="icon loading" />
        </div>
        {{message}}
      </div>
    </div>
  </div>
  <div v-else class="container">
    <div class="spinner">
      <FontIcon icon="icon loading" />
    </div>
    {{message}}
  </div>
</template>

<script lang="ts">
  import {defineComponent, toRefs, PropType} from 'vue'
  import FontIcon from '/@/second/icons/FontIcon.vue'

  export default defineComponent({
    name: "LoadingInfo",
    props: {
      message: {
        type: String as PropType<string>,
      },
      wrapper: {
        type: Boolean as PropType<boolean>,
        default: false
      }
    },
    components: {FontIcon},
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
    align-items: center;
    margin-right: 10px;
  }

  .spinner {
    font-size: 20pt;
    margin: 10px;
  }
  .wrapper {
    position: absolute;
    left: 0;
    top: 0;
    right: 0;
    bottom: 0;
    display: flex;
    align-items: center;
    justify-content: space-around;
  }
  .box {
    background-color: var(--theme-bg-2);
    padding: 8px 10px;
    border: 1px solid var(--theme-border);
  }
</style>

