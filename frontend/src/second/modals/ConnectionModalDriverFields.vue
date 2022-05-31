<template>
  <a-form layout="vertical">
    <a-form-item label="Database engine">
      <a-select placeholder="please select your zone" :options="databaseEngine"
                v-model:value="driverForm.engine"/>
    </a-form-item>
    <a-form-item label="Database file">
      <a-row type="flex" justify="space-between" align="top">
        <a-col :span="12">
          <a-input/>
        </a-col>
        <a-col :span="12">
          <a-button type="primary">browse</a-button>
        </a-col>
      </a-row>
    </a-form-item>

    <a-form-item label="Resources">
      <a-radio-group v-model:value="resources" name="radioGroup" :options="[
         { label: 'Fill database connection details', value: '' },
         { label: 'Use database URL', value: '1' },
      ]"/>
    </a-form-item>

    <a-row type="flex" justify="space-between" align="top">
      <a-col :span="16">
        <a-form-item label="Server">
          <a-input v-model:value="driverForm.server"/>
        </a-form-item>
      </a-col>
      <a-col :span="8">
        <a-form-item label="Port">
          <a-input v-model:value="driverForm.server" :placeholder="driver && driver.defaultPort"/>
        </a-form-item>
      </a-col>
    </a-row>


    <a-row type="flex" justify="space-between" align="top">
      <a-col :span="12">
        <a-form-item label="User">
          <a-input/>
        </a-form-item>
      </a-col>
      <a-col :span="12">
        <a-form-item label="Password">
          <a-input-password placeholder="input password"/>
        </a-form-item>
      </a-col>
    </a-row>

    <a-form-item label="Password mode">
      <a-select placeholder="please select your zone" :options="[
      { value: 'saveEncrypted', label: 'Save and encrypt' },
      { value: 'saveRaw', label: 'Save raw (UNSAFE!!)' },
    ]"/>
    </a-form-item>

    <a-form-item label="">
      <a-checkbox value="1" name="type">Is read only</a-checkbox>
    </a-form-item>

    <a-form-item label="Database URL">
      <a-input/>
    </a-form-item>
  </a-form>
</template>

<script lang="ts">
import {computed, defineComponent, reactive, ref, unref} from "vue"
import {
  Button,
  Checkbox,
  Col,
  Form,
  FormItem,
  Input,
  InputPassword,
  Radio,
  RadioGroup,
  Row,
  Select,
  SelectOption,
} from 'ant-design-vue'
import $extensions from './drivers.json'

import {handleDriverTestApi} from '/@/api/connection'

export default defineComponent({
  name: 'ConnectionModalDriverFields',
  components: {
    [Form.name]: Form,
    [FormItem.name]: FormItem,
    [Select.name]: Select,
    [SelectOption.name]: SelectOption,
    [Radio.name]: Radio,
    [RadioGroup.name]: RadioGroup,
    [Input.name]: Input,
    [Row.name]: Row,
    [Col.name]: Col,
    [InputPassword.name]: InputPassword,
    [Checkbox.name]: Checkbox,
    [Button.name]: Button,
  },
  setup() {
    const electron = null

    const databaseEngine = [
      {label: '(select driver)', value: ''},
      ...$extensions.drivers
        .filter(driver => !driver.isElectronOnly || electron)
        .map(driver => ({
          value: driver.engine,
          label: driver.title,
        })),
    ]

    const $values = {
      "server": "localhost",
      "engine": "mysql@dbgate-plugin-mysql",
      "sshMode": "userPassword",
      "sshPort": "22",
      "sshKeyfile": "/Users/liuliutiyong/.ssh/id_rsa",
      "useDatabaseUrl": ""
    }

    const driverForm = reactive({
      engine: '',
      server: '',
      port: '',
    })

    const engine = ref($values.engine)

    const handleTest = () => {
      handleDriverTestApi({
        engine: "mongo",
        server: "localhost",
        port: "27017"
      })
      // handleDriverTestApi({
      //   engine: "mysql",
      //   password: "123456",
      //   server: "localhost",
      //   sshKeyfile: "/Users/liuliutiyong/.ssh/id_rsa",
      //   sshMode: "userPassword",
      //   sshPort: "22",
      //   user: "root",
      //   port: "3306"
      // })
    }

    const driver = computed(() => {
      return $extensions.drivers.find(x => x.engine == unref(engine))
    })

    return {
      databaseEngine,
      engine,
      resources: ref(''),
      driver,
      handleTest,
      driverForm
    }

  }
})
</script>

<style scoped>

</style>
