import {defineComponent, PropType, toRefs} from 'vue'
import InlineButton from '/@/second/buttons/InlineButton.vue'
import FontIcon from '/@/second/icons/FontIcon.vue'
export default defineComponent({
  name: 'CloseSearchButton',
  components: {
    InlineButton,
    FontIcon
  },
  props: {
    filter: {
      type: String as PropType<string>,
    },
    showDisabled: {
      type: Boolean as PropType<boolean>,
      default: false
    }
  },
  emit: ['filter'],
  setup(props, {emit}) {
    const {filter, showDisabled} = toRefs(props)

    function handleClick() {
      emit('filter', '')
    }

    return () => (filter.value || showDisabled.value ? (
      <InlineButton
        title="Clear filter"
        disabled={!filter.value}
        onClick={handleClick}>
        <FontIcon icon="icon close" />
      </InlineButton>
    ) : null)
  }
})
