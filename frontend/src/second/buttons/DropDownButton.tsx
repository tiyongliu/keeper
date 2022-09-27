import {defineComponent, nextTick, PropType, ref, toRefs, unref} from 'vue';
import InlineButton from '/@/second/buttons/InlineButton.vue'
import FontIcon from '/@/second/icons/FontIcon.vue'
import {useLocaleStore} from '/@/store/modules/locale'
import {useBootstrapStore} from "/@/store/modules/bootstrap";

export default defineComponent({
  name: 'DropDownButton',
  props: {
    icon: {
      type: String as PropType<string>,
      default: 'icon chevron-down'
    },
    menu: {
      type: Array as unknown as PropType<[]>,
    },
    narrow: {
      type: Boolean as PropType<boolean>,
      default: false
    },
  },
  setup(props) {
    const locale = useLocaleStore()
    const bootstrap = useBootstrapStore()
    const {narrow, icon, menu} = toRefs(props)
    const domButton = ref<Nullable<HTMLElement>>(null)

    async function handleClick() {
      await nextTick()
      const rect = domButton.value!.getBoundingClientRect();
      const left = rect.left;
      const top = rect.bottom;
      bootstrap.subscribeCurrentDropDownMenu({left, top, items: menu.value!});

      locale.subscribeCurrentDropDownMenu()
    }

    return () => (
      <InlineButton square narrow={unref(narrow)} onClick={handleClick} ref={domButton}>
        <FontIcon icon={unref(icon)}/>
      </InlineButton>
    )
  }
})
