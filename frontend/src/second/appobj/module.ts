import {unref} from 'vue'
export function getExpandIcon(_expandable, _subItemsComponent, _isExpanded, _expandIconFunc) {
  const expandable = unref(_expandable)
  const subItemsComponent = unref(_subItemsComponent)
  const isExpanded = unref(_isExpanded)
  const expandIconFunc = unref(_expandIconFunc)

  if (!subItemsComponent) return null
  if (!expandable) return 'icon invisible-box'
  return expandIconFunc(isExpanded);
}
