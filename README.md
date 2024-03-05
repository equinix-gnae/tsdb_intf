# Gotchas
* Number of datapoints returned are dependent on data ingested time and query start/end time
    * two TSDB by be off by 1 data point max
* Need to add time precession to query Data Structure (DS)?
    * Right now all the implementation are only returning time in milliseconds.
* add funtions like rate(prometheus)/non-negative-derative option.

# run tests

## Prometheuse
```bash
go test tests/tsdb_test.go -v -run TestPrometheus
```

## Mimir
```bash
go test tests/tsdb_test.go -v -run TestMimir
```

## InfluxDB
```bash
go test tests/tsdb_test.go -v -run TestInfluxDB
```

## All in one Test
```bash
go test tests/tsdb_test.go -v -run TestAllTSDBs
```

# Import guide
Due to https://stackoverflow.com/questions/32232655/go-get-results-in-terminal-prompts-disabled-error-for-github-private-repo

Configure go get to authenticate and fetch over https, all you need to do is to add the following line to $HOME/.netrc

```bash
machine github.com login USERNAME password TOKEN
```

Since its a private repo, its good not to cache it

```bash
go env -w GOPRIVATE=github.com/equinix-gnae/*
```