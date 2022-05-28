export function getExpandIcon(expandable, subItemsComponent, isExpanded, expandIconFunc) {
  if (!subItemsComponent) return null;
  if (!expandable) return 'icon invisible-box';
  return expandIconFunc(isExpanded);
}
