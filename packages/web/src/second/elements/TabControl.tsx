import {defineComponent, ref, unref, watch, onBeforeUnmount, toRefs} from 'vue';
import {Tabs, TabPane} from 'ant-design-vue'
import {compact} from 'lodash-es'
// const TabPane = Tabs.TabPane
interface TabDef {
  label: string;
  slot?: number;
  component?: any;
  props?: any;
}

export default defineComponent({
  name: 'TabControl',
  props: {
    isInline: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    tabs: {
      type: Array as PropType<TabDef[]>
    }
  },
  setup(props) {
    const {tabs} = toRefs(props)

    const tabsVNode = () => compact(tabs.value).map((tab, index) => {
      console.log(tab, `--`, index)
      return <TabPane key={index} tab={tab.label} />
    })

    return () => (
      <div>
        <div>
          <Tabs defaultActiveKey="1">
            {tabsVNode()}
          </Tabs>
        </div>

        <div></div>
      </div>
    )
  }
})
