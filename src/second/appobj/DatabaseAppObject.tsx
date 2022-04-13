import {defineComponent} from 'vue'

export const extractKey = props => props.name

export default defineComponent({
  name: 'DatabaseAppObject',
  setup() {
    return () => (<div>123</div>)
  }
})
