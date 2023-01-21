<template>
  <div class="wrapper">
    <InlineButton @click="handleSkipMinus">
      <FontIcon icon="icon arrow-left"/>
    </InlineButton>
    <span class="label">Start:</span>
    <TextField
      size="small"
      type="number"
      v-model:value="skipRW"
      @blur="handleLoad"
      @keydown="handleKeyDown"
    />
    <span class="label">Rows:</span>
    <TextField
      size="small"
      type="number"
      v-model:value="limitRW"
      @blur="handleLoad"
      @keydown="handleKeyDown"
    />
    <InlineButton @click="handleSkipPlus">
      <FontIcon icon="icon arrow-right"/>
    </InlineButton>
  </div>
</template>

<script lang="ts">
import {defineComponent, PropType, ref, toRefs, watch} from 'vue'
import {isNumber} from '/@/utils/is'
import FontIcon from '/@/second/icons/FontIcon.vue'
import InlineButton from '/@/second/buttons/InlineButton.vue'
import TextField from '/@/second/forms/TextField'
import keycodes from '/@/second/utility/keycodes'

export default defineComponent({
  name: "Pager",
  components: {
    FontIcon,
    InlineButton,
    TextField
  },
  props: {
    skip: {
      type: Number as PropType<number>,
    },
    limit: {
      type: Number as PropType<number>,
    }
  },
  emits: ['update:skip', 'update:limit', 'load'],
  setup(props, {emit}) {
    const {skip, limit} = toRefs(props)

    const skipRW = ref(skip)
    const limitRW = ref(limit)


    function handleSkipPlus() {
      if (isNumber(skipRW.value) && isNumber(limitRW.value)) {
        skipRW.value = parseInt(skipRW.value) + parseInt(limitRW.value)
        if (skipRW.value < 0) skipRW.value = 0
        emit('load')
      }
    }

    function handleSkipMinus() {
      if (isNumber(skipRW.value) && isNumber(limitRW.value)) {
        skipRW.value = parseInt(skipRW.value) - parseInt(limitRW.value)
        emit('load')
      }
    }

    function handleKeyDown(e) {
      if (e.keyCode == keycodes.enter) {
        e.preventDefault();
        e.stopPropagation();
        emit('load')
      }
    }

    watch(() => skipRW.value, () => {
      if (isNumber(skipRW.value) && skipRW.value >= 0) emit('update:skip', skipRW.value)
    })
    watch(() => limitRW.value, () => {
      if (isNumber(limitRW.value) && limitRW.value >= 0) emit('update:limit', limitRW.value)
    })

    function handleLoad() {
      emit('load')
    }

    return {
      skipRW,
      limitRW,
      handleLoad,
      handleSkipPlus,
      handleSkipMinus,
      handleKeyDown
    }
  }
})
</script>

<style scoped>
.wrapper :global(input) {
  width: 100px;
}

.wrapper {
  display: flex;
  align-items: center;
}

.label {
  margin-left: 5px;
  margin-right: 5px;
}
</style>
