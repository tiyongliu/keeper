import { MenuModeEnum, MenuTypeEnum } from '/@/enums/menuEnum';

export interface BeforeMiniState {
  menuCollapsed?: boolean;
  menuSplit?: boolean;
  menuMode?: MenuModeEnum;
  menuType?: MenuTypeEnum;
}
