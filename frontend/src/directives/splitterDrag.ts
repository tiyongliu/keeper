import type { App, Directive } from "vue";
const splitterDrag: Directive = {
  // 在绑定元素的 attribute 或事件监听器被应用之前调用, 在指令需要附加须要在普通的 v-on 事件监听器前调用的事件监听器时，这很有用
  created() {},
  // 当指令第一次绑定到元素并且在挂载父组件之前调用
  beforeMount() {},
  mounted(el, bindings, ...arg) {
    let resizeStart = null
    const { value } = bindings
    const axes = value
    const { props } = arg[0]

    el.addEventListener("mousedown", (event) => {
      resizeStart = event[axes]

      document.onmousemove = (e) => {
        e.preventDefault();
        const diff = e[axes] - resizeStart!;
        resizeStart = e[axes];
        props!.resizeSplitter && props!.resizeSplitter({
          detail: diff,
        })
      }

      document.onmouseup = (e) => {
        e.preventDefault();
        resizeStart = null;
        document.onmousemove = null;
        document.onmouseup = null;
      }

    })
  },
  // 在更新包含组件的 VNode 之前调用
  beforeUpdate() {},
  // 在包含组件的 VNode 及其子组件的 VNode 更新后调用
  updated() {},
  // 在卸载绑定元素的父组件之前调用
  beforeUnmount() {},
  // 当指令与元素解除绑定且父组件已卸载时, 只调用一次
  unmounted() {},
};

export function setupSplitterDrag(app: App) {
  app.directive("splitterDrag", splitterDrag);
}

export default splitterDrag;
