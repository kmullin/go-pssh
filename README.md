# PSSH

A simple parallel ssh written in Go.

Simply pipe in a hostlist to standard in (stdin), and provide it a command and any arguments you want to run on every host.

### Example

    $ echo admin{01..10}.localdomain | gpssh 'uptime'
    admin05.localdomain:  23:45:04 up 957 days, 0 min,  2 users,  load average: 2.16, 2.70, 3.10
    admin06.localdomain:  23:45:04 up 956 days, 23:35,  5 users,  load average: 7.83, 8.95, 10.15
    admin04.localdomain:  23:45:04 up 957 days,  1:01,  2 users,  load average: 12.84, 17.20, 17.09
    admin08.localdomain:  23:45:04 up 207 days,  2:48,  1 user,  load average: 1.46, 2.52, 3.89
    admin02.localdomain:  23:45:04 up 909 days, 10:06,  3 users,  load average: 0.37, 0.45, 0.55
    admin01.localdomain:  23:45:04 up 909 days, 10:06,  3 users,  load average: 0.37, 0.45, 0.55
    admin03.localdomain:  23:45:04 up 2152 days,  4:35,  5 users,  load average: 0.20, 0.28, 0.22
    admin10.localdomain:  23:45:04 up 5 days,  4:00,  0 users,  load average: 0.02, 0.01, 0.00
    admin09.localdomain:  23:45:04 up 207 days,  2:48,  1 user,  load average: 6.87, 6.48, 7.04
    admin07.localdomain:  23:45:04 up 207 days,  2:33,  3 users,  load average: 6.55, 7.20, 7.35

    total hosts: 10 (10/0)

### Usage

    Usage:	./gpssh [option] command [argument ...]

    Options:
      -h, --help                Show help (this output)
      -V, --version             Show current version
      -f, --fanout int          Hosts to run in parallel (default 50)
      -n, --no-color            Disable colors
          --ok-color string     Color to use for stdout (default "#A8CC8C")
          --fail-color string   Color to use for stderr (default "#E88388")

    SSH Options:
      -r, --retries int   Number of ssh connection attempts (default 1)
      -s, --strict        Strict host key checking
      -v, --verbose       Verbose output (turns off quiet)
