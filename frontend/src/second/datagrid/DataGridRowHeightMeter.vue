<template>
  <div>
    <table>
      <tr>
        <td>
          title
          <InlineButton square>
            <FontIcon icon="icon chevron-down"/>
          </InlineButton>
        </td>
      </tr>
    </table>
  </div>
</template>

<script lang="ts">
import {defineComponent, ref} from 'vue'
import InlineButton from '/@/second/buttons/InlineButton.vue'
import FontIcon from '/@/second/icons/FontIcon.vue'

export const dataGridRowHeight = ref(0)

export default defineComponent({
  name: 'DataGridRowHeightMeter',
  components: {
    InlineButton,
    FontIcon,
  },
  setup() {

  }
})
</script>

<style scoped>
div {
  position: fixed;
  left: -1000px;
  top: -1000px;
  visibility: hidden;
}

td {
  display: flex;
}
</style>
