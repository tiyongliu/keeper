<template>
  <div class="widgetTitle">
    <slot></slot>
  </div>
</template>

<style lang="less">
.widgetTitle {
  padding: 5px;
  font-weight: bold;
  text-transform: uppercase;
  background-color: var(--theme-bg-1);
  border: 2px solid var(--theme-border);

  &:hover {
    background-color: var(--theme-bg-2);
  }
}
</style>

