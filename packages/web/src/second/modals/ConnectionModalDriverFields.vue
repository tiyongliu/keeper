<template>
  <!--  <div>ConnectionModalDriverFields</div>-->

  <a-form>
    <a-form-item label="Database engine">
      <a-select placeholder="please select your zone">
        <a-select-option value="shanghai">Zone one</a-select-option>
        <a-select-option value="beijing">Zone two</a-select-option>
      </a-select>
    </a-form-item>

    <a-form-item label="Resources">
      <a-radio-group>
        <a-radio value="1">Fill database connection details</a-radio>
        <a-radio value="2">Use database URL</a-radio>
      </a-radio-group>
    </a-form-item>

    <a-row type="flex" justify="center" align="top">
      <a-col :span="12">
        <a-form-item label="Server">
          <a-input/>
        </a-form-item>
      </a-col>
      <a-col :span="12">
        <a-form-item label="Port">
          <a-input/>
        </a-form-item>
      </a-col>
    </a-row>


    <a-row type="flex" justify="center" align="top">
      <a-col :span="12">
        <a-form-item label="User">
          <a-input/>
        </a-form-item>
      </a-col>
      <a-col :span="12">
        <a-form-item label="Password">
          <a-input-password placeholder="input password" />
        </a-form-item>
      </a-col>
    </a-row>

    <a-form-item label="Password mode">
      <a-select placeholder="please select your zone">
        <a-select-option value="saveEncrypted">Save and encrypt</a-select-option>
        <a-select-option value="saveRaw">Save raw (UNSAFE!!)</a-select-option>
      </a-select>
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
import {defineComponent} from "vue"
import {Col, Form, FormItem, Input, InputPassword, Radio, RadioGroup, Row, Select, SelectOption, Checkbox} from 'ant-design-vue'

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

  }
})
</script>

<style scoped>

</style>
