# watch
TUI for a stopwatch/timer terminal application.

## Usage

    watch [-help] [duration]

## Stopwatch
A bare

```shell
$ watch
```

without any arguments, starts a stopwatch. You can also take laps, and copy
them onto your clipboard.

## Timer
Specify duration with `watch` to start a timer. Duration must be in
`[[hh:]mm:]ss` format, that is,

- a duration of 5       starts a 5 second timer
- a duration of 1200    starts a 20 minute timer
- a duration of 120:00  starts a 2 hour timer
- a duration of 4:32    starts a 4 minute and 32 seconds timer
- a duration of 1:23:00 starts a one hour, 23 minute timer

You can queue multiple timers like so,

```shell
$ watch 1 2 3 4
```

This starts 1 second timer which would be followed by a 2, 3 and 4 second
timer.

End of the timer is followed by a chime.
