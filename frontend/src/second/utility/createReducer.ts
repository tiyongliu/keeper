import {ref} from 'vue'

export default function createReducer(reducer, initialState): any {
  const state = ref(initialState)

  function dispatch(action) {
    state.value = reducer(state.value, action)
  }

  return [state, dispatch]
}
