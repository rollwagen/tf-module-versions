# Terraform module source version checker

![CodeQL](https://github.com/rollwagen/tf-module-versions/workflows/CodeQL/badge.svg)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)

![image](https://user-images.githubusercontent.com/7364201/180030995-cd871550-4c58-4abf-9554-acd1d5c14cd7.png)

A tool that validates and compares used vs. available terraform module versions in git repositories, specific modules hosted in Gitlab repositories

## Install

```sh
brew tap rollwagen/homebrew-tap
brew install rollwagen/tap/tf-module-versions
```

## Usage

### Pre-requisites

Currently, version validation support is only for terraform modules stored
in Gitlab repositories.
For authentication towards Gitlab, an environment variable `GITLAB_TOKEN`
needs to be present that holds a valid GitLab authentication token.

```text
A tool that validates and compares used vs. available terraform module versions
in git repositories, specific modules hosted in Gitlab repositories

Usage:
  tfm [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  validate    Print module version validation on stdout as logs

Flags:
  -h, --help   help for tf-modver

Use "tfm [command] --help" for more information about a command.
```

Example output of running `tfm validate`

![tfm](https://user-images.githubusercontent.com/7364201/180036688-e8b43e06-a085-453f-97a6-f90672685a7a.gif)
