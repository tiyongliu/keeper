import {defineStore} from 'pinia'
import {store} from '/@/store'

interface IVariableBasic {
  leftPanelWidth: number
}

export const VariableBasic = defineStore({
  id: 'app-layout',
  state: (): IVariableBasic => ({
    leftPanelWidth: 300
  }),
  getters: {},
  actions: {
    setLeftPanelWidth(value: number) {
      
    }
  },
})

export function useVariableBasicWithOut() {
  return VariableBasic(store)
}
