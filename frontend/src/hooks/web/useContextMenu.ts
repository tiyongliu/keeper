import { onUnmounted, getCurrentInstance } from 'vue';
// import { createContextMenu, destroyContextMenu } from '/@/components/ContextMenu';
import { createContextMenu, destroyContextMenu } from '/@/second/modals/createContextMenu';
import type { ContextMenuItem } from '/@/components/ContextMenu';
export type { ContextMenuItem };
export function useContextMenu(authRemove = true) {
  if (getCurrentInstance() && authRemove) {
    onUnmounted(() => {
      destroyContextMenu();
    });
  }
  return [createContextMenu, destroyContextMenu];
}
