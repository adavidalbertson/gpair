# `gpair`
[![Actions Status](https://github.com/adavidalbertson/gpair/workflows/go%20build/badge.svg)](https://github.com/adavidalbertson/gpair/actions)
[![Actions Status](https://github.com/adavidalbertson/gpair/workflows/go%20test/badge.svg)](https://github.com/adavidalbertson/gpair/actions)
[![Actions Status](https://github.com/adavidalbertson/gpair/workflows/golangci-lint/badge.svg)](https://github.com/adavidalbertson/gpair/actions)

`gpair` is a simple utility that makes it easier to share credit for collaboration on GitHub.
It stores the contact info of your frequent collaborators and adds a `Co-authored-by` clause to your default commit message.

## Usage
Assuming your current directory is in a git repository and you have added your coauthor's information (see the `add` subcommand for how to do that) under the alias `ALIAS`, run:

```
gpair ALIAS
```

This will append a `Co-authored-by` clause to the default commit message template for the repository.
The coauthor will appear on all commits to the repository until you run `gpair solo`.
The coauthor's information is saved so you can run `gpair ALIAS` again at any time to resume pairing.

To switch to a different coauthor, simply run `gpair ALIAS_2`, where `ALIAS_2` is the alias of another coauthor.

You can even credit multiple coauthors by running `gpair ALIAS_1 [ALIAS_2 ...]`

You can use the `--global` or `-g` flag to pair in global mode, for instance if you are working on multiple repos with the same coauthor.
Note that as with any git config, the local repo setting will override the global setting if present.

There are some additional flags you can pass in for more information on `gpair` or any subcommand:

* `-h` or `-help`: Display usage information
* `-v` or `-verbose`: Enable verbose output

## Subcommands
### `add`
Use the `add` subcommand to save a collaborator's information:

```
gpair add [ALIAS] NAME EMAIL
```

The positional arguments are as follows:

* `ALIAS`: An optional short name to refer to the collaborator by. If no alias is provided, `NAME` will be used instead.
* `NAME`: The collaborator's GitHub username
* `EMAIL`: The email address associated with the collaborator's GitHub account

You can also specify the arguments in any order using flags:

```
gpair add -email EMAIL -name NAME [-alias ALIAS]
```

To share credit with this collaborator, use `gpair ALIAS`.

### `remove`
Use the `remove` subcommand to remove a collaborator from the list:

```
gpair remove ALIAS
```

This will remove the coauthor from `gpair`'s list.

Note that if you have run `gpair` with this coauthor, they will still appear on commit messages until you run `gpair solo`.

### `solo`
Use the `solo` subcommand to end a pairing session.

```
gpair solo
```

This will reset the default commit message for the repository.

You can use the `--global` or `-g` flag to unpair if you previously used `gpair` in global mode.

## Installation

### Go Get
```
go get github.com/adavidalbertson/gpair
```

If your `$GOPATH` is on your path, that's all you need to do!

Of course, you'll need Go and Git installed already.

### Homebrew
```
brew install adavidalbertson/gpair/gpair
```

## How it works
`gpair` stores coauthor information and any applicable configuration settings in `~/.gpair/config.json`.
This file is created the first time `gpair` runs.

When you run `gpair ALIAS` in a repo, it creates a file `~/.gpair/REPO_NAME-template.txt` containing the coauthor's information, and sets git's `commit.template` config property to point to this file.
Subsequent uses of `gpair` will overwrite the template file.

`gpair solo` simply unsets git's `commit.template` property.

`gpair` aims to be nondestructive, so if your `commit.template` is set to a file *not* created by `gpair`, it will exit rather than overwrite the property.