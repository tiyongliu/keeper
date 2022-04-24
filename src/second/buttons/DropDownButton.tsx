import {defineComponent, PropType, toRefs} from 'vue';
import InlineButton from '/@/second/buttons/InlineButton.vue'
import FontIcon from '/@/second/icons/FontIcon.vue'
import {cssVariableStore} from '/@/store/modules/cssVariable'

export default defineComponent({
  name: 'DropDownButton',
  props: {
    icon: {
      type: String as PropType<string>,
      default: 'icon chevron-down'
    },
    menu: {
      type: [Function, Array] as PropType<[] | Function>,
    },
    narrow: {
      type: Boolean as PropType<false>,
      default: false
    },
  },
  setup(props) {
    const cssVariable = cssVariableStore()
    const {narrow, icon} = toRefs(props)

    function handleClick() {
      /*
       const rect = domButton.getBoundingClientRect();
       const left = rect.left;
       const top = rect.bottom;
       currentDropDownMenu.set({ left, top, items: menu });
      * */
      cssVariable.subscribeCurrentDropDownMenu()
    }

    return () => (
      <InlineButton square narrow={narrow} onClick={handleClick}>
        <FontIcon icon={icon}/>
      </InlineButton>
    )
  }
})
