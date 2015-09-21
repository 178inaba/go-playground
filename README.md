# go playground to gist

go playground to gist is web application that run code of go language, and post the code to gist.

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
# preparation: it is allowed to run the mongo
$ cd app/
$ cp setting/setting_example.toml setting/setting.toml
# edit setting/setting.toml in your client id and client secret
$ go run *.go # non build
# or
$ go build && ./app
```

#### option

* -d output debug log.

### use docker

```bash
# mongo
$ docker run --name mongo -d mongo

# sandbox
$ docker build --force-rm --no-cache -t tag/sandbox .
$ docker run -d -p 8081:8080 --name name-sandbox tag/sandbox

# playground
$ docker build --force-rm --no-cache -t tag/playground .
$ docker run -d -p 8082:8080 --link mongo:mongo --link name-sandbox:sandbox --name name-playground tag/playground
```
#### docker link host setting(setting.toml)

mongo host setting
```toml
[mongo]
host = "mongo"
```

and sandbox url setting
```toml
[sandbox]
url = "http://sandbox:8080/compile"
```

## licence

MIT
