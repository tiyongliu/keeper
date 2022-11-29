<template>
    <div>
        HtmlCellView
        <!-- {{
            selection
                .map((cell) => {
                    const { value } = cell;
                    if (_.isPlainObject(value) || _.isArray(value))
                        return JSON.stringify(value, undefined, 2);
                    return cell.value;
                })
                .join('\n')
        }} -->
    </div>
</template>

<script>
import { reactive, toRefs } from 'vue';
export default {
    props:{
        selection:{
            type:Array,
            default:[]
        }
    },
    setup() {
        const state = reactive({
            count: 0,
        });

        return {
            ...toRefs(state),
        };
    },
};
</script>
