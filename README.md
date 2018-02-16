# gosloc

# Install

```bash
go get github.com/GreenRaccoon23/gosloc;

```
# Download

```bash
git clone https://github.com/GreenRaccoon23/gosloc.git;
```

# Description

Command-line program to count source lines of code quickly.

Mostly works, but still in development.

```bash
[hiro@nakamura ~]$ gosloc -h;
gosloc <options> <path>...
  -i, --include string    File patterns to include, separated by commas
  -x, --exclude string    File patterns to exclude, separated by commas
  -c, --concurrency int   Max number of files to read simultaneously (default 1)
  -t, --total             Show a grand total, not the total for each file

WARNING: Setting concurrency too high will cause the program to crash,
corrupting the files it was editing.
```

Written in [Go](https://golang.org/).

# Binaries

I have compiled *(but not tested)* binaries for different platforms:

- [Linux 32-bit](./bin/linux_32/gosloc)
- [Linux 64-bit](./bin/linux_64/gosloc)
- [MacOS 32-bit](./bin/darwin_32/gosloc)
- [MacOS 64-bit](./bin/darwin_64/gosloc)
- [Windows 32-bit](./bin/windows_32/gosloc)
- [Windows 64-bit](./bin/windows_64/gosloc)
- Android 32-bit *(none yet)*
- Android 64-bit *(none yet)*

If you have [Go](https://golang.org/dl/) installed, it is easier and more efficient to install this package with the [install](#install) command. I added these binaries because I wanted to see if I could figure out how to make them.

There are no Android binaries because I have not figured out how to make them yet. I doubt anyone would want to run this program on an Android device anyway.

# Why

Similar tools I found were painfully slow, required awkward workarounds, and/or did not include everything I needed for my use cases. So I wrote this instead.

It is much faster than the other tools I have tried. It processes about 5000 files a second on a laptop which was considered mid-range performance 7 years ago.
