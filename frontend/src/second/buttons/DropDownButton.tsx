import {defineComponent, PropType, ref, toRefs, unref} from 'vue'
import InlineButton from '/@/second/buttons/InlineButton.vue'
import FontIcon from '/@/second/icons/FontIcon.vue'
import {ContextMenuItem} from '/@/components/ContextMenu'
import {useContextMenu} from '/@/hooks/web/useContextMenu'
export default defineComponent({
  name: 'DropDownButton',
  props: {
    icon: {
      type: String as PropType<string>,
      default: 'icon chevron-down'
    },
    menu: {
      type: [Function, Array] as unknown as PropType<() => ContextMenuItem[]>,
    },
    narrow: {
      type: Boolean as PropType<boolean>,
      default: false
    },
  },
  setup(props) {
    const {narrow, icon} = toRefs(props)
    const {menu} = props
    const domButton = ref<Nullable<HTMLElement>>(null)

    const [createContextMenu] = useContextMenu()
    async function handleClick(e) {
      e.preventDefault()
      e.stopPropagation()
      // @ts-ignore
      createContextMenu({event: e, items: menu})
    }

    return () => (
      <InlineButton square narrow={unref(narrow)} onClick={e => handleClick(e)} ref={domButton}>
        <FontIcon icon={unref(icon)}/>
      </InlineButton>
    )
  }
})
