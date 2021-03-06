# tt [![Build Status](https://github.com/roccoblues/tt/workflows/Test/badge.svg)](https://github.com/roccoblues/tt/actions)

tt is a simple command line time tracker.

## Installation

```
go get github.com/roccoblues/tt
```

## Usage

```
Usage: ./tt [flags] [start|stop] [time]

  -date-format string
    	parse and write dates with format (default "02.01.2006")
  -file string
    	path to data file (default "$HOME/.tt.json")
  -month int
    	output month (default current)
  -round-to int
    	round to minutes (default 15)
  -time-format string
    	parse and write times with format (default "15:04")
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

## FAQ

### Help, I forgot to start/stop the timer.

If you just forgot the most recent event you can call `start`/`stop` with an optional time to fix it. If it's already the next day you need to manually [edit the data file](#edit-data).

### I need to track times for different client/projects.

Just use a different file `-f FILE` for each client.
