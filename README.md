# Goal
To implement a stopwatch and a timer.

# Usage
`watch [-help] [duration]`

Run `watch` to start a stopwatch.

Specify duration with `watch` to start a timer. Duration must be in [[hh:]mm:]ss format,
such as,

- `watch 5`       start a 5 second timer
- `watch 1200`    start a 20 minute timer
- `watch 120:00`  start a 2 hour timer
- `watch 4:32`    start a 4 minute and 32 seconds timer
- `watch 1:23:00` start a one hour, 23 minute timer

End of the timer is followed by a chime.
