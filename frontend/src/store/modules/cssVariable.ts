import {defineStore} from "pinia";
import {store} from "/@/store";
import {reactive} from "vue";
import {isNumber} from 'lodash-es'
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

export const LEFTPANELWIDTH = "leftPanelWidth"

export const dynamicProps = reactive({splitterVisible: false})
const _leftPanelWidth = getWithStorageVariableCache(300, LEFTPANELWIDTH)
export const cssVariableStore = defineStore({
  id: "app-cssVariable",
  state: (): IVariableBasic => ({
    currentDropDownMenu: null,
    visibleTitleBar: 0,
    selectedWidget: 1,
    leftPanelWidth: parseFloat(_leftPanelWidth).toString() !== 'NaN' ?
      parseFloat(_leftPanelWidth) : 300,
    dynamicProps: {
      splitterVisible: false
    }
  }),
  getters: {
    getDynamicProps(): { splitterVisible: boolean } {
      return this.dynamicProps
    }
  },
  actions: {
    setSelectedWidget(value: number) {
      this.selectedWidget = value;
    },
    setLeftPanelWidth(value) {
      this.leftPanelWidth += value
      document.documentElement.style.setProperty("--dim-left-panel-width", `${this.leftPanelWidth}px`);
      if (isNumber(this.leftPanelWidth)) {
        setWithStorageVariableCache(LEFTPANELWIDTH, String(this.leftPanelWidth));
      }
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
