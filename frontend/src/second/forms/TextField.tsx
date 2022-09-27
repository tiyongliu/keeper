import {defineComponent, PropType} from 'vue'
export default defineComponent({
  name: 'TextField',
  porps: {
    focused: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    value: {
      type: String as PropType<string>,
    }
  },
  setup() {
    return () => (
      <input
        type="text"
        autocomplete="new-password"
      />
    )
  }
})
