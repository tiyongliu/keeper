import {defineComponent} from 'vue'
import DropDownMenu from './DropDownMenu.vue'
import {useBootstrapStore} from "/@/store/modules/bootstrap";

export default defineComponent({
  name: 'CurrentDropDownMenu',
  setup() {
    const bootstrap = useBootstrapStore()
    return () => bootstrap.$state.currentDropDownMenu &&
      <DropDownMenu
        left={bootstrap.$state.currentDropDownMenu.left}
        top={bootstrap.$state.currentDropDownMenu.top}
        item={bootstrap.$state.currentDropDownMenu.items}
        targetElement={bootstrap.$state.currentDropDownMenu.targetElement}
        onClose={() => bootstrap.setCurrentDropDownMenu(null)}
      />
  }
})
