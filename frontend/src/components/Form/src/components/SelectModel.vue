<template>
    <a-select
        ref="refSelectModel"
        v-if="!editable"
        v-model="selectVal"
        :placeholder="placeholderVal"
        :size="size"
        :disabled="disabled"
        @change="change"
        class="select-model"
        :collapse-tags="collapseTags"
        :clearable="clearable"
        filterable
        :default-value="{ key: defaultSheet }"
    >
        <a-select-option
            v-for="(item, index) in options"
            :key="index"
            :label="item.label"
            :value="item.value"
            :disabled="item.disabled || false"
        >
            <slot v-bind:item="item"></slot>
        </a-select-option>
    </a-select>
    <span v-else>{{ text }}</span>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import { Select } from 'ant-design-vue';

export default defineComponent({
    name: 'SelectModel',
    components:{ASelect: Select,ASelectOption: Select.Option,},
    props: {
        value: {
            default:undefined
        },
        disabled: Boolean,
        size: {
            default: 'medium',
        },
        placeholder: {
            default: 'Select',
        },
        type: {
            default: '',
        },
        clearable: {
            default: true,
        },
        editable: {
            default: false,
        },
        preOptions: {
            type: Array,
            default: function () {
                return [];
            },
        },
        appendOptions: {
            type: Array,
            default: function () {
                return [];
            },
        },
        filter: {
            default: '',
        },
        filterBy: String,
        selectAll: {
            default: false,
        },
        collapseTags: {
            default: false,
        },
        optionsArr: {
            type: Array,
            default: function () {
                return [];
            },
        },
        defaultSheet:{
            default:''
        }
    },
    model: {
        prop: 'value',
        event: 'update',
    },
    setup() {
        return {
            selectVal: ''
        };
    },
    methods: {
        optionsFilter(options) {
            let arr = [];
            if (this.filter && this.filter !== '') {
                arr = options.map((item) => {
                    let temp;
                    let length = item.bindingTypes.length;
                    if (length > 0) {
                        for (let i = 0, j = length; i < j; i++) {
                            let binding = item.bindingTypes[i];
                            if (binding && binding.indexOf(this.filter.toUpperCase()) > -1) {
                                temp = item;
                                break;
                            } else {
                                temp = null;
                            }
                        }
                    } else {
                        temp = null;
                    }
                    return temp;
                });
            } else {
                return options;
            }
            return arr.filter((item) => item);
        },
        change(val) {
            if ( Array.isArray(val)) {
                let hasAll = val.filter((item) => {
                    return item === '**';
                });
                if (hasAll.length > 0) {
                    val = this.options.map((item) => {
                        if (item.value !== '**') {
                            return item.value;
                        }
                    });
                    val = val.filter((item) => item);
                }
            }
            this.selectVal = val
            this.$emit('update', val);
            this.$emit('change', val);
        },
        focus() {
            this.$refs.refSelectModel.focus();
        },
    },
    computed: {
        options() {
            return [...this.preOptions, ...this.appendOptions];
        },
        placeholderVal() {
            return this.disabled ? '' : this.placeholder;
        },
        text() {
            let res = this.options.find((v) => {
                return v.value === this.value;
            });
            return res ? res.label : '/';
        },
    },
    watch: {
        value(val, oldVal) {
            this.selectVal = val;
        },
    },
});
</script>

<style scoped>
</style>

