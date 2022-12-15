import {isEmpty, keys, pick} from 'lodash-es'
import {useInstalledPlugins} from '/@/api/bridge'
import {useBootstrapStore} from "/@/store/modules/bootstrap"
import {onBeforeUnmount, onMounted, ref, watch} from 'vue'
import {ExtensionsDirectory} from "/@/second/typings/types/extensions";
import mysql from './keeper-plugin-mysql'
import mongo from './keeper-plugin-mongo'

const plugins = {
  mysql,
  mongo
}

export default function initPluginsProvider() {
  const installedPlugins = ref()
  const bootstrap = useBootstrapStore()
  let pluginsDict = {}

  onMounted(() => {
    useInstalledPlugins({}, installedPlugins)
  })

  onBeforeUnmount(() => {
    installedPlugins.value = null
    pluginsDict = {}
  })


  watch(() => installedPlugins.value, () => {
    loadPlugins(pluginsDict, installedPlugins.value, bootstrap)
      .then(newPlugins => {
        if (isEmpty(newPlugins)) return
        pluginsDict = pick(
          {...pluginsDict, ...(newPlugins as object)},
          installedPlugins.value.map(y => y.name)
        )
      })
      .then(() => {
        bootstrap.setExtensions(
          buildExtensions(buildPlugins(installedPlugins.value))
        )
      })
  })

  function buildPlugins(installedPlugins) {
    return (installedPlugins || [])
      .map(manifest => ({
        packageName: manifest.name,
        manifest,
        content: pluginsDict[manifest.name],
      }))
      .filter(x => x.content)
  }
}

async function loadPlugins(pluginsDict, installedPlugins, bootstrap) {
  const newPlugins = {}
  for (const installed of installedPlugins || []) {
    if (!keys(pluginsDict).includes(installed.name)) {
      bootstrap.setLoadingPluginStore({
        loaded: false,
        loadingPackageName: installed.name
      })

      //这种方法开发环境可以，编译后无法访问。
      // const defaultFrontend = await import(`./keeper-plugin-${installed.name}`)
      const defaultFrontend = await plugins[installed.name]
      newPlugins[installed.name] = defaultFrontend ?? {}
    }
  }

  if (installedPlugins) {
    bootstrap.setLoadingPluginStore({
      loaded: true,
      loadingPackageName: null
    })
  }

  return newPlugins
}

function buildDrivers(plugins) {
  const res = [];
  for (const {content} of plugins) {
    if (content.drivers) { // @ts-ignore
      res.push(...content.drivers);
    }
  }
  return res;
}

function buildExtensions(plugins): ExtensionsDirectory {
  return {
    plugins,
    drivers: buildDrivers(plugins),
  }
}
