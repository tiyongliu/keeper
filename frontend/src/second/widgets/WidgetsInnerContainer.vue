<template>
  <div class="widgetsInnerContainer">
    <slot></slot>
  </div>
</template>

<style scoped>
.widgetsInnerContainer {
  flex: 1 1;
  overflow-x: auto;
  overflow-y: auto;
  width: var(--dim-left-panel-width);
}
</style>
