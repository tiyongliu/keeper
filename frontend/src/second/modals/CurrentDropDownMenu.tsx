import DropDownMenu from './DropDownMenu.vue'
import {defineComponent} from 'vue'
import {dataBaseStore} from "/@/store/modules/dataBase";

export default defineComponent({
  name: 'CurrentDropDownMenu',
  setup() {
    const dataBase = dataBaseStore()
    const currentDropDownMenu = dataBase.$state.currentDropDownMenu
    return () => currentDropDownMenu &&
      <DropDownMenu
        left={currentDropDownMenu.left}
        top={currentDropDownMenu.top}
        item={currentDropDownMenu.items}
        targetElement={currentDropDownMenu.targetElement}
      />
  }
})
