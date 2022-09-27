<template>
  <div class="wrapper">
    <InlineButton>
      <FontIcon icon="icon arrow-left" />
    </InlineButton>
    <span class="label">Start:</span>
    <TextField type="number" />
    <span class="label">Rows:</span>
    <TextField type="number" />
    <InlineButton>
      <FontIcon icon="icon arrow-right" />
    </InlineButton>
  </div>
</template>

<script lang="ts">
import {defineComponent, PropType} from 'vue'
import FontIcon from '/@/second/icons/FontIcon.vue'
import InlineButton from '/@/second/buttons/InlineButton.vue'
import TextField from '/@/second/forms/TextField'
export default defineComponent({
  name: "Pager",
  components: {
    FontIcon,
    InlineButton,
    TextField
  },
  props: {
    skip: {
      type: Number as PropType<number>,
    },
    limit: {
      type: Number as PropType<number>,
    }
  },
  setup() {

  }
})
</script>

<style scoped>
.wrapper :global(input) {
  width: 100px;
}
.wrapper {
  display: flex;
  align-items: center;
}
.label {
  margin-left: 5px;
  margin-right: 5px;
}
</style>
