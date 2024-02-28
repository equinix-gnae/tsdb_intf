# Gotchas
* Influx is returning one additional datapoint
* Influx timestamp is seconds where as prometheus timestamp is in milliseconds
    * Need to add time precession to query Data Structure (DS)?
        * since prometheus doesn't allow to set precision, we can't add it to query DS
    * This require all the DBs to store time with same precision?
        * influx allows us to wirte data in different time precision
* add funtions like rate(prometheus)/non-negative-derative option