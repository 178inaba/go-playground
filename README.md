# go playground to gist

go playground to gist is web application that run code of go language.
then code made can be saved to gist.

## description

inspired [golang/playground](https://github.com/golang/playground).  
code that you made can be posted to [github gist](https://gist.github.com).

## run

### sandbox

```bash
$ cd sandbox/
$ docker build --force-rm --no-cache -t tag/sandbox .
$ docker run -d -p 8080:8080 --name name-sandbox tag/sandbox
```

### app

```bash
$ cd app/
$ cp setting/setting_example.toml setting/setting.toml
# edit setting/setting.toml in your client id and client secret
$ go run *.go # non build
# or
$ go build && ./app
```

#### option

* -d output debug log.

## licence

MIT
