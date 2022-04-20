import {defineComponent} from 'vue'
import DatabaseAppObject from './DatabaseAppObject'
export const extractKey = props => props.name

export default defineComponent({
  name: 'PinnedAppObject',
  setup(props) {
    // const props = {
    //   checkedObjectsStore: null,
    //   data: {
    //     connection: {
    //       engine: "mysql@dbgate-plugin-mysql",
    //       password: "crypt:7000413edf483ada3770dc5c4b9a69f0beea98f82c2e3b9ba243488a63c0fc056ee70323004cbfe3b5438a7297fcdfe3LC25uegcuz6H5UxZfY2UyA==",
    //       server: "localhost",
    //       sshKeyfile: "C:\\Users\\Administrator\\.ssh\\id_rsa",
    //       sshMode: "userPassword",
    //       sshPort: "22",
    //       status: {
    //         name: "ok"
    //       },
    //       user: "root",
    //       _id: "065caa90-a8c6-11ec-9b4b-6f98950c4d7a"
    //     }
    //   },
    //   disableContextMenu: false,
    //   expandIcon: null
    // }

    return () => (
      <DatabaseAppObject
        {...props}
      />
    )
  }
})
