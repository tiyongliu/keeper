import {defineComponent} from 'vue'
import DropDownMenu from './DropDownMenu.vue'
import {dataBaseStore} from "/@/store/modules/dataBase";

export default defineComponent({
  name: 'CurrentDropDownMenu',
  setup() {
    const dataBase = dataBaseStore()
    return () => dataBase.$state.currentDropDownMenu &&
      <DropDownMenu
        left={dataBase.$state.currentDropDownMenu.left}
        top={dataBase.$state.currentDropDownMenu.top}
        item={dataBase.$state.currentDropDownMenu.items}
        targetElement={dataBase.$state.currentDropDownMenu.targetElement}
        onClose={() => dataBase.subscribeCurrentDropDownMenu(null)}
      />
  }
})
