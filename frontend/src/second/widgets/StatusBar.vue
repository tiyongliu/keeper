<template>
  <div class="main" >
    <div class="container">
      <template v-if="databaseName">
        <div class="item">
          <FontIcon icon="icon lock"  padRight v-if="connection.isReadOnly"/>
          <FontIcon icon="icon database" padRight  v-else/>
          {{databaseName}}
        </div>
        <div v-if="dbid" class="item clickable" title="Database color. Overrides connection color">
          <div  class="colorbox" :style="{databaseButtonBackground}">
            <FontIcon icon="icon palette" />
          </div>
        </div>
      </template>
      <template v-if="connectionLabel">
        <div class="item">
          <FontIcon icon="icon server" padRight />
          {{connectionLabel}}
        </div>
        <div class="item clickable" title="Connection color. Can be overriden by database color">
          <div :style="connectionButtonBackground" class="colorbox">
            <FontIcon icon="icon palette" />
          </div>
        </div>
      </template>
      <div class="item" v-if="connection?.user">
        <FontIcon icon="icon account" padRight />
        {{connection.user}}
      </div>
      <div class="item clickable" >
        <template v-if="status.name == 'pending'">
          <FontIcon icon="icon loading" padRight /> Loading
        </template>
        <template v-else-if="status.name == 'checkStructure'">
          <FontIcon icon="icon loading" padRight /> Checking model
        </template>
        <template v-else-if="status.name == 'loadStructure'">
          <FontIcon icon="icon loading" padRight /> Loading model
        </template>
        <template v-else-if="status.name == 'ok'">
          <FontIcon icon="img ok-inv" padRight /> Connected
        </template>
        <template v-else="status.name == 'error'">
          <FontIcon icon="img error-inv" padRight /> Error
        </template>
      </div>
      <div class="item" v-if="!connection">
        <FontIcon icon="icon disconnected" padRight /> Not connected
      </div>
      <div class="item flex" :title="serverVersion.version" v-if="serverVersion">
        <FontIcon icon="icon version" padRight />
        <div class="version">
          {{serverVersion.versionText || serverVersion.version}}
        </div>
      </div>
      <div class="item flex clickable" v-if="status?.analysedTime" :title="`Last ${databaseName} model refresh: ${moment(status?.analysedTime).format('HH:mm:ss')}\nClick for refresh DB model`" @click="handleSyncModel">
        <FontIcon icon="icon history" padRight />
        <div class="version ml-1">
          {{moment(status?.analysedTime).fromNow() + (timerValue ? '' : '')}}
        </div>
      </div>
    </div>
    <div class="container" v-for="item in contextItems">
      <div class="item" >
        {#if item.icon}
        <FontIcon icon={item.icon} padRight />
        {/if}
        {item.text}
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import moment from 'moment';
import { defineComponent, computed} from 'vue';
import FontIcon from '/@/second/icons/FontIcon.vue'
import {dataBaseStore} from "/@/store/modules/dataBase"
import { useConnectionColor } from "/@/second/utility/useConnectionColor"
import getConnectionLabel from "/@/second/utility/getConnectionLabel";
import { useDatabaseServerVersion,useDatabaseStatus } from "../utility/metadataLoaders"

export default defineComponent({
  name: 'StatusBar',
  components: {
    FontIcon
  },
  setup(props,attrs) {
    const dataBase = dataBaseStore()

    useConnectionColor()
    // const connectionBackground = computed(() => {
    //   return useConnectionColor()
    // });
    //


    const databaseName = computed(()=> dataBase.$state.currentDatabase?.name)
    const connection = computed(()=> dataBase.$state.currentDatabase?.connection)
    const dbid = computed(()=> connection ? { conid: connection._id, database: databaseName } : null)
    const connectionLabel = computed(()=> getConnectionLabel(connection,{allowExplicitDatabase: false}))
    const serverVersion = computed(()=> useDatabaseServerVersion(dbid || {}))
    const status = computed(()=> useDatabaseStatus(dbid || {}))
    const contextItems = []
    // const databaseButtonBackground = useConnectionColor(dbid, 6, 'dark', true, false)
    const databaseButtonBackground = '------'
    const connectionButtonBackground = '------'

    async function handleSyncModel(){

    }
    console.log( connection,'----connection')
    return {
      databaseName,
      connection,
      dbid,
      databaseButtonBackground,
      connectionButtonBackground,
      connectionLabel,
      status,
      serverVersion,
      handleSyncModel,
      contextItems,
      moment
    }
  }


})

</script>
<style lang="less" scoped>
.main {
  display: flex;
  color: var(--theme-font-inv-15);
  align-items: stretch;
  justify-content: space-between;
  cursor: default;
  flex: 1;
  background-color: blue;
}
.container {
  display: flex;
  align-items: stretch;
}
.item {
  padding: 0px 10px;
  display: flex;
  align-items: center;
  background-color: blue;
}

.version {
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.clickable {
  cursor: pointer;
}
.clickable:hover {
  background-color: var(--theme-bg-statusbar-inv-hover);
}

.colorbox {
  padding: 0px 3px;
  border-radius: 2px;
  color: var(--theme-bg-statusbar-inv-font);
  background: var(--theme-bg-statusbar-inv-bg);
}
</style>
