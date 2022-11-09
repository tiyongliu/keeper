import {defineComponent, PropType} from 'vue'
import {databaseConnectionsSqlSelectApi} from '/@/api/simpleApis'
export default defineComponent({
  name: 'SqlFormView',
  props: {
    conid: {
      type: String as PropType<string>
    },
    database: {
      type: String as PropType<string>
    },
  },
  setup() {
    return () => (<div>FormView</div>)
  }
})

async function loadRow(props, select) {
  const { conid, database } = props;

  if (!select) return null;

  const response = await databaseConnectionsSqlSelectApi({
    conid,
    database,
    select,
  }) as any

  if (response.errorMessage) return response
  return response.rows[0]
}
