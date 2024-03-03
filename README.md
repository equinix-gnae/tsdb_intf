# Gotchas
* Number of datapoints returned are dependent on data ingested time and query start/end time
    * two TSDB by be off by 1 data point max
* Need to add time precession to query Data Structure (DS)?
    * Right now all the implementation are only returning time in milliseconds.
* add funtions like rate(prometheus)/non-negative-derative option.

# TODO
* add tests