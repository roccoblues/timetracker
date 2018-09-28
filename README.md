# tt [![Build Status](https://travis-ci.com/roccoblues/tt.svg?branch=master)](https://travis-ci.com/roccoblues/tt)

tt is a command line time tracker. The times in the output are rounded to 15 minute intervals.


## Installation

```
go get github.com/roccoblues/tt
```

## Usage

```
Usage:
  tt [flags]
  tt [command]

Available Commands:
  help        Help about any command
  start       Start a new timetracking interval
  stop        Stop the current timetracking interval

Flags:
  -f, --file FILE   path to data FILE (default "/Users/dennis/.tt.json")
  -h, --help        help for tt

Use "tt [command] --help" for more information about a command.
```

## Example output

```
$ tt
03.09.2018  8.50   09:00-13:30 14:15-18:15
04.09.2018  5.00   08:30-13:30 14:15-

Total: 13.50
```

## Edit data

The data is saved by default in `~/.tt.json` and can be edited with your preferred editor. Example:

```
{
  "03.09.2018": [
    "09:00",
    "13:30",
    "14:17",
    "18:15"
  ],
  "04.09.2018": [
    "08:30",
    "13:30",
    "14:16"
  ]
}
```

