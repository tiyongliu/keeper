import {defineStore} from "pinia";
import {store} from "/@/store";
import {reactive} from "vue";
import {getWithStorageVariableCache, setWithStorageVariableCache} from '../index'

interface IVariableBasic {
  currentDropDownMenu: null
  selectedWidget: number
  leftPanelWidth: number
  visibleTitleBar: number
  dynamicProps: {
    splitterVisible: boolean
  }
}

const LEFTPANELWIDTH = "leftPanelWidth"

export const dynamicProps = reactive({splitterVisible: false})

export const cssVariableStore = defineStore({
  id: "app-cssVariable",
  state: (): IVariableBasic => ({
    currentDropDownMenu: null,
    visibleTitleBar: 0,
    selectedWidget: 1,
    leftPanelWidth: getWithStorageVariableCache(300, LEFTPANELWIDTH),
    dynamicProps: {
      splitterVisible: false
    }
  }),
  getters: {
    getSelectedWidget(): number {
      return this.selectedWidget;
    },
    getLeftPanelWidth(): number {
      return this.leftPanelWidth;
    },
    getVisibleTitleBar(): number {
      return this.visibleTitleBar;
    },
    getDynamicProps(): { splitterVisible: boolean } {
      return this.dynamicProps
    }
  },
  actions: {
    setSelectedWidget(value: number) {
      this.selectedWidget = value;
    },
    setLeftPanelWidth(value) {
      this.leftPanelWidth += value;
      setWithStorageVariableCache(LEFTPANELWIDTH, String(this.leftPanelWidth));
      document.documentElement.style.setProperty("--dim-left-panel-width", `${this.leftPanelWidth}px`);
    },
    subscribeCssVariable(value, transform, cssVariable) {
      document.documentElement.style.setProperty(cssVariable, transform(value));
    },
    subscribeDynamicProps(value: { splitterVisible: boolean }) {
      this.dynamicProps = value
    },
    subscribeCurrentDropDownMenu() {

    }
  }
})

export function useCssVariableStoreWithOut() {
  return cssVariableStore(store);
}