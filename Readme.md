pbin
---


this is an experimental project for putting pastes on privatebin, see here for a public directory: https://privatebin.info/directory/


Install:
---
install the normal "go" way:
```
go get github.com/cbluth/pbin/cmd/pbin
```
or download a binary from the releases page:
```
https://github.com/cbluth/pbin/releases
```

Usage:
---

Upload Paste:
```
$ echo "anything" | pbin
https://privatebin.net/?5f9fc3956e8bc7bd#8NBafBFyqKWZrqPHiw4hC1JkL9Vx9mxEUGtXBT5wLNJF
```

Upload Base64 Paste:
```
$ cat cat-meme.gif | pbin -base64
https://privatebin.net/?c3dad23d043b0675#EEwJs9g3jSMC9gMHk5Gt5ptVDYpLXzCJMhP4Ufu3C3bf
```

Download Paste:
```
$ pbin https://privatebin.net/?908a9812a167d638#AKQaAp7bwC9t7gLBJkLXxJt1ZQQyW4bfjnBCzbn73c95
## prints to stdout
```

Download Base64 Paste:
```
$ pbin -base64 https://privatebin.net/?c3dad23d043b0675#EEwJs9g3jSMC9gMHk5Gt5ptVDYpLXzCJMhP4Ufu3C3bf > cat-meme.gif
```
