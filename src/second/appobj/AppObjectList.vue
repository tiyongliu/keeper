<template>

</template>

<script lang="ts">
  import {computed, defineComponent, PropType, unref} from 'vue';

  export default defineComponent({
    name: "DatabaseWidget",
    props: {
      list: {
        type: Array as unknown as PropType<[]>,
        default: [{
          "server": "localhost",
          "engine": "mysql@dbgate-plugin-mysql",
          "sshMode": "userPassword",
          "sshPort": "22",
          "sshKeyfile": "C:\\Users\\Administrator\\.ssh\\id_rsa",
          "port": "5001",
          "user": "root",
          "password": "crypt:73305a0992fb1e75e1a1c231ec53b56b27edec6cdf2f23a9a2247e3a6bfe04a7fcb9681bae0a889796817983aa178e75r/sewk0ogrcB9PTXu01Vzw==",
          "_id": "66436840-ad08-11ec-88e1-47c40318681c"
        }],
      }
    },
    components: {},
    setup(props) {
      const dynamicList = computed(() => unref(props.list))

      return {
        dynamicList
      }
    }
  })
</script>

<style scoped>

</style>
