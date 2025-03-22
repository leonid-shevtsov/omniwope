# OmniWOPE

_Write Once, Post Everywhere_

This tool lets you publish the content of your blog to various platforms. This helps you reach your readers without forcing them to use a specific

The source of truth should be a site that you own. Any source is supported as long as you can export the articles into a special JSON file.

From there, you run the tool after publishing the next article to the blog and it takes care of rebroadcasting it to every configured Output.

## Supported Sources

- [Hugo](sources/hugo/README.md)

_Adding a source is as easy as generating JSON of the [input structure](docs/input_structure.md). PRs are welcome!_

## Supported Outputs

- [Mastodon](docs/mastodon.md) or API-compatible services (GoToSocial and more)
- [Telegram channel](docs/tg.md)

## How to run

```
Omniwope - Write Once Publish Everywhere

Usage:
  omniwope [flags]

Flags:
      --config string   config file (default is omniwope.yml)
      --dry-run         dry-run (log changes instead of applying)
  -h, --help            help for omniwope
      --verbose         enable debug logging
```

## How to develop

### Project directory structure

```
/cmd - entry points for commands (this project uses Cobra)
/internal
  /output - implementations of outputs
    /tg - Telegram
  /store - implementations of the data store
    /json - JSON store
/sources - guides to retrieve input
  /hugo - from a Hugo blog
```
