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

## Customizable Release Rules
It is possible to customize the release rules by providing options to the analyzer. The following options are available:

| Option                | Default |
|-----------------------|---------|
| `major_release_rules` | `*!`    |
| `minor_release_rules` | `feat`  |
| `patch_release_rules` | `fix`   |

⚠️ Commits that contain `BREAKING CHANGE(S)` in their body will always result in a major release. This behavior cannot be customized yet.

### Rule Syntax
A rule may match a specific commit type, scope or both. The following syntax is supported: `<type>(<scope>)<modifier>`

- `<type>`: The commit type, e.g. `feat`, `fix`, `refactor`.
- `<scope>`: The commit scope, e.g. `lang`, `config`. If left empty, the rule matches all scopes (`*`).
- `<modifier>`: The modifier, e.g. `!` for breaking changes. If left empty, the rule matches only commits without a modifier.
- A `*` may be used as a wildcard for a type, scope or modifier.

### Example Rules
| Commit                             | `feat` (or `feat(*)` | `*!` (or `*(*)!`) | `chore(deps)` | `*🚀` |
|------------------------------------|----------------------|-------------------|---------------|-------|
| `feat(ui): add button component`   | ✅                    | ❌                 | ❌             | ❌     |
| `feat!: drop support for Go 1.17`  | ❌                    | ✅                 | ❌             | ❌     |
| `chore(deps): update dependencies` | ❌                    | ❌                 | ✅             | ❌     |
| `refactor: remove unused code`     | ❌                    | ❌                 | ❌             | ❌     |
| `fix🚀: correct minor typos`       | ❌                    | ❌                 | ❌             | ✅     |


## References
- [Conventional Commit v1.0.0 - Examples](https://www.conventionalcommits.org/en/v1.0.0/#examples)

## Licence

The [MIT License (MIT)](http://opensource.org/licenses/MIT)

Copyright © 2024 [Christoph Witzko](https://twitter.com/christophwitzko)
