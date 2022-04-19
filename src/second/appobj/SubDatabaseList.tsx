import {computed, defineComponent, onMounted, PropType, unref, reactive} from 'vue'
import {map, sortBy} from 'lodash-es'
import {filterName} from '/@/packages/tools/src'
import './SubDatabaseList.less'
import FontIcon from '../icons/FontIcon.vue'
import AppObjectList from './AppObjectList'
import databaseAppObject from './DatabaseAppObject'
export default defineComponent({
  name: "SubDatabaseList",
  props: {
    passProps: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    data: {
      type: Object as PropType<{}>
    },
    filter: {
      type: String as PropType<string>,
      default: ''
    }
  },
  setup(props) {
    const {data, filter, passProps} = props
    const databases = computed((): {name: string, sortOrder?: string}[] => {
      return [{"name":"information_schema"}, {"name":"crmeb_java_beta"},{"name":"mysql"},{"name":"performance_schema"},{"name":"sys"}]
    })




    console.log(computed(() => {
      return sortBy(
        (unref(databases) || []).filter(x => filterName(unref(filter!), x.name)),
        x => x.sortOrder ?? x.name
      ).map(db => ({ ...db, connection: unref(data) }))
    }))

    onMounted(() => {

    })

    const tables = [{"name":"crmeb","connection":{"server":"localhost","engine":"mysql@dbgate-plugin-mysql","sshMode":"userPassword","sshPort":"22","sshKeyfile":"C:\\Users\\Administrator\\.ssh\\id_rsa","user":"root","password":"crypt:7000413edf483ada3770dc5c4b9a69f0beea98f82c2e3b9ba243488a63c0fc056ee70323004cbfe3b5438a7297fcdfe3LC25uegcuz6H5UxZfY2UyA==","_id":"065caa90-a8c6-11ec-9b4b-6f98950c4d7a","status":{"name":"ok"}}},{"name":"erd","connection":{"server":"localhost","engine":"mysql@dbgate-plugin-mysql","sshMode":"userPassword","sshPort":"22","sshKeyfile":"C:\\Users\\Administrator\\.ssh\\id_rsa","user":"root","password":"crypt:7000413edf483ada3770dc5c4b9a69f0beea98f82c2e3b9ba243488a63c0fc056ee70323004cbfe3b5438a7297fcdfe3LC25uegcuz6H5UxZfY2UyA==","_id":"065caa90-a8c6-11ec-9b4b-6f98950c4d7a","status":{"name":"ok"}}},{"name":"information_schema","connection":{"server":"localhost","engine":"mysql@dbgate-plugin-mysql","sshMode":"userPassword","sshPort":"22","sshKeyfile":"C:\\Users\\Administrator\\.ssh\\id_rsa","user":"root","password":"crypt:7000413edf483ada3770dc5c4b9a69f0beea98f82c2e3b9ba243488a63c0fc056ee70323004cbfe3b5438a7297fcdfe3LC25uegcuz6H5UxZfY2UyA==","_id":"065caa90-a8c6-11ec-9b4b-6f98950c4d7a","status":{"name":"ok"}}},{"name":"kb-dms","connection":{"server":"localhost","engine":"mysql@dbgate-plugin-mysql","sshMode":"userPassword","sshPort":"22","sshKeyfile":"C:\\Users\\Administrator\\.ssh\\id_rsa","user":"root","password":"crypt:7000413edf483ada3770dc5c4b9a69f0beea98f82c2e3b9ba243488a63c0fc056ee70323004cbfe3b5438a7297fcdfe3LC25uegcuz6H5UxZfY2UyA==","_id":"065caa90-a8c6-11ec-9b4b-6f98950c4d7a","status":{"name":"ok"}}},{"name":"mallplusbak","connection":{"server":"localhost","engine":"mysql@dbgate-plugin-mysql","sshMode":"userPassword","sshPort":"22","sshKeyfile":"C:\\Users\\Administrator\\.ssh\\id_rsa","user":"root","password":"crypt:7000413edf483ada3770dc5c4b9a69f0beea98f82c2e3b9ba243488a63c0fc056ee70323004cbfe3b5438a7297fcdfe3LC25uegcuz6H5UxZfY2UyA==","_id":"065caa90-a8c6-11ec-9b4b-6f98950c4d7a","status":{"name":"ok"}}},{"name":"mysql","connection":{"server":"localhost","engine":"mysql@dbgate-plugin-mysql","sshMode":"userPassword","sshPort":"22","sshKeyfile":"C:\\Users\\Administrator\\.ssh\\id_rsa","user":"root","password":"crypt:7000413edf483ada3770dc5c4b9a69f0beea98f82c2e3b9ba243488a63c0fc056ee70323004cbfe3b5438a7297fcdfe3LC25uegcuz6H5UxZfY2UyA==","_id":"065caa90-a8c6-11ec-9b4b-6f98950c4d7a","status":{"name":"ok"}}},{"name":"performance_schema","connection":{"server":"localhost","engine":"mysql@dbgate-plugin-mysql","sshMode":"userPassword","sshPort":"22","sshKeyfile":"C:\\Users\\Administrator\\.ssh\\id_rsa","user":"root","password":"crypt:7000413edf483ada3770dc5c4b9a69f0beea98f82c2e3b9ba243488a63c0fc056ee70323004cbfe3b5438a7297fcdfe3LC25uegcuz6H5UxZfY2UyA==","_id":"065caa90-a8c6-11ec-9b4b-6f98950c4d7a","status":{"name":"ok"}}},{"name":"schema","connection":{"server":"localhost","engine":"mysql@dbgate-plugin-mysql","sshMode":"userPassword","sshPort":"22","sshKeyfile":"C:\\Users\\Administrator\\.ssh\\id_rsa","user":"root","password":"crypt:7000413edf483ada3770dc5c4b9a69f0beea98f82c2e3b9ba243488a63c0fc056ee70323004cbfe3b5438a7297fcdfe3LC25uegcuz6H5UxZfY2UyA==","_id":"065caa90-a8c6-11ec-9b4b-6f98950c4d7a","status":{"name":"ok"}}},{"name":"shop_go","connection":{"server":"localhost","engine":"mysql@dbgate-plugin-mysql","sshMode":"userPassword","sshPort":"22","sshKeyfile":"C:\\Users\\Administrator\\.ssh\\id_rsa","user":"root","password":"crypt:7000413edf483ada3770dc5c4b9a69f0beea98f82c2e3b9ba243488a63c0fc056ee70323004cbfe3b5438a7297fcdfe3LC25uegcuz6H5UxZfY2UyA==","_id":"065caa90-a8c6-11ec-9b4b-6f98950c4d7a","status":{"name":"ok"}}},{"name":"sql_join","connection":{"server":"localhost","engine":"mysql@dbgate-plugin-mysql","sshMode":"userPassword","sshPort":"22","sshKeyfile":"C:\\Users\\Administrator\\.ssh\\id_rsa","user":"root","password":"crypt:7000413edf483ada3770dc5c4b9a69f0beea98f82c2e3b9ba243488a63c0fc056ee70323004cbfe3b5438a7297fcdfe3LC25uegcuz6H5UxZfY2UyA==","_id":"065caa90-a8c6-11ec-9b4b-6f98950c4d7a","status":{"name":"ok"}}},{"name":"ssodb","connection":{"server":"localhost","engine":"mysql@dbgate-plugin-mysql","sshMode":"userPassword","sshPort":"22","sshKeyfile":"C:\\Users\\Administrator\\.ssh\\id_rsa","user":"root","password":"crypt:7000413edf483ada3770dc5c4b9a69f0beea98f82c2e3b9ba243488a63c0fc056ee70323004cbfe3b5438a7297fcdfe3LC25uegcuz6H5UxZfY2UyA==","_id":"065caa90-a8c6-11ec-9b4b-6f98950c4d7a","status":{"name":"ok"}}},{"name":"yami_shops","connection":{"server":"localhost","engine":"mysql@dbgate-plugin-mysql","sshMode":"userPassword","sshPort":"22","sshKeyfile":"C:\\Users\\Administrator\\.ssh\\id_rsa","user":"root","password":"crypt:7000413edf483ada3770dc5c4b9a69f0beea98f82c2e3b9ba243488a63c0fc056ee70323004cbfe3b5438a7297fcdfe3LC25uegcuz6H5UxZfY2UyA==","_id":"065caa90-a8c6-11ec-9b4b-6f98950c4d7a","status":{"name":"ok"}}}]



    const onPin = (e) => {
      e.preventDefault();
      e.stopPropagation();

      //pinnedDatabases
     console.log(data, 'data-data')

    }

    // return () => (
    //   tables.map(item => {
    //     return <div class="main" draggable="true">
    //       <FontIcon icon="mdi mdi-database color-icon-gold" />{item.name}
    //           <span class="pin">
    //          <FontIcon icon="mdi mdi-pin" onClick={onPin}/>
    //       </span>
    //     </div>
    //   })
    // )


    return () => (
      <AppObjectList
        module={databaseAppObject}
        list={sortBy(
          (unref(databases) || []).filter(x => filterName(unref(filter!), x.name)),
          x => x.sortOrder ?? x.name
        ).map(db => ({ ...db, connection: unref(data) }))}
        {...unref(passProps)}
      />
    )
  }
})
