export interface ContextMenuItem {
  divider?: boolean;
  text?: string
  label: string;
  keyText?: string
  onClick?: Fn;
  submenu?: ContextMenuItem[]
}

export interface CreateContextOptions {
  event: MouseEvent;
  icon?: string;
  styles?: any;
  items?: ContextMenuItem[] | Fn;
}

export interface ContextMenuProps {
  event?: MouseEvent;
  styles?: any;
  items: ContextMenuItem[] | Fn;
  left?: number;
  top?: number;
  targetElement?: string | string[]
}

export interface ItemContentProps {
  showIcon: boolean | undefined;
  item: ContextMenuItem;
  onClick: Fn;
}
