<template>
  <div class="toolstrip">
    <slot name="toolstrip"></slot>
  </div>
  <div class="content">
    <slot></slot>
  </div>
</template>

<style scoped>
.content {
  position: absolute;
  left: 0;
  top: 0;
  border-bottom: 1px solid var(--theme-border);
  right: 0;
  bottom: var(--dim-toolstrip-height);
  display: flex;
}

.toolstrip {
  position: absolute;
  left: 0;
  height: var(--dim-toolstrip-height);
  right: 0;
  bottom: 0;
  display: flex;
  background: var(--theme-bg-1);
}
</style>
