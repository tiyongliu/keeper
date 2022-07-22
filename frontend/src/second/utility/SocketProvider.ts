import {onMounted, ref} from 'vue'
import {EventsOff, EventsOn} from '/@/wailsjs/runtime/runtime'
import {cacheClean} from "/@/second/utility/cache";

export function useSocket() {
  const socket = ref<{ on: unknown, off: unknown } | null>(null)
  onMounted(() => {
    socket.value = newSocket
    newSocket.on("clean-cache", reloadTrigger => cacheClean(reloadTrigger))
  })

  return socket
}


const newSocket = {
  on: EventsOn,
  off: EventsOff,
}
