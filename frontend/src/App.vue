<template>
  <Home />
</template>

<script setup lang="ts">
// This template is using Vue 3 <script setup> SFCs
// Check out https://v3.vuejs.org/api/sfc-script-setup.html#sfc-script-setup
import Home from './components/Home.vue'
</script>

<style>
@font-face {
  font-family: "Nunito";
  font-style: normal;
  font-weight: 400;
  src: local(""),
  url("assets/fonts/nunito-v16-latin-regular.woff2") format("woff2");
}

html {
  background-color: rgba(33, 37, 43, 1);
  color: white;
  font-family: "Nunito", -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto",
  "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue",
  sans-serif;
}

body {
  margin: 0;
}

#app {
  width: 100vw;
  height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.link {
  cursor: pointer;
  color: #df0000;
}

.link:hover {
  text-decoration: underline;
}
</style>
