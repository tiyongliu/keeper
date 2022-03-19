import {nextTick} from 'vue'

let isInvalidated = false
export default async function invalidateCommands() {
  if (isInvalidated) return
  isInvalidated = true
  await nextTick()

  isInvalidated = false
}
