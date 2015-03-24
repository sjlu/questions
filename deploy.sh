#!/bin/bash
set -e
npm install
npm prune
npm run build
goapp deploy
