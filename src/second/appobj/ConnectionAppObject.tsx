import {computed, defineComponent, onMounted, PropType, unref, watch, ref} from 'vue'
import {getLocalStorage} from '/@/second/utility/storageCache'
import {filterName} from "/@/packages/tools/src/filterName";
import { dataBaseStore } from "/@/store/modules/dataBase";
import {uniq, get} from 'lodash-es'

import AppObjectCore from '/@/second/appobj/AppObjectCore.vue'
import getConnectionLabel from '/@/second/utility/getConnectionLabel'

import {IConnectionAppObjectData} from '/@/second/types/IStore.d'
export default defineComponent({
  props: {
    data: {
      type: Object as PropType<IConnectionAppObjectData>,
      //todo delete dynamic default
      default: {
        "server":"localhost",
        "engine":"mysql@dbgate-plugin-mysql",
        "sshMode":"userPassword",
        "sshPort":"22",
        "sshKeyfile":"C:\\Users\\Administrator\\.ssh\\id_rsa",
        "port":"5001",
        "user":"root",
        "password":"crypt:6f4a500c408bce8261606389954de288bdb5d66c02d0e8b1ffe4275cba54d24a5942f62d29ae1703d1f9de29e468af72adlfTFoqebJEvRNqKoigxQ==",
        "_id":"b9c51b10-b354-11ec-812a-3d58c681a37b",
        "status":{
          "name":"ok"
        },
      }
    },
    passProps: {
      type: Object as PropType<{ showPinnedInsteadOfUnpin: boolean }>,
      default: {
        showPinnedInsteadOfUnpin: true
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
    }
  },
  components: {
    AppObjectCore
  },
  setup(props) {
    const {data, extInfo, engineStatusIcon, engineStatusTitle, statusIcon, statusTitle} = props

    let statusTitleRef = ref()
    let statusIconRef = ref()
    let extInfoRef = ref()
    let engineStatusIconRef = ref()
    let engineStatusTitleRef = ref()
    const dataBase = dataBaseStore()

    const handleConnect = () => {
      if (unref(data!).singleDatabase) {
        dataBase.subscribeCurrentDatabase({connection: data, name: data.defaultDatabase})
      } else {
        dataBase.subscribeOpenedConnections(uniq([... dataBase.$state.openedConnections, data._id]))
      }
    }

    const getContextMenu = () => {

    }

    watch(() => unref(dataBase.$state.extensions), () => {
      if (unref(dataBase.$state.extensions!).drivers.find(x => x.engine == data.engine)) {
        const match = (data.engine || '').match(/^([^@]*)@/)
        extInfoRef.value = match ? match[1] : data.engine;
        engineStatusIconRef.value = null
        engineStatusTitleRef.value = null
      } else {
        extInfoRef.value = data.engine;
        engineStatusIconRef.value = 'img warn'
        engineStatusTitleRef.value = `Engine driver ${data.engine} not found, review installed plugins and change engine in edit connection dialog`
      }
    })

    watch(() => unref(data), () => {
      const {_id, status} = unref(data)
      if (dataBase.$state.openedConnections.includes(_id)) {
        if (!status) statusIconRef.value = 'icon loading'
        else if (status.name == 'pending') statusIconRef.value = 'icon loading';
        else if (status.name == 'ok') statusIconRef.value = 'img ok'
        else statusIconRef.value = 'img error';
        if (status && status.name == 'error') {
          statusTitleRef.value = status.message
        } else {
          statusIconRef.value = null
          statusTitleRef.value = null
        }
      }
    })


    onMounted(() => {
      statusTitleRef.value = unref(statusTitle)
      statusIconRef.value = unref(statusIcon)
      extInfoRef.value = unref(extInfo)
      engineStatusIconRef.value = unref(engineStatusIcon)
      engineStatusTitleRef.value = unref(engineStatusTitle)
    })

    return () => {
      return <AppObjectCore
        data={data}
        title={getConnectionLabel(data)}
        icon={data.singleDatabase ? 'img database' : 'img server'}
        isBold={data.singleDatabase
          ? get(dataBase.$state.currentDatabase, 'connection._id') == data._id && get(dataBase.$state.currentDatabase, 'name') == data.defaultDatabase
          : get(dataBase.$state.currentDatabase, 'connection._id') == data._id}
        // statusIcon={unref(statusIconRef) || unref(engineStatusIconRef)}
        statusIcon={`img ok`}
        statusTitle={unref(statusTitleRef) || unref(engineStatusTitleRef)}
        statusIconBefore={data.isReadOnly ? 'icon lock' : null}
        // extInfo={unref(extInfoRef)}
        extInfo={`mysql`}
        menu={getContextMenu}
        onClick={handleConnect}
      />
    }
  }
})

export const extractKey = data => data._id;
export const createMatcher = props => filter => {
  const { _id, displayName, server } = props;
  const databases = getLocalStorage(`database_list_${_id}`) || [];
  return filterName(filter, displayName, server, ...databases.map(x => x.name));
};
export const createChildMatcher = props => filter => {
  if (!filter) return false;
  const { _id } = props;
  const databases = getLocalStorage(`database_list_${_id}`) || [];
  return filterName(filter, ...databases.map(x => x.name));
};
