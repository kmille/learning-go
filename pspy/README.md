## pspy (process spy) [![Go tests for pspy](https://github.com/kmille/learning-go/actions/workflows/pspy_tests.yml/badge.svg)](https://github.com/kmille/learning-go/actions/workflows/pspy_tests.yml)

This is my own implementation of [pspy](https://github.com/DominicBreuker/pspy) to learn Go. It prints every executed command. It's nice to find out what a server is doing. In CTFs, it's nice to catch credentials when they are passed via the command line.
## How to get it?

[Download link](https://github.com/kmille/learning-go/releases/download/v0.1/pspy) or via command line:

```bash
kmille@linbox:pspy wget https://github.com/kmille/learning-go/releases/download/v0.1/pspy && chmod +x pspy
...
2021-03-02 13:14:30 (23.6 MB/s) - ‘pspy’ saved [2464899/2464899]
kmille@linbox:pspy ./pspy -h
Usage of ./pspy:
  -cmd string
        filter CMD
  -debug
        debug print for every event
  -uid int
        filter UID (default -1)
  -w string
        output file (default "-")
```

## How does it work?
pspy uses Linux' [inotifywatch](https://linux.die.net/man/1/inotifywatch) capabilities to get notified if a file is opened. 

1. parse the $PATH variable
2. watch IN_OPEN events on every directory in $PATH. At this point we know if a binary is executed
3. parse /proc und parse process information. Print the data if we haven't already
4. go to step 2 and wait for more events

## Demo
[![asciicast](https://asciinema.org/a/395925.svg)](https://asciinema.org/a/395925)

## Improvements
What's interesting is that sometimes we are too fast. Look hat the output of `/proc/141026/cmdline`:

```
2021/02/16 13:30:33 [0] cmdline of 141026 is zsh                       
2021/02/16 13:30:33 [1] cmdline of 141026 is zsh                                                                       
2021/02/16 13:30:33 [2] cmdline of 141026 is zsh                       
2021/02/16 13:30:33 [3] cmdline of 141026 is ls                       
2021/02/16 13:30:33 [4] cmdline of 141026 is ls -vlah --color=auto 
```

Even it's the same pid, the content of cmdline changes (from `zsh` which is my shell to `ls` which is an alias to then `ls -vlah --color=auto?`).
We are not thread safe at the moment. But I think in the worst case we print a process more than once (which is ok for now).
