# Terraform module source version checker

A tool that validates and compares used vs. available terraform module versions
in git repositories, specific modules hosted in Gitlab repositories

## Install

```sh
brew tap "rollwagen/tf-module-versions" "https://github.com/rollwagen/tf-module-versions"
brew install rollwagen/tf-module-versions/tf-module-versions
```

## Usage

```text
A tool that validates and compares used vs. available terraform module versions
in git repositories, specific modules hosted in Gitlab repositories

Usage:
  tf-modver [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  validate    Print module version validation on stdout as logs

Flags:
  -h, --help   help for tf-modver

Use "tf-modver [command] --help" for more information about a command.
```

Example output of running `tf-module-versions validate`

<img width="901" alt="image" src="https://user-images.githubusercontent.com/7364201/179170042-5649e5f1-dc31-4e00-9a4d-8e4c7c5773df.png">
