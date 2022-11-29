<template>
  <BasicModal
    @register="register"
    title="Add connection">
    <!-- Modal content -->
    <div class="window" :class="fullScreen && 'fullScreen'">

    </div>
  </BasicModal>
</template>

<script lang="ts">
import {defineComponent, PropType, toRefs} from 'vue';
import {BasicModal, useModalInner} from '/@/components/Modal';
import FontIcon from '/@/second/icons/FontIcon.vue'

export default defineComponent({
  name: 'ModalBase',
  props: {
    fullScreen: {
      type: Boolean as PropType<boolean>,
      default: false,
    }
  },
  components: {
    BasicModal,
    FontIcon
  },
  setup(props) {
    const [register, {closeModal, setModalProps}] = useModalInner()

    return {
      ...toRefs(props),
      register,
      closeModal
    }
  }
})
</script>

