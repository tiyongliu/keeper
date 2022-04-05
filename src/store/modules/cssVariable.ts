import { defineStore } from "pinia";
import { store } from "/@/store";
import { safeJsonParse } from "/@/utils/lib/pkg/stringTools";

interface IVariableBasic {
  selectedWidget: number
  leftPanelWidth: number
  visibleTitleBar: number
}

const LEFTPANELWIDTH = "leftPanelWidth";

export const cssVariableStore = defineStore({
  id: "app-cssVariable",
  state: (): IVariableBasic => ({
    visibleTitleBar: 0,
    selectedWidget: 1,
    leftPanelWidth: getCssVariableCache(300, LEFTPANELWIDTH),
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
    }
  },
  actions: {
    setSelectedWidget(value: number) {
      this.selectedWidget = value;
    },
    setLeftPanelWidth(value) {
      this.leftPanelWidth += value;
      setCssVariableCache(LEFTPANELWIDTH, String(this.leftPanelWidth));
      document.documentElement.style.setProperty("--dim-left-panel-width", `${this.leftPanelWidth}px`);
    },
    subscribeCssVariable(value, transform, cssVariable) {
      document.documentElement.style.setProperty(cssVariable, transform(value));
    },
  }
});

export function useCssVariableStoreWithOut() {
  return cssVariableStore(store);
}

function getCssVariableCache<T>(defaultValue: T, storageName) {
  const init = localStorage.getItem(storageName);
  return (init ? safeJsonParse(init, defaultValue, true) : defaultValue);
}

function setCssVariableCache(storageName, value: string) {
  localStorage.setItem(storageName, value);
}
