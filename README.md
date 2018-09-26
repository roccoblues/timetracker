# tt [![Build Status](https://travis-ci.com/roccoblues/tt.svg?branch=master)](https://travis-ci.com/roccoblues/tt)

Command line time tracker.

## Installation

```
go get github.com/roccoblues/tt
```

## Usage

```
Usage:
  tt [command]

Available Commands:

  start	 Start a new timetracking interval
  stop	 Stop the current timetracking interval
  print	 Print the current timesheet (default command)

Flags:
  -file string
    	full path to data file (default "$HOME/.tt.json")
  -round-to int
    	round sum per day to minutes (default 15)
```

## Example output

```
$ tt
03.09.2018  5.25  09:00-13:30 17:30-18:15
04.09.2018  6.00  08:30-13:30 14:30-15:30
05.09.2018  4.00  10:30-13:30 14:30-15:30
06.09.2018  7.75  08:30-12:45 13:30-15:00 16:00-18:00
07.09.2018  1.25  11:30-12:15 13:00-13:30

10.09.2018  7.25  10:00-12:30 13:00-16:30 16:45-18:00
11.09.2018  8.00  08:30-12:14 12:40-16:59
12.09.2018  5.00  08:37-10:05 10:32-12:25 13:45-15:29
13.09.2018  7.25  08:26-12:34 13:33-13:42 14:39-17:37
14.09.2018  7.00  08:03-12:55 13:17-14:48 19:45-20:16
```

## Edit data

The data is saved by default in `~/.tt.json` and can be edited with your preferred editor. Example:

```
{
  "03.09.2018": [
    "09:00",
    "13:30",
    "17:30",
    "18:15"
  ],
  "04.09.2018": [
    "08:30",
    "13:30",
    "14:30",
    "15:30"
  ]
}
```

