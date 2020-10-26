# :bulb: commit-analyzer-cz
[![CI](https://github.com/go-semantic-release/commit-analyzer-cz/workflows/CI/badge.svg?branch=master)](https://github.com/go-semantic-release/commit-analyzer-cz/actions?query=workflow%3ACI+branch%3Amaster)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-semantic-release/commit-analyzer-cz)](https://goreportcard.com/report/github.com/go-semantic-release/commit-analyzer-cz)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/go-semantic-release/commit-analyzer-cz)](https://pkg.go.dev/github.com/go-semantic-release/commit-analyzer-cz)

A [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) analyzer for [go-semantic-release](https://github.com/go-semantic-release/semantic-release).

## How the commit messages are analyzed

### Bump major version (0.1.2 -> 1.0.0)
- By adding `BREAKING CHANGE` or `BREAKING CHANGES` in the commit message footer, e.g.:
  ```
  feat: allow provided config object to extend other configs

  BREAKING CHANGE: `extends` key in config file is now used for extending other config files
  ```
- By adding `!` at the end of the commit type, e.g.:
  ```
  refactor!: drop support for Node 6
  ```

### Bump minor version (0.1.2 -> 0.2.0)
- By using type `feat`, e.g.:
  ```
  feat(lang): add polish language
  ```

### Bump patch version (0.1.2 -> 0.1.3)
- By using type `fix`, e.g.:
  ```
  fix: correct minor typos in code

  see the issue for details

  on typos fixed.

  Reviewed-by: Z
  Refs #133
  ```

## References
- [Conventional Commit v1.0.0 - Examples](https://www.conventionalcommits.org/en/v1.0.0/#examples)

## Licence

The [MIT License (MIT)](http://opensource.org/licenses/MIT)

Copyright © 2020 [Christoph Witzko](https://twitter.com/christophwitzko)
