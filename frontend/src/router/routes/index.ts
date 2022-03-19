import type { AppRouteRecordRaw } from '/@/router/types';
export const LAYOUT = () => import('/@/layouts/default/index.vue');
// import { PAGE_NOT_FOUND_ROUTE } from '/@/router/routes/basic';
// import { PageEnum } from '/@/enums/pageEnum';

// 根路由
export const RootRoute: AppRouteRecordRaw = {
  path: '/',
  name: 'Root',
  component: LAYOUT,
  meta: {
    title: 'Root',
  },
};

// Basic routing without permission
// 未经许可的基本路由
export const basicRoutes = [
  RootRoute,
];
