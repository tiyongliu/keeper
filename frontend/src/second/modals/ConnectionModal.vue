<template>
  <FormProviderCore>
    <BasicModal
      @register="register"
      @cancel="handleCancelTest"
      @ok="handleSubmit"
      class="connectionModal"
      width="50%"
      title="Add connection">
      <TabControl isInline :tabs="tabs"/>
      <template #insertFooter>
        <a-button class="float-left" type="default" @click="handleTest">测试连接</a-button>
      </template>
    </BasicModal>
  </FormProviderCore>
</template>

<script lang="ts">
import {defineComponent, provide, unref} from 'vue'
import {storeToRefs} from 'pinia'
import {pickBy} from 'lodash-es'
import {Alert, Tabs} from 'ant-design-vue'
import {BasicModal, useModalInner} from '/@/components/Modal'
import FormProviderCore from '/@/second/forms/FormProviderCore'
import TabControl from '/@/second/elements/TabControl.vue'
import ConnectionModalDriverFields from '/@/second/modals/ConnectionModalDriverFields.vue'
import ConnectionModalSshTunnelFields from '/@/second/modals/ConnectionModalSshTunnelFields.vue'
import ConnectionModalSslFields from '/@/second/modals/ConnectionModalSslFields.vue'
import {connectionTestApi, connectionSaveApi} from '/@/api/simpleApis'

import {useBootstrapStore} from "/@/store/modules/bootstrap"
const TabPane = Tabs.TabPane

export default defineComponent({
  name: 'ConnectionModal',
  components: {
    FormProviderCore,
    TabControl,
    BasicModal,
    [Tabs.name]: Tabs,
    [TabPane.name]: TabPane,
    [Alert.name]: Alert,
  },
  emits: ['register', 'closeCurrentModal'],
  setup(_, {emit}) {
    const [register, {closeModal, setModalProps}] = useModalInner()
    let connParams = {}
    const bootstrap = useBootstrapStore()
    const {connections} = storeToRefs(bootstrap)
    provide('dispatchConnections', (dynamicProps) => {
      connParams = dynamicProps
    })

    const handleTest = async () => {
      try {
        await connectionTestApi(pickBy(unref(connParams), (item) => !!item))
      } catch (e) {
        console.log(e)
      }
    }

    const handleCancelTest = () => {}

    const handleSubmit = async () => {
      try {
        const resp = await connectionSaveApi(pickBy(unref(connParams), (item) => !!item))
        void bootstrap.setConnections([...unref(connections), resp])
        emit('closeCurrentModal')
      } catch (e) {
        console.log(e)
      }
    }

    return {
      register,
      closeModal,
      handleTest,
      handleCancelTest,
      handleSubmit,
      setModalProps: () => {
        //bodyStyle
        setModalProps({title: 'Modal New Title', bodyStyle: {padding: `0`}});
      },
      tabs: [
        {
          label: 'Main',
          component: ConnectionModalDriverFields
        },
        {
          label: 'SSH Tunnel',
          component: ConnectionModalSshTunnelFields,
        },
        {
          label: 'SSL',
          component: ConnectionModalSslFields
        },
      ],
      bodyStyle: {
        padding: `0`
      },
    }
  }
})
</script>

<style lang="less" scoped>
::v-deep(.scrollbar.scroll-container) {
  padding: 3px;
}
</style>
