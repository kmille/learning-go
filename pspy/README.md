# pspy

This is my own implementation of [pspy](https://github.com/DominicBreuker/pspy) to learn Go. It prints every command which is executed.

TODO: describe what it does/how it works

# TO READ

- defer over functions
- global variables - good or evil


- https://go.googlesource.com/proposal/+/master/design/go2draft-generics-overview.md

# TODO

- makefile
- make a nice README.md
    - let Github Actions run the tests
    - provide a link to the working binary in the README
    - make an asciinema capture
- rename argv[0]
- We are too fast. The cmdline of a pid changes ....

```
2021/02/16 13:30:33 [0] cmdline of 141026 is zsh                       
2021/02/16 13:30:33 [1] cmdline of 141026 is zsh                                                                       
2021/02/16 13:30:33 [2] cmdline of 141026 is zsh                       
2021/02/16 13:30:33 [3] cmdline of 141026 is ls                       
2021/02/16 13:30:33 [4] cmdline of 141026 is ls -vlah --color=auto 
```

- we are vuln to races but I'm not sure if it's a problem

- Tests
    - compare results with pspy
    - are we now get every command? What about the hidden commands  (tab tab, git executed by the shell)

