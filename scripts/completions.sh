#!/bin/sh
CLI_NAME="tfm"
set -e
rm -rf completions
mkdir completions
for sh in bash zsh fish; do
	go run main.go completion "${sh}" >"completions/${CLI_NAME}.${sh}"
done
