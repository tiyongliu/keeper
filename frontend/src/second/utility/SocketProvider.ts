import {onMounted, ref} from 'vue'
import {EventsOff, EventsOn} from '/@/wailsjs/runtime/runtime'
import {cacheClean} from "/@/second/utility/cache";

interface RuntimeEvent {
  on: (eventName: string, callback: (...data: any) => void) => void
  off: (eventName: string) => void
}

const newSocket: RuntimeEvent = {
  on: EventsOn,
  off: EventsOff,
}

export function useSocket() {
  const socket = ref<RuntimeEvent | null>(null)
  onMounted(() => {
    socket.value = newSocket
    newSocket.on("clean-cache", reloadTrigger => cacheClean(reloadTrigger))
  })
  console.log(`rerwrwerwwrrwewer`)

  return socket
}

