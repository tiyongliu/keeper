## 参考表格

```tsx
import {defineComponent, h, onMounted, ref} from 'vue'
import {Progress} from 'ant-design-vue';
import {BasicColumn, BasicTable, useTable} from '/@/components/Table';
import ToolStripContainer from '/@/second/buttons/ToolStripContainer.vue'
import {useMessage} from "/@/hooks/web/useMessage";
import eb_system_config from "/@/second/tabs/eb_system_config.json";
const columns: BasicColumn[] = [
  {
    title: '输入框',
    dataIndex: 'name',
    edit: true,
    editComponentProps: {
      prefix: '$',
    },
  },
  {
    title: '默认输入状态',
    dataIndex: 'name7',
    edit: true,
    editable: true,
  },
  {
    title: '输入框校验',
    dataIndex: 'name1',
    edit: true,
    // 默认必填校验
    editRule: true,
  },
  {
    title: '数字输入框',
    dataIndex: 'id',
    edit: true,
    editRule: true,
    editComponent: 'InputNumber',
    editComponentProps: () => {
      return {
        max: 100,
        min: 0,
      };
    },
    editRender: ({text}) => {
      return h(Progress, {percent: Number(text)});
    },
  },


  {
    title: '开关',
    dataIndex: 'name6',
    edit: true,
    editComponent: 'Switch',
    editValueMap: (value) => {
      return value ? '开' : '关';
    },
  },
]
export default defineComponent({
  name: 'CollectionDataTab',
  setup() {
    const [registerTable, {setColumns, setTableData}] = useTable({
      columns: [],
      tableSetting: { fullScreen: true },
      showIndexColumn: false,
      showTableSetting: true,
    });

    const {createMessage} = useMessage();

    function handleEditEnd({record, index, key, value}: Recordable) {
      console.log(record, index, key, value);
      return false;
    }

    // 模拟将指定数据保存
    function feakSave({value, key, id}) {
      createMessage.loading({
        content: `正在模拟保存${key}`,
        key: '_save_fake_data',
        duration: 0,
      });
      return new Promise((resolve) => {
        setTimeout(() => {
          if (value === '') {
            createMessage.error({
              content: '保存失败：不能为空',
              key: '_save_fake_data',
              duration: 2,
            });
            resolve(false);
          } else {
            createMessage.success({
              content: `记录${id}的${key}已保存`,
              key: '_save_fake_data',
              duration: 2,
            });
            resolve(true);
          }
        }, 2000);
      });
    }

    async function beforeEditSubmit({record, index, key, value}) {
      console.log('单元格数据正在准备提交', {record, index, key, value});
      return await feakSave({id: record.id, key, value});
    }

    function handleEditCancel() {
      console.log('cancel');
    }

    onMounted(() => {
      setColumns(eb_system_config.columns.map(column => {
        return {
          title: column.columnName,
          dataIndex: column.columnName,
          edit: true,
          editable: false,
        }
      }))

      setTableData(eb_system_config.rows)
    })

    return () => (
      <ToolStripContainer>
        <BasicTable
          onRegister={registerTable}
          onEditEnd={handleEditEnd}
          onEditCancel={handleEditCancel}
          beforeEditSubmit={beforeEditSubmit}/>
      </ToolStripContainer>
    )
  },
})


export const matchingProps = ['conid', 'database', 'schemaName', 'pureName'];
export const allowAddToFavorites = _ => true

```
