iTunes daemon - stream from iTunes to a Linux server (RAOP)

## install

```
# apt-get install libpulse-dev
$ go get -u github.com/alicebob/ituned
$ $GOPATH/bin/ituned
```

The daemon should show up on your ipod/itunes.

This thing has a bonjour server build in, if you have something else running on
:5353 it'll fail to start. Most likely that something else will be an avahi
daemon.

## status

It compiles, ship it.

## credits

Most of the code is copied from http://github.com/joelgibson/go-airplay
