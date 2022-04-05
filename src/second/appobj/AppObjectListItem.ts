export function getExpandIcon(expandable, subItemsComponent, isExpanded, expandIconFunc) {
  if (!subItemsComponent) return null;
  if (!expandable) return 'icon invisible-box';
  return expandIconFunc(isExpanded);
}

import {computed, defineComponent, onMounted, PropType, unref} from 'vue'

export default defineComponent({
  props: {
    isHidden: {
      type: Boolean as PropType<boolean>,
    }
  },
  setup() {

  }
})
