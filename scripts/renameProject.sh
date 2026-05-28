#!/usr/bin/env bash

cd "$(dirname "$0")"

# CAUTION: this script will replace every occurrence of the word
# `cc-tools-demo` in the project folder with whatever argument
# you pass. Be very careful.

if [ $# -lt 1 ] ; then
  printf "Usage:\n$ ./renameProject.sh <my-project-name>\n"
  exit
fi

# Detect operating system and set sed options accordingly
if [[ "$OSTYPE" == "darwin"* ]]; then
  # macOS requires backup extension for sed -i
  grep -rl cc-tools-demo ../ --exclude-dir={.git,node_modules} | xargs sed -i '' s/cc-tools-demo/$1/g
else
  # Linux (Ubuntu and others) - no backup extension needed
  grep -rl cc-tools-demo ../ --exclude-dir={.git,node_modules} | xargs sed -i s/cc-tools-demo/$1/g
fi
