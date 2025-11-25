# visualizedGit

 Visualize local git contribution powered by cobra and zap.

## Install

```shell
git clone https://github.com/zouyonghe/visualizedGit.git
cd visualizedGit
go install .
```

## Usage

```shell
visualizedGit --help                          # Show help information of this tool
visualizedGit add -p /path/to/repo(s)         # Add Git repository which you would like to track (defaults to current dir)
visualizedGit show -e email@address.com       # show visualized local git contribution (defaults to git config user.email)
visualizedGit chkcfg                          # Show current tracking git repositories
visualizedGit rmcfg                           # Remove tracking git repositories
```

## License

This software is power by the [GPLv2](./LICENSE) license.
