import {
  computed,
  createVNode,
  defineComponent,
  onBeforeUnmount,
  onMounted,
  PropType,
  ref,
  toRefs,
  unref,
  watch
} from 'vue'
import {storeToRefs} from 'pinia'
import {filterName} from '/@/second/keeper-tools'
import {Modal} from "ant-design-vue";
import {ExclamationCircleOutlined} from "@ant-design/icons-vue";
import {getLocalStorage} from '/@/second/utility/storageCache'
import {useBootstrapStore} from "/@/store/modules/bootstrap"
import {get, uniq} from 'lodash-es'
import AppObjectCore from '/@/second/appobj/AppObjectCore.vue'
import getConnectionLabel from '/@/second/utility/getConnectionLabel'
import {ConnectionsWithStatus} from '/@/second/typings/mysql'
import {IPinnedDatabasesItem} from '/@/second/typings/types/standard.d'
import {connectionDeleteApi} from '/@/api/simpleApis'
import {serverConnectionsRefreshApi} from '/@/api/simpleApis'

export default defineComponent({
  name: 'ConnectionAppObject',
  props: {
    data: {
      type: Object as PropType<ConnectionsWithStatus>,
    },
    passProps: {
      type: Object as PropType<{ showPinnedInsteadOfUnpin: boolean }>,
      default: {showPinnedInsteadOfUnpin: true}
    },
    statusIcon: {
      type: String as PropType<string>
    },
    statusTitle: {
      type: String as PropType<string>
    },
    extInfo: {
      type: String as PropType<string>
    },
    engineStatusIcon: {
      type: String as PropType<string>
    },
    engineStatusTitle: {
      type: String as PropType<string>
    },
  },
  setup(props, {attrs}) {
    const {
      data,
      extInfo,
      engineStatusIcon,
      engineStatusTitle,
      statusIcon,
      statusTitle,
    } = toRefs(props)
    const statusTitleRef = ref()
    const statusIconRef = ref()
    const extInfoRef = ref()
    const engineStatusIconRef = ref()
    const engineStatusTitleRef = ref()
    const bootstrap = useBootstrapStore()
    const {extensions, openedConnections} = storeToRefs(bootstrap)
    let timerId: ReturnType<typeof setTimeout> | null

    const handleConnect = () => {
      if (unref(data)!.singleDatabase) {
        bootstrap.subscribeCurrentDatabase({
          connection: unref(data)!,
          name: unref(data)!.defaultDatabase
        } as unknown as IPinnedDatabasesItem)
      } else {
        bootstrap.subscribeOpenedConnections(uniq([...bootstrap.getOpenedConnections, unref(data)!._id]))
        timerId = setTimeout(() => {
          void serverConnectionsRefreshApi({
            conid: unref(data)!._id,
            keepOpen: true,
          })
        })
      }
    }

    watch(() => [data.value, extensions.value], () => {
      if (extensions.value?.drivers.find(x => x.engine == data.value?.engine)) {
        const match = (unref(data)!.engine || '').match(/^([^@]*)@/)
        extInfoRef.value = match ? match[1] : unref(data)!.engine;
        engineStatusIconRef.value = null
        engineStatusTitleRef.value = null
      } else {
        extInfo.value = data.value?.engine
        engineStatusIconRef.value = 'img warn'
        engineStatusTitleRef.value = `Engine driver ${data.value?.engine} not found, review installed plugins and change engine in edit connection dialog`
      }
    }, {
      immediate: true,
    })

    watch(() => [data.value, openedConnections.value], () => {
      const {_id, status} = unref(data)!
      if (openedConnections.value.includes(_id)) {
        if (!status) statusIconRef.value = 'icon loading'
        else if (status.name == 'pending') statusIconRef.value = 'icon loading';
        else if (status.name == 'ok') statusIconRef.value = 'img ok';
        else statusIconRef.value = 'img error';
        if (status && status.name == 'error') {
          statusTitleRef.value = status.message
        }
      } else {
        statusIconRef.value = null
        statusTitleRef.value = null
      }
    }, {
      immediate: true
    })

    onMounted(() => {
      // dataBase.subscribeExtensions(buildExtensions() as any)
      statusTitleRef.value = unref(statusTitle)
      statusIconRef.value = unref(statusIcon)
      extInfoRef.value = unref(extInfo)
      engineStatusIconRef.value = unref(engineStatusIcon)
      engineStatusTitleRef.value = unref(engineStatusTitle)
    })

    onBeforeUnmount(() => {
      timerId && clearTimeout(timerId)
    })

    const currentDatabase = computed(() => bootstrap.$state.currentDatabase)

    const handleDelete = async () => {
      const r = Modal.confirm({
        title: 'Confirm',
        icon: createVNode(ExclamationCircleOutlined),
        content: `Really delete connection ${getConnectionLabel(data.value)}${data.value?.port ? '_' + data.value?.port : ''}?`,
        okText: '确认',
        cancelText: '取消',
        onOk: async () => {
          try {
            await connectionDeleteApi({_id: data.value?._id})
            r.destroy()
          } catch (e) {
            console.log(e)
          }
        },
        onCancel: () => r.destroy(),
      })
    }

    // const addWailsEventListener = () => {
    //   EventsOn("connection-list-changed", data => {
    //     console.log(data, 'connections/list');
    //     console.log(data, 'connections/get');
    //   })
    // }

    const handleClick = async () => {
      // handleConnect()
    }

    const getContextMenu = () => {
      return [
        {
          label: 'Edit',
          handler: () => {
            console.log('click delete')
          },
        },
        {
          label: 'Delete',
          handler: handleDelete,
        },
        {
          label: 'Duplicate',
          handler: () => {
            console.log('click open');
          },
        }
      ]
    }

    return () => {
      const {...restProps} = attrs
      return <AppObjectCore
        {...restProps}
        data={unref(data) as ConnectionsWithStatus}
        title={getConnectionLabel(unref(data))}
        icon={unref(data)!.singleDatabase ? 'img database' : 'img server'}
        isBold={unref(data)!.singleDatabase
          ? get(unref(currentDatabase), 'connection._id') == unref(data)!._id && get(unref(currentDatabase), 'name') == unref(data)!.defaultDatabase
          : get(unref(currentDatabase), 'connection._id') == unref(data)!._id}

        statusIcon={statusIconRef.value || engineStatusIconRef.value}
        statusTitle={statusTitleRef.value || engineStatusTitleRef.value}
        // statusIconBefore={data!.isReadOnly ? 'icon lock' : null}
        extInfo={unref(extInfoRef)}
        menu={getContextMenu}
        onClick={handleClick}
        onDblclick={handleConnect}
      />
    }
  },
  extractKey: data => data._id,
  createMatcher: props => filter => {
    const {_id, displayName, server} = props;
    const databases = getLocalStorage(`database_list_${_id}`) || [];
    return filterName(unref(filter), displayName, server, ...databases.map(x => x.name))
  },
  createChildMatcher: props => filter => {
    if (!filter) return false;
    const {_id} = props;
    const databases = getLocalStorage(`database_list_${_id}`) || [];
    return filterName(unref(filter), ...databases.map(x => x.name));
  }
})
