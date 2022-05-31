import {defineComponent, onMounted, onBeforeUnmount, provide, inject} from 'vue';
import keycodes from '/@/second/utility/keycodes'
import createRef from '/@/second/utility/createRef'

const contextKey = 'formProviderContextKey'

export function getFormContext(): any {
  return inject(contextKey)
}

export default defineComponent({
  name: 'FormProviderCore',
  setup(_, {slots}) {
    const handleEnter = (e) => {
      if (e.keyCode == keycodes.enter) {
        e.preventDefault()
        //todo 参考Navicat Premium, 后面需要调整
      }
    }



    const context = {
      submitActionRef: createRef(null)
    }

    onMounted(() => {
      window.addEventListener('keydown', handleEnter)
    })

    //todo 参考Navicat Premium
    provide(contextKey, context)

    onBeforeUnmount(() => {
      window.removeEventListener('keydown', handleEnter)
    })



    return () => slots.default?.()
  }
})
