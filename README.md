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