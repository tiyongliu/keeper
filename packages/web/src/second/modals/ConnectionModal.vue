<template>
  <FormProviderCore>
    <BasicModal @register="register" title="Add connection">

<!--      <a-tabs default-active-key="1">-->
<!--        <a-tab-pane key="1" tab="Main" />-->
<!--        <a-tab-pane key="2" tab="SSH Tunnel" />-->
<!--        <a-tab-pane key="3" tab="SSL" />-->
<!--      </a-tabs>-->

      <TabControl isInline :tabs="tabs"/>

    </BasicModal>
  </FormProviderCore>
</template>

<script lang="ts">
import {defineComponent} from 'vue'
import {Tabs} from 'ant-design-vue'
// import {useModal} from '/@/components/Modal'
import {BasicModal, useModalInner} from '/@/components/Modal'
import FormProviderCore from '/@/second/forms/FormProviderCore'
import TabControl from '/@/second/elements/TabControl'

const TabPane = Tabs.TabPane
export default defineComponent({
  name: 'ConnectionModal',
  components: {
    FormProviderCore,
    TabControl,
    BasicModal,
    [Tabs.name]: Tabs,
    [TabPane.name]: TabPane
  },
  emits: ['register'],
  setup() {
    const [register, {closeModal, setModalProps}] = useModalInner()

    return {
      register,
      closeModal,
      setModalProps: () => {
        setModalProps({ title: 'Modal New Title' });
      },

      tabs: [
        {
          label: 'Main',
        },
        {
          label: 'SSH Tunnel',
        },
        {
          label: 'SSL',
        },
      ]
    }
  }
})
</script>

<style scoped></style>
