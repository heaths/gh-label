# gh label

[GitHub CLI] extension for issue label management.

## Install

Make sure you have version 2.0 or newer of the [GitHub CLI] installed.

```bash
gh extension install heaths/gh-label
```

## Commands

### create

Create a label in a repository.

```bash
gh label create feedback
gh label create p1 --color e00808
gh label create p2 --color "#ffa501" --description "Affects more than a few users"
```

### edit

Edit a label in a repository.

```bash
gh label edit general --new-name feedback
gh label edit feedback --color c046ff --description "User feedback"
```

### delete

Delete a label from a repository.

```bash
gh label delete p1
```

### list

List labels in a repository. You can optionally pass a substring to match in the label name or description.

```bash
gh label list
gh label list service
```

## License

Licensed under the [MIT](LICENSE.txt) license.

Portions of this source copied from [vilmibm/gh-user-status](https://github.com/vilmibm/gh-user-status/tree/cead3abf46ffb5fd3c178a0ba6f2c69c3dbabf7e) under the [GNU Affero General Public License v3.0](https://github.com/vilmibm/gh-user-status/blob/cead3abf46ffb5fd3c178a0ba6f2c69c3dbabf7e/LICENSE).

[GitHub CLI]: https://github.com/cli/cli
