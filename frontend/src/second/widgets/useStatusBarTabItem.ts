import {inject, onBeforeUnmount, Ref, unref, watch} from 'vue'
import {buildUUID} from '/@/utils/uuid'
import {updateStatusBarInfoItem} from '/@/second/utility/statusBarStore'

export function useStatusBarTabItem(text?: Ref<number | undefined>, icon?: Ref<string | undefined>) {
  const key = buildUUID()
  const tabid = inject('tabid')

  watch(() => [unref(tabid), text?.value, icon?.value], () => {
    //hooks 需要查看props怎么传递值
    if (text?.value && icon?.value) {
      //todo updateStatuBarInfoItem(tabid, key, { text, icon, clickable, onClick })
      updateStatusBarInfoItem(unref(tabid), key, {text, icon: icon?.value})
    }
  }, {immediate: true})

  onBeforeUnmount(() => updateStatusBarInfoItem(unref(tabid), key, null))
}
