<template>
    <WidgetColumnBar :hidden="hidden">
        <div class="cell_wrapper">
            <WidgetTitle>Cell data view</WidgetTitle>
            <div class="cell_main">
                <div class="cell_toolbar">
                    Format:<span>&nbsp;</span>
                    <!-- <SelectField
                    isNative
                    :value="selectedFormatType"
                    @change="
                        (e) => {
                            selectedFormatType = e.detail;
                        }
                    "
                    :options="[
                        { value: 'autodetect', label: `Autodetect - ${autodetectFormat.title}` },
                        ...formats.map((fmt) => ({ label: fmt.title, value: fmt.type })),
                    ]"
                /> -->
                    <!-- <a-select @select="handleSelect" style="width:200px" size="small" v-model="selectedFormatType">
                        <a-select-option
                            :value="item.title"
                            v-for="(item, index) in formats"
                            :key="index"
                            :label="item.title"
                            ></a-select-option
                        >
                    </a-select> -->
                    <SelectModel
                        :default-value="{ key: selectedFormatType }"
                        @select="handleSelect"
                        style="width: 200px"
                        size="small"
                        v-model="selectedFormatType"
                        :options="[
                            {
                                value: 'autodetect',
                                label: `Autodetect - ${autodetectFormat().title}`,
                            },
                            ...formats.map((fmt) => ({ label: fmt.title, value: fmt.type })),
                        ]"
                    ></SelectModel>
                </div>
                <div class="cell_data">
                    <ErrorInfo
                        v-if="usedFormat().single && selection?.length != 1"
                        message="Must be selected one cell"
                        alignTop
                    />
                    <ErrorInfo
                        v-else-if="usedFormat() == null"
                        message="Format not selected"
                        alignTop
                    />
                    <ErrorInfo
                        v-else-if="!selection || selection.length == 0"
                        message="No data selected"
                        alignTop
                    />
                    <component
                        v-else
                        :is="usedFormat().component"
                        :selection="selection"
                    ></component>
                </div>
            </div>
        </div>
    </WidgetColumnBar>
</template>

<script lang="ts">
import { reactive, toRefs, markRaw } from 'vue';
import {isString, isPlainObject, isArray} from 'lodash-es';
import WidgetColumnBar from './WidgetColumnBar.vue';
// import SelectField from '../forms/SelectField.vue';
// import { Select } from 'ant-design-vue';
import ErrorInfo from '../elements/ErrorInfo.vue';
import WidgetTitle from './WidgetTitle.vue';
import SelectModel from '../../components/Form/src/components/SelectModel.vue';

import HtmlCellView from '../celldata/HtmlCellView.vue';
import TextCellViewWrap from '../celldata/TextCellViewWrap.vue';
import TextCellViewNoWrap from '../celldata/TextCellViewNoWrap.vue';
import JsonCellView from '../celldata/JsonCellView.vue';
import JsonRowView from '../celldata/JsonRowView.vue';
import PictureCellView from '../celldata/PictureCellView.vue';

export default {
    props:{
        hidden:{
            type:Boolean,
            default:false
        }
    },
    components: {
        ErrorInfo,
        WidgetTitle,
        // ASelect: Select,
        // ASelectOption: Select.Option,
        WidgetColumnBar,
        SelectModel,
        TextCellViewWrap: markRaw(TextCellViewWrap),
        TextCellViewNoWrap: markRaw(TextCellViewNoWrap),
        JsonCellView: markRaw(JsonCellView),
        JsonRowView: markRaw(JsonRowView),
        PictureCellView: markRaw(PictureCellView),
        HtmlCellView: markRaw(HtmlCellView),
    },
    setup() {
        const state = reactive({
            // hidden: false,
            selectedFormatType: 'autodetect',
            formats: [
                {
                    type: 'textWrap',
                    title: 'Text (wrap)',
                    component: TextCellViewWrap,
                    single: false,
                },
                {
                    type: 'text',
                    title: 'Text (no wrap)',
                    component: TextCellViewNoWrap,
                    single: false,
                },
                {
                    type: 'json',
                    title: 'Json',
                    component: JsonCellView,
                    single: true,
                },
                {
                    type: 'jsonRow',
                    title: 'Json - Row',
                    component: JsonRowView,
                    single: false,
                },
                {
                    type: 'picture',
                    title: 'Picture',
                    component: PictureCellView,
                    single: true,
                },
                {
                    type: 'html',
                    title: 'HTML',
                    component: HtmlCellView,
                    single: false,
                },
            ],
        });
        return {
            ...toRefs(state),
        };
    },
    computed: {
        // usedFormatType() {
        //     return this.selectedFormatType == 'autodetect'
        //         ? this.autodetectFormatType
        //         : this.selectedFormatType;
        // },
        // autodetectFormatType() {
        //     return this.autodetect(this.selection);
        // },
        // selection() {
        //     return this.handleSelect?this.handleSelect():[]
        // },
        // usedFormat() {
        //     return this.formats.find(x => x.type == this.usedFormatType);
        // },
        // autodetectFormat(){
        //     return this.formats.find(x => x.type == this.autodetectFormatType)
        // }
    },
    methods: {
        usedFormatType() {
            return this.selectedFormatType == 'autodetect'
                ? this.autodetectFormatType()
                : this.selectedFormatType;
        },
        autodetectFormatType() {
            return this.autodetect(this.selection());
        },
        selection() {
            return this.handleSelect ? this.handleSelect() : [];
        },
        usedFormat() {
            return this.formats.find((x) => x.type == this.usedFormatType());
        },
        autodetectFormat() {
            return this.formats.find((x) => x.type == this.autodetectFormatType());
        },
        autodetect(selection) {
            if (selection[0]?.engine?.databaseEngineTypes?.includes('document')) {
                return 'jsonRow';
            }
            const value = selection.length == 1 ? selection[0].value : null;
            if (isString(value)) {
                if (value.startsWith('[') || value.startsWith('{')) return 'json';
            }
            if (isPlainObject(value) || isArray(value)) {
                return 'json';
            }
            return 'textWrap';
        },
        handleSelect(val) {
            let arr = [];
            if (!val) {
                return arr;
            }
            this.selectedFormatType = val;
            return [val];
        },
    },
};
</script>

<style lang="less" scope>
.cell_wrapper {
    flex: 1;
    display: flex;
    flex-direction: column;
}

.cell_main {
    display: flex;
    flex: 1;
    flex-direction: column;
}

.cell_toolbar {
    display: flex;
    background: var(--theme-bg-1);
    align-items: center;
    border-bottom: 1px solid var(--thene-border);
    margin: 2px;
}

.cell_data {
    display: flex;
    flex: 1;
    position: relative;
}
</style>
