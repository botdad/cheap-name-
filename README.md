`go build`

then `./cheap-name -whatever` 

or `docker build . -t cheap-name`

then `docker run --rm cheap-name -whatever`

```
Usage of ./cheap-name:
  -bytes string
        hex bytes to match (default "00000000")
  -chars int
        number of random characters to use (default 5)
  -prefix string
        optional prefix before random string
  -selector string
        selector given
```

