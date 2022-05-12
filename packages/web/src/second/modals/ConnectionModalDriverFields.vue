<template>
  <!--  <div>ConnectionModalDriverFields</div>-->

  <a-form layout="vertical">
    <a-form-item label="Database engine">
      <a-select placeholder="please select your zone" :options="engine"/>
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
          <a-input/>
        </a-form-item>
      </a-col>
      <a-col :span="8">
        <a-form-item label="Port">
          <a-input :placeholder="driver && driver.defaultPort"/>
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
import {defineComponent, ref} from "vue"
import {
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
  SelectOption
} from 'ant-design-vue'
import $extensions from './drivers.json'

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
  },
  setup() {
    const electron = null

    const engine = [
      {label: '(select driver)', value: ''},
      ...$extensions.drivers
        .filter(driver => !driver.isElectronOnly || electron)
        .map(driver => ({
          value: driver.engine,
          label: driver.title,
        })),
    ]

    return {
      engine,
      resources: ref(''),
      driver: $extensions.drivers.find(x => x.engine == engine)
    }

  }
})
</script>

<style scoped>

</style>
