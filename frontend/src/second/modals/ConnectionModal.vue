<template>
  <FormProviderCore>
    <BasicModal
      @register="register"
      @cancel="handleCancelTest"
      @ok="handleSubmit"
      class="connectionModal"
      title="Add connection">
      <TabControl isInline :tabs="tabs"/>
      <template #insertFooter>
        <a-button class="float-left" type="default" @click="handleTest">测试连接</a-button>
      </template>
    </BasicModal>
  </FormProviderCore>
</template>

<script lang="ts">
import {defineComponent, onMounted, reactive, ref, provide, unref} from 'vue'
import {pickBy} from 'lodash-es'
import {Alert, Tabs} from 'ant-design-vue'
// import {useModal} from '/@/components/Modal'
import {BasicModal, useModalInner} from '/@/components/Modal'
import FormProviderCore from '/@/second/forms/FormProviderCore'
import TabControl from '/@/second/elements/TabControl.vue'
import ConnectionModalDriverFields from '/@/second/modals/ConnectionModalDriverFields.vue'
import ConnectionModalSshTunnelFields from '/@/second/modals/ConnectionModalSshTunnelFields.vue'
import ConnectionModalSslFields from '/@/second/modals/ConnectionModalSslFields.vue'
import {handleDriverTestApi, handleDriverSaveApi} from '/@/api/connection'

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
  emits: ['register'],
  setup() {
    const [register, {closeModal, setModalProps}] = useModalInner()
    let connParams = reactive<{[key in string]: any}>({})

    provide('dispatchConnections', (dynamicProps) => {
      console.log(`dynamicProps`, dynamicProps)
      connParams = dynamicProps
    })

    const handleTest = async () => {
      const resp = await handleDriverTestApi(pickBy(unref(connParams), (item) => !!item))
      console.log(resp, `resp`)
    }

    const handleCancelTest = () => {}

    const handleSubmit = async () => {
      const resp = await handleDriverSaveApi(pickBy(unref(connParams), (item) => !!item))
      console.log(resp, `resp`)
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
