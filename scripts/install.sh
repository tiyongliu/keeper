#!/usr/bin/env bash

cd ../packages/tools
pnpm install
pnpm build

cd ../web
pnpm install
pnpm build

cd ../../
pnpm install

echo node_modules Successful installation!

echo The next step 'yarn start'
