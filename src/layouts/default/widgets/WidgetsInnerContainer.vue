<template>
  <div><slot /></div>
</template>

<style scoped>
  div {
    flex: 1 1;
    overflow-x: auto;
    overflow-y: auto;
    width: var(--dim-left-panel-width);
  }
</style>
