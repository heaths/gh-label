# gh label

[GitHub CLI] extension for issue label management.

## Install

Make sure you have version 2.0 or [newer] of the [GitHub CLI] installed.

```bash
gh extension install heaths/gh-label
```

### Upgrade

The `gh extension list` command shows if updates are available for extensions. To upgrade, you can use the `gh extension upgrade` command:

```bash
gh extension upgrade heaths/gh-label

# Or upgrade all extensions:
gh extension upgrade --all
```

## Commands

### create

Create a label in a repository.
You can specify colors with or without a preceeding hash ("#").
If you do not specify a color a random color will be choosen.

```bash
gh label create feedback
gh label create p1 --color e00808
gh label create p2 --color "#ffa501" --description "Affects more than a few users"
```

### delete

Delete a label from a repository.

```bash
gh label delete p1
```

### edit

Edit a label in a repository.
You can specify colors with or without a preceeding hash ("#").

```bash
gh label edit general --new-name feedback
gh label edit feedback --color c046ff --description "User feedback"
```

### export

Export labels from the repository to <path>, or stdout if <path> is "-".

```bash
gh label export ./labels.csv
gh label export ./labels.json
gh label export --format csv -
```

### list

List labels in a repository.
You can optionally pass a substring to match in the label name or description.

```bash
gh label list
gh label list service
```

## License

Licensed under the [MIT](LICENSE.txt) license.

Portions of this source copied from [vilmibm/gh-user-status](https://github.com/vilmibm/gh-user-status/tree/533285348c0354064d79053da39aa75f17b5c55f) under the [GNU Affero General Public License v3.0](https://github.com/vilmibm/gh-user-status/blob/533285348c0354064d79053da39aa75f17b5c55f/LICENSE).

[GitHub CLI]: https://github.com/cli/cli
[newer]: https://github.com/cli/cli/releases/latest
