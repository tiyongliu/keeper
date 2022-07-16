import {
  computed,
  createVNode,
  defineComponent,
  onMounted,
  PropType,
  ref,
  toRefs,
  unref,
  watch
} from 'vue'
import {filterName} from 'keeper-tools'
import {Modal} from "ant-design-vue";
import {ExclamationCircleOutlined} from "@ant-design/icons-vue";
import {getLocalStorage} from '/@/second/utility/storageCache'
import {dataBaseStore} from "/@/store/modules/dataBase"
import {get, uniq} from 'lodash-es'
import AppObjectCore from '/@/second/appobj/AppObjectCore.vue'
import getConnectionLabel from '/@/second/utility/getConnectionLabel'
import {ConnectionsWithStatus} from '/@/second/typings/mysql'
import {IPinnedDatabasesItem} from '/@/second/typings/types/standard.d'
import {handleDeleteApi} from '/@/api/connection'
import {handleRefreshApi} from '/@/api/serverConnections'
import {
  connectionListChangedEvent,
  databaseListChangedEvent,
  databaseStructureChangedEvent,
  serverStatusChangedEvent,
} from "/@/api/event"

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
      statusTitle
    } = toRefs(props)
    const statusTitleRef = ref()
    const statusIconRef = ref()
    const extInfoRef = ref()
    const engineStatusIconRef = ref()
    const engineStatusTitleRef = ref()
    const dataBase = dataBaseStore()

    const handleConnect = () => {
      if (unref(data)!.singleDatabase) {
        dataBase.subscribeCurrentDatabase({
          connection: unref(data)!,
          name: unref(data)!.defaultDatabase
        } as unknown as IPinnedDatabasesItem)
      } else {


        dataBase.subscribeOpenedConnections(uniq([...dataBase.$state.openedConnections, unref(data)!._id]))
        console.log("fsjfdsdafjsfsdf", unref(data))
      }
    }

    watch(() => unref(dataBase.$state.extensions), () => watchExtensions())

    const watchExtensions = () => {
      const match = (unref(data)!.engine || '').match(/^([^@]*)@/)
      extInfoRef.value = match ? match[1] : unref(data)!.engine;
      engineStatusIconRef.value = null
      engineStatusTitleRef.value = null

      // if (unref(dataBase.$state.extensions!).drivers.find(x => x.engine == data!.engine)) {
      //   const match = (data!.engine || '').match(/^([^@]*)@/)
      //   extInfoRef.value = match ? match[1] : data!.engine;
      //   engineStatusIconRef.value = null
      //   engineStatusTitleRef.value = nulld
      // } else {
      //
      //   extInfoRef.value = data!.engine;
      //   engineStatusIconRef.value = 'img warn'
      //   engineStatusTitleRef.value = `Engine schema ${data!.engine} not found, review installed plugins and change engine in edit connection dialog`
      // }
    }

    watch(() => unref(data), () => watchStatus())

    const watchStatus = () => {
      const {_id, status} = unref(data)!
      if (dataBase.$state.openedConnections.includes(_id)) {
        if (!status) statusIconRef.value = 'icon loading'
        else if (status.name == 'pending') statusIconRef.value = 'icon loading';
        else if (status.name == 'ok')  statusIconRef.value = 'img ok';
        else statusIconRef.value = 'img error';
        if (status && status.name == 'error') {
          statusTitleRef.value = status.message
        }
      } else {
        statusIconRef.value = null
        statusTitleRef.value = null
      }
    }

    onMounted(() => {
      // dataBase.subscribeExtensions(buildExtensions() as any)
      statusTitleRef.value = unref(statusTitle)
      statusIconRef.value = unref(statusIcon)
      extInfoRef.value = unref(extInfo)
      engineStatusIconRef.value = unref(engineStatusIcon)
      engineStatusTitleRef.value = unref(engineStatusTitle)
      watchExtensions()
      watchStatus()

      if (window.runtime) {
        connectionListChangedEvent()
        serverStatusChangedEvent()
        databaseListChangedEvent()
        databaseStructureChangedEvent()
      }
    })

    const currentDatabase = computed(() => dataBase.$state.currentDatabase)

    const handleDelete = async () => {
      const r = Modal.confirm({
        title: 'Confirm',
        icon: createVNode(ExclamationCircleOutlined),
        content: `Really delete connection ${getConnectionLabel(data.value)}${data.value?.port ? '_' + data.value?.port : ''}?`,
        okText: '确认',
        cancelText: '取消',
        onOk: async () => {
          try {
            await handleDeleteApi({_id: data.value?._id})
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
      dataBase.subscribeCurrentDatabase({connection: data.value})
      await handleRefreshApi({conid: data.value!._id, keepOpen: true})
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
      const {onClick, onExpand, ...restProps} = attrs
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
  }
})

export const extractKey = data => data._id;
export const createMatcher = props => filter => {
  const {_id, displayName, server} = props;
  const databases = getLocalStorage(`database_list_${_id}`) || [];
  return filterName(unref(filter), displayName, server, ...databases.map(x => x.name));
};
export const createChildMatcher = props => filter => {

  if (!filter) {
    return false;
  }
  const {_id} = props;
  const databases = getLocalStorage(`database_list_${_id}`) || [];
  return filterName(unref(filter), ...databases.map(x => x.name));
};

function openConnection(connection) {
  if (connection.singleDatabase) {
    //currentDatabase.set({ connection, name: connection.defaultDatabase })

  }
}
