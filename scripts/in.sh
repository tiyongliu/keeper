#!/usr/bin/env bash

cd ../packages/tools
pnpm install
pnpm build

cd ../../frontend
pnpm install
pnpm build

echo node_modules Successful installation!

echo The next step 'cd frontend'

echo The next step 'yarn start'