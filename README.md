# Mailosaur Go Client Library

[![Build Status](https://travis-ci.com/jslang/mailosaur-go.svg?branch=master)](https://travis-ci.com/jslang/mailosaur-go)

[Mailosaur](https://mailosaur.com/) allows you to automate tests involving email. Allowing you to perform end-to-end automated and functional email testing.


## Installation

```
go get -u github.com/jslang/mailosaur-go
```

## Usage

```
import mailosaur

c := mailosaur.NewClient("<yourapikey>", "<yourserverid>")
msgs, err := c.ListMessages()
...
msg, err := c.GetMessage(msgs[0].Id)
...
err = c.DeleteMessage(msg.Id)
...
err = c.DeleteMessages()
...
```

## Tests

Unit tests

```
go test ./mailosaur
```

Integration tests

```
MAILOSAUR_API_KEY=<apikey> MAILOSAUR_SERVER_ID=<serverid> go test ./test/integration_test.go
```
