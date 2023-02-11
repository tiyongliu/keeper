<template>
  <div
    class="outer buttonLike"
    :title="title"
    :class="[disabled && 'disabled', square && 'square', narrow && 'narrow']"
    ref="domButton">
    <div class="inner">
      <slot></slot>
    </div>
  </div>
</template>

<script lang="ts">
  import {defineComponent, ref} from 'vue';
  export default defineComponent({
    name: "InlineButton",
    props: {
      title: {
        type: String as PropType<string>,
      },
      disabled: {
        type: Boolean as PropType<boolean>,
        default: false
      },
      square: {
        type: Boolean as PropType<boolean>,
        default: false
      },
      narrow: {
        type: Boolean as PropType<boolean>,
        default: false
      }
    },
    setup(props) {
      const domButton = ref<Nullable<HTMLElement>>(null)
      const hidden = ref(false)

      function getBoundingClientRect() {
        return domButton.value!.getBoundingClientRect()
      }

      return {
        domButton,
        hidden,
        ...props,
        getBoundingClientRect,
      }
    }
  })
</script>

<style scoped>
.outer {
  --bg-1: var(--theme-bg-1);
  --bg-2: var(--theme-bg-3);

  background: linear-gradient(to bottom, var(--bg-1) 5%, var(--bg-2) 100%);
  background-color: var(--bg-1);
  border: 1px solid var(--bg-2);
  display: inline-block;
  cursor: pointer;
  vertical-align: middle;
  color: var(--theme-font-1);
  font-size: 12px;
  padding: 3px;
  margin: 0;
  text-decoration: none;
  display: flex;
}

.narrow {
  padding: 3px 1px;
}

.outer.disabled {
  color: var(--theme-font-3);
}

.outer:hover:not(.disabled) {
  border: 1px solid var(--theme-font-1);
}

.outer:active:not(.disabled) {
  background: linear-gradient(to bottom, var(--bg-2) 5%, var(--bg-1) 100%);
  background-color: var(--bg-2);
}

.inner {
  margin: auto;
  flex: 1;
  text-align: center;
}

.square {
  width: 18px;
}
</style>
