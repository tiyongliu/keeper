<template>
  <div class="main" >
    <template v-if="databaseName">
      <div class="container">
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
  </div>
</template>
<script lang="ts">
import moment from 'moment';
import { defineComponent, computed, computed, unref, toRefs, ref,reactive} from 'vue';
import FontIcon from '/@/second/icons/FontIcon.vue'
import {dataBaseStore} from "/@/store/modules/dataBase"
import { useConnectionColor } from "/@/second/utility/useConnectionColor"
import getConnectionLabel from "/@/second/utility/getConnectionLabel";


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
      const connectionLabel = getConnectionLabel(connection,{allowExplicitDatabase: false})

      // const databaseButtonBackground = useConnectionColor(dbid, 6, 'dark', true, false)
      const databaseButtonBackground = '------'
      const connectionButtonBackground = '------'

      return {
        databaseName,
        connection,
        dbid,
        databaseButtonBackground,
        connectionButtonBackground,
        connectionLabel
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
}
.container {
  display: flex;
  align-items: stretch;
}
.item {
  padding: 0px 10px;
  display: flex;
  align-items: center;
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
