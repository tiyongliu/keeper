import { RoleEnum } from '/@/enums/roleEnum';
// User permissions related operations
export function usePermission() {
  /**
   * Determine whether there is permission
   */
  function hasPermission(value?: RoleEnum | RoleEnum[] | string | string[], def = true): boolean {
    // Visible by default
    if (!value) {
      return def;
    }
    return true;
  }
  return { hasPermission };
}
