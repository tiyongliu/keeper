<template>
  <div class="widgetsInnerContainer"><slot /></div>
</template>

<style lang="less">
  .widgetsInnerContainer {
    flex: 1 1;
    overflow-x: auto;
    overflow-y: auto;
    width: var(--dim-left-panel-width);
  }
</style>
