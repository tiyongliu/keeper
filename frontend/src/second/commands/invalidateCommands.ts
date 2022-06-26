import {nextTick} from 'vue'
// import {dataBaseStore} from '/@/store/modules/dataBase'

let isInvalidated = false
// const dataBase = dataBaseStore()
export default async function invalidateCommands() {
  if (isInvalidated) return
  isInvalidated = true
  await nextTick()

  isInvalidated = false

  // dataBase.$state.commands

}
