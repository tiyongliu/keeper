<template>
  <FormProviderCore>
    <BasicModal ref="ConnectionModal" class="connectionModal" @register="register"
                title="Add connection">
      <TabControl isInline :tabs="tabs"/>
      <template #insertFooter>
        <a-button class="float-left" type="default" @click="handleTest">测试连接</a-button>
        {{sqlConnectResult}}
      </template>
    </BasicModal>
  </FormProviderCore>
</template>

<script lang="ts">
import {defineComponent, onMounted, reactive, ref} from 'vue'
import {Alert, Tabs} from 'ant-design-vue'
// import {useModal} from '/@/components/Modal'
import {BasicModal, useModalInner} from '/@/components/Modal'
import FormProviderCore from '/@/second/forms/FormProviderCore'
import TabControl from '/@/second/elements/TabControl.vue'
import ConnectionModalDriverFields from '/@/second/modals/ConnectionModalDriverFields.vue'
import ConnectionModalSshTunnelFields from '/@/second/modals/ConnectionModalSshTunnelFields.vue'
import ConnectionModalSslFields from '/@/second/modals/ConnectionModalSslFields.vue'
import {handleDriverTestApi} from '/@/api/connection'

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
    const connectionModal = ref()
    const sqlConnectResult = ref(null)

    onMounted(() => {
      console.log(connectionModal.value, `connectionModal`)
    })

    const handleTest = async () => {
      const resp = await handleDriverTestApi({
        engine: "mongo",
        server: "localhost",
        port: "27017"
      })

      // const resp = await handleDriverTestApi({
      //   engine: "mysql",
      //   password: "123456",
      //   server: "localhost",
      //   sshKeyfile: "/Users/liuliutiyong/.ssh/id_rsa",
      //   sshMode: "userPassword",
      //   sshPort: "22",
      //   user: "root",
      //   port: "3306"
      // })

      // sqlConnectResult.value = resp
      console.log(resp, `resp`)
    }

    return {
      register,
      closeModal,
      handleTest,
      setModalProps: () => {
        //bodyStyle
        setModalProps({title: 'Modal New Title', bodyStyle: {padding: `0`}});
      },
      sqlConnectResult,
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
      connectionModal
    }
  }
})
</script>

<style lang="less" scoped>
::v-deep(.scrollbar.scroll-container) {
  padding: 3px;
}
</style>
