import {inject, onBeforeUnmount, unref, watch} from 'vue'
import {buildUUID} from '/@/utils/uuid'
import {updateStatusBarInfoItem} from '/@/second/utility/statusBarStore'

export function useStatusBarTabItem() {
  const key = buildUUID()
  const tabid = inject('tabid')
  watch(() => [unref(tabid)], () => {
    //hooks 需要查看props怎么传递值
    //todo updateStatuBarInfoItem(tabid, key, { text, icon, clickable, onClick })
    updateStatusBarInfoItem(unref(tabid), key, {})
  }, {immediate: true})

  onBeforeUnmount(() => updateStatusBarInfoItem(unref(tabid), key, null))
}
