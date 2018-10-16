# gospace

gospace is a workspace mangament tool for flattening vendor directory to a cohesive `GOPATH`.

This tool is necessary for linters that do not handle modules properly. Once Go tooling has caught up with modules, this tool/repository will be deprecated.

An example setup for `storj.io/storj` would set these environment variables:

```
# set go modules to default behavior
GO111MODULE=auto

# go knows where our gopath is
GOPATH=~/storj

# gospace knows where our gopath is (this is to avoid accidental damage to existing GOPATH)
# you should not use default GOPATH here
GOSPACE_ROOT=~/storj

# set the github repository that this GOSPACE manages
GOSPACE_PKG=storj.io/storj

# set the where the repository is located
GOSPACE_REPO=git@github.com:storj/storj.git
```

## Setup

First time when setting up your local system.

```
gospace setup
```

`gospace setup` does:

1. Delete `$GOSPACE_ROOT/bin`, `$GOSPACE_ROOT/src` folders.
2. Download `$GOSPACE_REPO` repository into `$GOSPACE_ROOT/src/$GOSPACE_PKG` folder.
3. Run `GO111MODULE=on go mod vendor` inside `$GOSPACE_ROOT/src/$GOSPACE_PKG` folder. This downloads all dependencies into vendor directory.
4. Moves all vendored directories to `$GOSPACE_ROOT/src`.

## Updating

Every time go.mod changes, you should run:

```
gospace update
```

`gospace update` does:

1. Delete all non-repository directories in `$GOSPACE_ROOT/src`, effectively deleting all vendored directories. Unless you have placed something manually, which is also deleted.
2. Run `GO111MODULE=on go mod vendor` inside `$GOSPACE_ROOT/src/$GOSPACE_PKG` folder. This downloads all dependencies into vendor directory.
3. Moves all vendored directories to `$GOSPACE_ROOT/src`.

## Low-level commands

To calculate hash of the `go.mod` and `go.sum` file:

```
gospace hash
```

To download dependencies and zip vendor folder:

```
gospace zip-vendor vendor.zip
```

To unzip a vendor directory:

```
gospace unzip-vendor vendor.zip
```

To flatten vendor directory:

```
gospace flatten-vendor
```

# License

This code is licensed under [Apache v2](https://www.apache.org/licenses/LICENSE-2.0) license.