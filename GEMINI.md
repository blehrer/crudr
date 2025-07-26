# `crudr list-paths` command

This feature will be one of many cobra (https://github.com/spf13/cobra) commands that present the user with a TUI for interacting with an openapi v3+ spec.

## Implementation details
- We will use https://github.com/pb33f/libopenapi to parse, render, update (etc) openapi specs.
- We will use https://github.com/charmbracelet/huh for the TUI.

## User story 1
1. User enters the following command
```sh
crudr list-paths -f relative/path/to/openapi.yaml
```
2. They are shown a list of all paths, including paths specified in documents linked by `$ref`, starting from the file passed with the flag `-f`
3. The command exits with code 0
