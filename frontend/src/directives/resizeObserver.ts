import type {App, Directive} from 'vue'
import ResizeObserver from 'resize-observer-polyfill'

const resizeObserver: Directive = {
  // 在绑定元素的 attribute 或事件监听器被应用之前调用, 在指令需要附加须要在普通的 v-on 事件监听器前调用的事件监听器时，这很有用
  created() {
  },
  // 当指令第一次绑定到元素并且在挂载父组件之前调用
  beforeMount() {
  },
  mounted(node, bindings, ...arg) {
    node.resizeObserver = null
    const {value} = bindings
    const observerEnabled = value
    const {props} = arg[0]

    const measure = () => {
      const rect = node.getBoundingClientRect()
      props!.resize && props!.resize({
        detail: {
          width: rect.width,
          height: rect.height,
        },
      })
    }

    node.doUpdate = function () {
      if (observerEnabled && !node.resizeObserver) {
        node.resizeObserver = new ResizeObserver(() => {

        })
        node.resizeObserver.observe(node)
      }
      if (!observerEnabled && node.resizeObserver) {
        node.resizeObserver.disconnect()
        node.resizeObserver = null
      }
    }

    node.doUpdate()
    if (observerEnabled) measure()
  },
  // 在更新包含组件的 VNode 之前调用
  beforeUpdate(node, bindings) {
    node.resizeObserver = null
    const {value} = bindings
    const observerEnabled = value
    node.observerEnabled = observerEnabled
    node.doUpdate()
  },
  // 在包含组件的 VNode 及其子组件的 VNode 更新后调用
  updated() {
  },
  // 在卸载绑定元素的父组件之前调用
  beforeUnmount(node) {
    if (node.resizeObserver) {
      node.resizeObserver.disconnect()
      node.resizeObserver = null
    }
  },
  // 当指令与元素解除绑定且父组件已卸载时, 只调用一次
  unmounted() {

  }
};

export function setupResizeObserver(app: App) {
  app.directive('resizeObserver', resizeObserver)
}

export default resizeObserver

