import {ref, watchEffect} from 'vue'
import {createGridConfig} from '/@/second/keeper-datalib'

function doLoadGridConfigFunc(tabid: string) {
  try {
    const existing = localStorage.getItem(`tabdata_grid_${tabid}`);
    if (existing) {
      return {
        ...createGridConfig(),
        ...JSON.parse(existing),
      };
    }
  } catch (err) {
    // @ts-ignore
    console.warn('Error loading grid config:', err.message);
  }
  return createGridConfig();
}

export default function useGridConfig(tabid: string) {
  const config = ref(doLoadGridConfigFunc(tabid))

  watchEffect(() => {
    // const value = unref(config)
    // localStorage.setItem(`tabdata_grid_${tabid}`, JSON.stringify(value))
  })

  // watch()

  return config
}
