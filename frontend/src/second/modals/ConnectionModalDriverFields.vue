<template>
  <a-form layout="vertical">
    <a-form-item label="Database engine">
      <a-select
        @change="handleSelect"
        placeholder="please select your zone"
        :options="databaseEngine"
        v-model:value="engine"/>
    </a-form-item>

    <a-form-item label="Database file" v-if="false">
      <a-row type="flex" justify="space-between" align="top">
        <a-col :span="12">
          <a-input/>
        </a-col>
        <a-col :span="12">
          <a-button type="primary">browse</a-button>
        </a-col>
      </a-row>
    </a-form-item>

    <a-form-item label="Resources" v-if="false">
      <a-radio-group v-model:value="resources" name="radioGroup" :options="[
         { label: 'Fill database connection details', value: '' },
         { label: 'Use database URL', value: '1' },
      ]"/>
    </a-form-item>

    <a-row type="flex" justify="space-between" align="top">
      <a-col :span="16">
        <a-form-item label="Server" v-bind="validateInfos.host">
          <a-input v-model:value="driverForm.host" autocomplete="off"/>
        </a-form-item>
      </a-col>
      <a-col :span="8">
        <a-form-item label="Port">
          <a-input
            v-model:value="driverForm.port"
            :placeholder="driver && driver.defaultPort"
            autocomplete="off"/>
        </a-form-item>
      </a-col>
    </a-row>

    <a-row type="flex" justify="space-between" align="top">
      <a-col :span="12">
        <a-form-item label="User">
          <a-input v-model:value="driverForm.username" autocomplete="off"/>
        </a-form-item>
      </a-col>
      <a-col :span="12">
        <a-form-item label="Password">
          <a-input-password v-model:value="driverForm.password" placeholder=""/>
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


    <a-form-item label="Database URL" v-if="false">
      <a-input/>
    </a-form-item>

    <a-form-item label="Default database">
      <a-input/>
    </a-form-item>

    <a-form-item label="Display name">
      <a-input v-model:value="driverForm.displayName"/>
    </a-form-item>
  </a-form>

</template>

<script lang="ts">
import {
  computed,
  defineComponent,
  inject,
  onMounted,
  reactive,
  ref,
  toRefs,
  toRaw,
  unref,
  watch
} from "vue"
import {useDebounceFn} from '@vueuse/core'
import {cloneDeep} from 'lodash-es'
import {
  Button,
  Checkbox,
  Col,
  Form,
  FormItem,
  Input,
  Radio,
  RadioGroup,
  Row,
  Select,
  SelectOption,
} from 'ant-design-vue'
import $extensions from './drivers.json'
const InputPassword = Input.Password
const useForm = Form.useForm

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
      {label: '(select schema)', value: ''},
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

    const driverForm = reactive<{ [key in string]: string | null | undefined } & { port: string | number }>({
      engine: '',
      host: 'localhost',
      username: '',
      password: '',
      port: '',
      displayName: null
    })

    const driver = computed(() => {
      return $extensions.drivers.find(x => x.engine == unref(engine))
    })

    const engine = ref($values.engine)
    const dispatchConnections = inject('dispatchConnections') as any

    const rulesRef = reactive({
      host: [{ required: true, message: 'Please input your username!' }],
    })
    const { resetFields, validate, validateInfos } = useForm(driverForm, rulesRef,{
      onValidate: (...args) => console.log(...args),
    })

    const notificationTest = async () => {
      try {
        await validate()
        const dynamicProps = cloneDeep(toRaw(driverForm))
        const [shortName] = unref(engine).split('@')
        dynamicProps.engine = shortName
        if (!dynamicProps.port) {
          dynamicProps.port = `${driver.value!.defaultPort}`
        }
        dispatchConnections(dynamicProps)
      } catch (e) {
        console.log(e)
      }
    }

    const handleSelect = () => resetFields()

    watch(() => [unref(driver), toRefs(driverForm)],
      useDebounceFn(() => notificationTest(), 300),
      {deep: true}
    )

    onMounted(() => {
      void notificationTest()
    })

    return {
      databaseEngine,
      engine,
      resources: ref(''),
      driver,
      driverForm,
      handleSelect,
      validateInfos
    }

  }
})
</script>
