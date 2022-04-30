import {computed, defineComponent, onMounted, PropType, unref, watch, ref, toRefs} from 'vue'
import {getLocalStorage} from '/@/second/utility/storageCache'
import {filterName} from "/@/packages/tools/src/filterName";
import { dataBaseStore } from "/@/store/modules/dataBase"
import {uniq, get} from 'lodash-es'
import AppObjectCore from '/@/second/appobj/AppObjectCore.vue'
import getConnectionLabel from '/@/second/utility/getConnectionLabel'

import {IConnectionAppObjectData, IPinnedDatabasesItem} from '/@/second/types/standard.d'
export default defineComponent({
  name: 'ConnectionAppObject',
  props: {
    data: {
      type: Object as PropType<IConnectionAppObjectData>,
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
    },
  },
  setup(props, {attrs}) {
    const {data, extInfo, engineStatusIcon, engineStatusTitle, statusIcon, statusTitle} = toRefs(props)
    let statusTitleRef = ref()
    let statusIconRef = ref()
    let extInfoRef = ref()
    let engineStatusIconRef = ref()
    let engineStatusTitleRef = ref()
    const dataBase = dataBaseStore()

    const handleConnect = () => {
      if (unref(data)!.singleDatabase) {
        dataBase.subscribeCurrentDatabase({connection: unref(data)!, name: unref(data)!.defaultDatabase} as unknown as IPinnedDatabasesItem)
      } else {
        dataBase.subscribeOpenedConnections(uniq([... dataBase.$state.openedConnections, unref(data)!._id]))
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
      //   engineStatusTitleRef.value = null
      // } else {
      //
      //   extInfoRef.value = data!.engine;
      //   engineStatusIconRef.value = 'img warn'
      //   engineStatusTitleRef.value = `Engine driver ${data!.engine} not found, review installed plugins and change engine in edit connection dialog`
      // }
    }

    watch(() => unref(data), () => watchStatus())

    const watchStatus = () => {
      const {_id, status} = unref(data)!
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

    })

    const currentDatabase = computed(() => dataBase.$state.currentDatabase)

    return () => {
      const {onClick, onExpand, ...restProps} = attrs
       return <AppObjectCore
         {...restProps}
        data={unref(data) as IPinnedDatabasesItem}
        title={getConnectionLabel(unref(data))}
        icon={unref(data)!.singleDatabase ? 'img database' : 'img server'}
        isBold={unref(data)!.singleDatabase
          ? get(unref(currentDatabase), 'connection._id') == unref(data)!._id && get(unref(currentDatabase), 'name') == unref(data)!.defaultDatabase
          : get(unref(currentDatabase), 'connection._id') == unref(data)!._id}

        // statusIcon={unref(statusIconRef) || unref(engineStatusIconRef)}
        statusIcon={`img ok`}

        statusTitle={unref(statusTitleRef) || unref(engineStatusTitleRef)}
        // statusIconBefore={data!.isReadOnly ? 'icon lock' : null}
        extInfo={extInfoRef}
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
  if (!filter) {
    return false;
  }
  const { _id } = props;
  const databases = getLocalStorage(`database_list_${_id}`) || [];
  return filterName(filter, ...databases.map(x => x.name));
};
