# pbin - privatebin cli tool
---

this is an experimental project for putting pastes on privatebin, see here for a public directory: https://privatebin.info/directory/


# Install:
install the normal "go" way:
```
go get github.com/cbluth/pbin/cmd/pbin
```
or download a binary from the releases page:
```
https://github.com/cbluth/pbin/releases
```

# Basic Usage:

Upload Paste:
```
$ echo "anything" | pbin
https://privatebin.net/?5f9fc3956e8bc7bd#8NBafBFyqKWZrqPHiw4hC1JkL9Vx9mxEUGtXBT5wLNJF
```

Download Paste:
```
$ URL="https://privatebin.net/?908a9812a167d638#AKQaAp7bwC9t7gLBJkLXxJt1ZQQyW4bfjnBCzbn73c95"
$ pbin $URL
## prints content to stdout
```

## Advanced Usage

You can set additional options if some of these arguments, only when creating a paste:
- -burn
- -open
- -base64
- -password

Upload Base64 Paste:
```
$ cat cat-meme.gif | pbin -base64
```

Download Base64 Paste:
```
$ pbin $URL -base64 > cat-meme.gif
```

Download Base64 Paste to file:
```
$ pbin $URL -base64 -o cat-meme.gif
```

Upload Paste with Burn After Read Once:
```
$ echo "anything" | pbin -burn
```

Download Paste to filepath:
```
$ pbin $URL -output cat-meme.gif
```

Upload Paste with discussion forum enabled:
```
$ echo "anything" | pbin -open
```

Upload Paste with password protection:
```
$ echo "anything" | pbin -password mySecretPassw0rd
```

## Expiry Options

You can set the expiry with one of these arguments, only when creating a paste:
- -hour
- -day
- -week
- -month
- -year
- -never

examples:
```
$ echo "anything" | pbin -hour # <- will expire after 1 hour
$ echo "anything" | pbin -week # <- will expire after 1 week
$ echo "anything" | pbin -month # <- will expire after 1 month
$ cat cat-meme.gif | pbin -base64 -never
```


# TODO:

A list of things to do:
- add reply/comment
- add upload file
- add shorten url

