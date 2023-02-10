import {createVNode, defineComponent, onMounted, PropType, ref, toRefs, unref, watch} from 'vue'
import {storeToRefs} from 'pinia'
import {filterName} from '/@/second/keeper-tools'
import {Modal, message} from "ant-design-vue";
import {ExclamationCircleOutlined} from "@ant-design/icons-vue";
import {getLocalStorage} from '/@/second/utility/storageCache'
import {useBootstrapStore} from "/@/store/modules/bootstrap"
import {get, uniq} from 'lodash-es'
import AppObjectCore from '/@/second/appobj/AppObjectCore.vue'
import getConnectionLabel from '/@/second/utility/getConnectionLabel'
import {ConnectionsWithStatus} from '/@/second/typings/mysql'
import {IPinnedDatabasesItem} from '/@/second/typings/types/standard.d'
import {
  connectionDeleteApi,
  databaseConnectionsRefreshApi,
  serverConnectionsRefreshApi
} from '/@/api/simpleApis'
import openNewTab from '/@/second/utility/openNewTab'

export default defineComponent({
  name: 'ConnectionAppObject',
  props: {
    data: {
      type: Object as PropType<ConnectionsWithStatus>,
    },
    passProps: {
      type: Object as PropType<{ showPinnedInsteadOfUnpin: boolean }>,
      default: () => {
        return {showPinnedInsteadOfUnpin: true}
      }
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
    const {extensions, openedConnections, currentDatabase} = storeToRefs(bootstrap)
    // let timerId: ReturnType<typeof setTimeout> | null

    const handleConnect = () => {
      openConnection(data.value, bootstrap)
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
      // dataBase.setExtensions(buildExtensions() as any)
      statusTitleRef.value = unref(statusTitle)
      statusIconRef.value = unref(statusIcon)
      extInfoRef.value = unref(extInfo)
      engineStatusIconRef.value = unref(engineStatusIcon)
      engineStatusTitleRef.value = unref(engineStatusTitle)
    })

    const handleDelete = async () => {
      const r = Modal.confirm({
        title: 'Confirm',
        icon: createVNode(ExclamationCircleOutlined),
        content: `Really delete connection ${getConnectionLabel(data.value)}${data.value?.port ? '_' + data.value?.port : ''}?`,
        okText: 'Ok',
        cancelText: 'Cancel',
        onOk: async () => {
          try {
            await connectionDeleteApi({_id: data.value?._id})
            // todo 暂时不使用删除连接池并判断是否为当前连接并清楚
            // await bootstrap.removeCurrentDatabase(data.value?._id)
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

    const handleSqlRestore = () => {

    }

    const getContextMenu = () => {
      const driver = extensions.value && extensions.value?.drivers.find(x => x.engine == data.value?.engine);
      const handleRefresh = () => {
        void serverConnectionsRefreshApi({conid: data.value?._id})
      }
      const handleDisconnect = () => {
        disconnectServerConnection(data.value?._id);
      }

      const handleCreateDatabase = () => {

      }

      const handleServerSummary = () => {
        void openNewTab({
          title: getConnectionLabel(data.value),
          icon: 'img server',
          tabComponent: 'ServerSummaryTab',
          props: {
            conid: data.value?._id,
          },
        });
      }
      const handleNewQuery = () => {
        const tooltip = `${getConnectionLabel(data.value)}`;
        openNewTab({
          title: 'Query #',
          icon: 'img sql-file',
          tooltip,
          tabComponent: 'QueryTab',
          props: {
            conid: data.value?._id,
          },
        });
      }

      return [
        [
          {
            label: data.value && bootstrap.getOpenedConnections.includes(data.value?._id) ? 'View details' : 'Edit',
            onClick: () => {
              message.warning('developing')
            },
          },
          // !(data.value && bootstrap.getOpenedConnections.includes(data.value?._id)) && {
          {
            label: 'Delete',
            onClick: handleDelete,
          },
          {
            label: 'Duplicate',
            onClick: () => {
              message.warning('developing')
            }
          },
          // {onClick: handleNewQuery, text: 'New query', isNewQuery: true},
          // {
          //   text: 'Connect',
          //   onClick: handleConnect,
          // },
          (data.value && bootstrap.getOpenedConnections.includes(data.value?._id) && data.value?.status) && {
            text: 'Refresh',
            onClick: handleRefresh,
          },
          data.value && bootstrap.getOpenedConnections.includes(data.value?._id) && {
            text: 'Disconnect',
            onClick: handleDisconnect,
          },
          {
            text: 'Create database',
            onClick: handleCreateDatabase,
          },
          {
            text: 'Server summary',
            onClick: handleServerSummary,
          }
        ],
        data.value?.singleDatabase && [
          {divider: true},
        ],
        (driver && driver?.databaseEngineTypes?.includes('sql')) && {
          onClick: handleSqlRestore,
          text: 'Restore/import SQL dump'
        }
      ]
    }

    return () => {
      const {...restProps} = attrs
      return <AppObjectCore
        {...restProps}
        data={data.value as ConnectionsWithStatus}
        title={getConnectionLabel(data.value)}
        icon={data.value!.singleDatabase ? 'img database' : 'img server'}
        isBold={data.value!.singleDatabase
          ? get(currentDatabase.value, 'connection._id') == data.value!._id && get(currentDatabase.value, 'name') == data.value!.defaultDatabase
          : get(currentDatabase.value, 'connection._id') == data.value!._id}
        statusIcon={statusIconRef.value || engineStatusIconRef.value}
        statusTitle={statusTitleRef.value || engineStatusTitleRef.value}
        statusIconBefore={data.value && data.value.isReadOnly ? 'icon lock' : null}
        extInfo={extInfoRef.value}
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


export function disconnectServerConnection(conid, showConfirmation = true) {

}

export function openConnection(connection, bootstrap) {

  if (connection!.singleDatabase) {
    bootstrap.setCurrentDatabase({
      connection: connection!,
      name: connection!.defaultDatabase
    } as unknown as IPinnedDatabasesItem)
    void databaseConnectionsRefreshApi({
      conid: connection._id!,
      database: connection.defaultDatabase!,
      keepOpen: true
    })
    bootstrap.updateOpenedSingleDatabaseConnections(x => uniq([...x, connection._id]))
  } else {
    bootstrap.updateOpenedConnections(x => uniq([...x, connection!._id]))
    void serverConnectionsRefreshApi({
      conid: connection!._id,
      keepOpen: true,
    })
    bootstrap.updateExpandedConnections(x => uniq([...x, connection._id]))
  }
}
