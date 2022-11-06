import {ref} from 'vue'
export const statusBarTabInfo = ref({})
export function updateStatusBarInfoItem(tabid, key, item) {
  const items = statusBarTabInfo.value[tabid] || []
  let newItems;
  if (item == null) {
    newItems = items.filter(x => x.key != key);
  } else if (items.find(x => x.key == key)) {
    newItems = items.map(x => (x.key == key ? { ...item, key } : x));
  } else {
    newItems = [...items, { ...item, key }];
  }
  statusBarTabInfo.value = {
    ...statusBarTabInfo.value,
    [tabid]: newItems,
  }
}
