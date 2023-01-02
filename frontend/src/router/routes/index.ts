import type { AppRouteRecordRaw } from '/@/router/types';
import {LAYOUT} from '../constant'

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
