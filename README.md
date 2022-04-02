# Goal
To implement a stopwatch and a timer.

# Usage
`wtc [-help] [duration]`

Run `wtc` to start a stopwatch.

Specify duration with `wtc` to start a timer. Duration must be in [[hh:]mm:]ss format,
such as,

- `wtc 5`       start a 5 second timer
- `wtc 1200`    start a 20 minute timer
- `wtc 120:00`  start a 2 hour timer
- `wtc 4:32`    start a 4 minute and 32 seconds timer
- `wtc 1:23:00` start a one hour, 23 minute timer

End of the timer is followed by a chime.
