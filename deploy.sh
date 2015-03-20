#!/bin/bash
set -e
npm install
npm prune
npm build
goapp deploy