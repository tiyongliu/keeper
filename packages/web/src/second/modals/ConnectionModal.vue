<template>
  <FormProviderCore>
    <BasicModal ref="ConnectionModal" class="connectionModal" @register="register" title="Add connection">
      <TabControl isInline :tabs="tabs"/>
    </BasicModal>
  </FormProviderCore>
</template>

<script lang="ts">
import {defineComponent, ref, onMounted} from 'vue'
import {Tabs} from 'ant-design-vue'
// import {useModal} from '/@/components/Modal'
import {BasicModal, useModalInner} from '/@/components/Modal'
import FormProviderCore from '/@/second/forms/FormProviderCore'
import TabControl from '/@/second/elements/TabControl.vue'
import ConnectionModalDriverFields from '/@/second/modals/ConnectionModalDriverFields.vue'
import ConnectionModalSshTunnelFields from '/@/second/modals/ConnectionModalSshTunnelFields.vue'
import ConnectionModalSslFields from '/@/second/modals/ConnectionModalSslFields.vue'

const TabPane = Tabs.TabPane
export default defineComponent({
  name: 'ConnectionModal',
  components: {
    FormProviderCore,
    TabControl,
    BasicModal,
    [Tabs.name]: Tabs,
    [TabPane.name]: TabPane,
  },
  emits: ['register'],
  setup() {
    const [register, {closeModal, setModalProps}] = useModalInner()
    const connectionModal = ref()


    onMounted(() => {
      console.log(connectionModal.value, `connectionModal`)
    })

    return {
      register,
      closeModal,
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
