/**
 * Plugin to minimize and use ejs template syntax in index.html.
 * https://github.com/anncwb/vite-plugin-html
 */
import type { PluginOption } from 'vite';
import { createHtmlPlugin } from 'vite-plugin-html';

export function configHtmlPlugin(env: ViteEnv, isBuild: boolean) {
  const { VITE_GLOB_APP_TITLE } = env;

  const htmlPlugin: PluginOption[] = createHtmlPlugin({
    minify: isBuild,
    inject: {
      // Inject data into ejs template
      data: {
        title: VITE_GLOB_APP_TITLE,
      },

      // Embed the generated app.config.js file
      tags: isBuild
        ? [
            {
              tag: 'script',
              attrs: {
                // src: getAppConfigSrc(),
              },
              // 解决wails编译后本地src 资源加载配置文件404，外联改为内联
              children: `window.__PRODUCTION__KEEPER__CONF__={"VITE_GLOB_APP_TITLE":"keeper","VITE_GLOB_APP_SHORT_NAME":"keeper","VITE_GLOB_API_URL":"/basic-api","VITE_GLOB_UPLOAD_URL":"/upload","VITE_GLOB_API_URL_PREFIX":""};Object.freeze(window.__PRODUCTION__KEEPER__CONF__);Object.defineProperty(window,"__PRODUCTION__KEEPER__CONF__",{configurable:false,writable:false,});`,
            },
          ]
        : [],
    },
  });
  return htmlPlugin;
}
