<template>
  <div class="panel_main">
    <div class="wrapper" :class="{'selected':item.name==test}" v-for="(item, index) in widgets"
         :key="index">
      <FontIcon :icon="item.icon" :title="item.title" @click="con(item)"/>
    </div>

    <div class="flex1">&nbsp;</div>

    <div class="wrapper">
      <FontIcon icon="icon settings"/>
    </div>
  </div>
</template>

<script>
import {reactive, toRefs} from 'vue';

import FontIcon from '../icons/FontIcon.vue';
import {getWithStorageVariableCache, setWithStorageVariableCache} from '/@/second/utility/storage'
// import mitt from '../../utils/mitt'
export default {
  setup(props, {emit}) {
    // const test = 'database'
    const state = reactive({
      widgets: [
        {
          icon: 'icon database',
          name: 'database',
          title: 'Database connections',
        },
        {
          icon: 'icon file',
          name: 'file',
          title: 'Favorites & Saved files',
        },
        {
          icon: 'icon history',
          name: 'history',
          title: 'Query history & Closed tabs',
        },
        {
          icon: 'icon archive',
          name: 'archive',
          title: 'Archive (saved tabular data)',
        },
        {
          icon: 'icon plugin',
          name: 'plugins',
          title: 'Extensions & Plugins',
        },
        {
          icon: 'icon cell-data',
          name: 'cell-data',
          title: 'Selected cell data detail view',
        },
        {
          icon: 'icon app',
          name: 'app',
          title: 'Application layers',
        },
      ],
      test: 'database'
    });

    function con(item) {
      setWithStorageVariableCache('selectedWidget', item.name)
      this.test = getWithStorageVariableCache('database', 'selectedWidget')
      emit('con', this.test)
    }

    return {
      ...toRefs(state), con
    };
  },
  methods: {
    //    con(item){
    //         setWithStorageVariableCache('selectedWidget',item.name)
    //         this.test = getWithStorageVariableCache('database','selectedWidget')
    //         console.log(item.name,'====')
    //    }
  },
  components: {FontIcon},
};
</script>

<style lang="less" scoped>
.wrapper {
  font-size: 23pt;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--theme-font-inv-2);
}

.wrapper:hover {
  color: var(--theme-font-inv-1);
}

.wrapper.selected {
  color: var(--theme-font-inv-1);
  background: var(--theme-bg-inv-3);
}

.panel_main {
  display: flex;
  flex: 1;
  flex-direction: column;
}
</style>


