import {defineComponent, onMounted, PropType, ref, toRefs, watch} from 'vue'
import {InputNumber} from 'ant-design-vue'

export default defineComponent({
  name: 'TextField',
  props: {
    size: {
      type: String as PropType<'large' | 'small'>,
    },
    type: {
      type: String as PropType<'number' | 'text'>,
      default: 'text'
    },
    focused: {
      type: Boolean as PropType<boolean>,
      default: false
    },
    value: {
      type: Number as PropType<number>
    },
    autocomplete: {
      type: String as PropType<string>,
      default: 'new-password'
    },

  },
  emits: ['update:value', 'change', 'input', 'click', 'keydown'],
  setup(props, {attrs, emit}) {
    const {size, type, focused, value, autocomplete} = toRefs(props)
    const valueRW = ref(value.value)

    watch(() => valueRW.value, () => {
      emit('update:value', valueRW.value)
    })

    const domEditor = ref<Nullable<HTMLElement>>(null);

    onMounted(() => {
      if (focused.value && domEditor.value) {
        domEditor.value.focus()
      }
    })

    return () => (
      <>
        {
          type.value === 'number' ?
            <InputNumber
              ref={domEditor}
              {...attrs}
              size={size.value}
              vModel:value={valueRW.value}
              onChange={e => emit('keydown', e)}
              onInput={e => emit('input', e)}
              onClick={e => emit('click', e)}
              onKeydown={e => emit('keydown', e)}
              autocomplete={autocomplete.value}
            /> :
            <a-input
              ref={domEditor}
              {...attrs}
              size={size.value}
              vModel:value={valueRW.value}
              onChange={e => emit('keydown', e)}
              onInput={e => emit('input', e)}
              onClick={e => emit('click', e)}
              onKeydown={e => emit('keydown', e)}
              autocomplete={autocomplete.value}
            />
        }
      </>
    )
  }
})
