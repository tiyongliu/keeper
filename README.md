# Wails + Vue 3 Typescript

## About

This is a Wails template project with Vue 3 and TypeScript, using Vite for
asset bundling. It comes with the bare minimum, and can be extended by following
the guides in this README. 

If you would like a more feature packed version that includes all features
documented below already added, please check out my
[feature packed Vite + Vue3 TypeScript template](https://github.com/codydbentley/wails-vite-vue-the-works)

## Live Development

To run in live development mode, run `wails dev` in the project directory. In another terminal, go into the `frontend`
directory and run `npm run dev`. Navigate to http://localhost:34115
in your browser to connect to your application.

Note: Typechecking is disabled. If you want to do type checking, use `npm run type-check`

## Extending Features

This template does not ship with things like routing, vuex, or sass.
To add any of these features, simply follow the instructions below. Please
note that all commands should be run in the `frontend` directory.

### Sass

Installation:
```shell
$ npm install --save-dev sass
```

Usage:

You can now add Sass to your single file component
styling like this:
```html
<style lang="scss">
  /* scss styling */
</style>
```

### ESLint + Prettier

Installation:
```shell
$ npm install --save-dev eslint prettier eslint-plugin-vue eslint-config-prettier @vue/eslint-config-typescript @typescript-eslint/parser @typescript-eslint/eslint-plugin
$ touch .eslintrc && touch .prettierrc
```

Usage: `eslintrc`
```json
{
  "extends": [
    "plugin:vue/vue3-essential",
    "eslint:recommended",
    "prettier",
    "@vue/typescript/recommended"
  ],
    "rules": {
    // override/add rules settings here, such as:
    // "vue/no-unused-vars": "error"
  }
}
```

Usage: `.prettierrc`
```json
{
  "semi": false,
  "tabWidth": 2,
  "useTabs": false,
  "printWidth": 120,
  "endOfLine": "auto",
  "singleQuote": true,
  "trailingComma": "all",
  "bracketSpacing": true,
  "arrowParens": "always"
}
```

### Vuex

Installation:
```shell
$ npm install --save vuex@next
$ touch src/store.ts
```

Usage: `src/store.ts`
```ts
import { InjectionKey } from 'vue'
import { createStore, Store, useStore as baseUseStore } from 'vuex'

// define your typings for the store state
export interface State {
  count: number
}

// define injection key
export const key: InjectionKey<Store<State>> = Symbol()

export const store = createStore<State>({
  state() {
    return {
      count: 0
    }
  },
  mutations: {
    increment(state) {
      state.count++
    }
  }
})

export function useStore() {
  return baseUseStore(key)
}
```

Usage: `src/main.ts`
```ts
import { createApp } from 'vue'
import App from './App.vue'
import { store, key } from './store'

createApp(App).use(store, key).mount('#app')
```

Usage: `src/components/Home.vue`
```ts
import { useStore } from '../store'
const store = useStore()
const increment = () => store.commit('increment')
```

### Vue Router

Installation:
```shell
$ npm install --save vue-router@4
$ touch src/router.ts
```

Usage: `src/router.ts`
```ts
import { createRouter, createWebHashHistory } from 'vue-router'
import Home from './components/Home.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home,
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

export default router
```

Usage: `src/main.ts`
```ts
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

createApp(App).use(router).mount('#app')
```

Usage: `src/App.vue`
```html
<template>
    <router-link to="/">Home</router-link>
    <router-view />
</template>
```

## Building 

To build this project in debug mode, use `wails build`. For production, use `wails build -production`.
To generate a platform native package, add the `-package` flag.

## Known Issues

- When making changes to the frontend, the browser reload will often happen too fast, causes issues. A refresh will fix the page.
- Typechecking is turned off due to Wails depending on the frontend to build before it will compile the backend and generate bindings.
- If you find any other problems, please create an issue.