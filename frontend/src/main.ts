//Windi CSS 是下一代工具优先的 CSS 框架
import 'virtual:windi-base.css';
import 'virtual:windi-components.css';

//utilities 引入这个css库，默认有媒体查询@media样式，需要注释
// import 'virtual:windi-utilities.css';
//vben admin 添加的样式，我们不需要，所以注释掉。
// import '/@/design/index.less';
import '/@/design/keeper.less'
// Register icon sprite
import 'virtual:svg-icons-register';
import App from './App.vue';
import {createApp} from 'vue';
import {initAppConfigStore} from '/@/logics/initAppConfig';
import {setupRouter} from '/@/router';
import {setupStore} from '/@/store';
import {setupGlobDirectives} from '/@/directives';
import {setupI18n} from '/@/locales/setupI18n';
import {registerGlobComp} from '/@/components/registerGlobComp';

async function bootstrap() {
  const app = createApp(App);

  // Configure store
  // 配置 store
  setupStore(app);

  // Initialize internal system configuration
  // 初始化内部系统配置
  initAppConfigStore();

  // Register global components
  // 注册全局组件
  registerGlobComp(app);

  // Multilingual configuration
  // 多语言配置
  // Asynchronous case: language files may be obtained from the server side
  // 异步案例：语言文件可能从服务器端获取
  await setupI18n(app);

  // Configure routing
  // 配置路由
  setupRouter(app);

  // Register global directive
  // 注册全局指令
  setupGlobDirectives(app);

  // https://next.router.vuejs.org/api/#isready
  // await router.isReady();

  app.mount('#app');
}

bootstrap();
