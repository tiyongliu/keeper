import type { AppRouteRecordRaw } from '/@/router/types';
import {
  LAYOUT,
  PAGE_NOT_FOUND_NAME,
} from '/@/router/constant';

// 404 on a page
export const PAGE_NOT_FOUND_ROUTE: AppRouteRecordRaw = {
  path: '/:path(.*)*',
  name: PAGE_NOT_FOUND_NAME,
  component: LAYOUT,
  meta: {
    title: 'Dashboard',
    hideBreadcrumb: true,
    hideMenu: true,
  },
};
