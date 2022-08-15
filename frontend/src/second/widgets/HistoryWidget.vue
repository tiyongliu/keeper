<template>
    <WidgetColumnBar :hidden="hidden">
        HistoryWidget---没写
    </WidgetColumnBar>
</template>

<script>
import { reactive, toRefs } from 'vue'
import WidgetColumnBar from './WidgetColumnBar.vue'
export default {
    props:{
        hidden:{
            type:Boolean,
            default:false
        }
    },
    setup () {
        const state = reactive({
            count: 0,
        })
    
        return {
            ...toRefs(state),
        }
    },
    components:{WidgetColumnBar}
}
</script>

<style lang="scss" scoped>

</style>