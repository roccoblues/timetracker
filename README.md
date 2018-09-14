# tt

Command line time tracker. Times are rounded to 15 minutes.

## Installation

```
go get github.com/roccoblues/timetracker/tt
```

## Usage

### Start a new time tracking interval
```
$ tt start
01.09.2018  0.00  8:00-
```

### Output all saved times
```
$ tt
01.09.2018  0.00  8:00-
```

### Stop a time tracking interval
```
$ tt stop
01.09.2018  4.00  08:00-12:00
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
11.09.2018  8.00  08:30-12:15 12:45-17:00
12.09.2018  5.00  08:30-10:00 10:30-12:30 13:45-15:30
13.09.2018  7.25  08:30-12:30 13:30-13:45 14:45-17:30
14.09.2018  0.00  08:00-
```

## Edit data

The data is saved in `~/.tt.json` and can be edited with your preferred editor.
