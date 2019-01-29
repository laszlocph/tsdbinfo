# tsdbinfo - Understand the series and labels you store in Prometheus

[![Build Status](https://cloud.drone.io/api/badges/laszlocph/tsdbinfo/status.svg)](https://cloud.drone.io/laszlocph/tsdbinfo)

`tsdbinfo` is a tool that looks into the data folder of Prometheus and gives you a basic understanding of what is stored there.

## Install

Linux
```
curl -LO https://github.com/laszlocph/tsdbinfo/releases/download/v0.2.0/tsdbinfo-v0.2.0-linux-amd64
mv tsdbinfo-v0.2.0-linux-amd64 tsdbinfo
chmod +x tsdbinfo
```

Mac
```
curl -LO https://github.com/laszlocph/tsdbinfo/releases/download/v0.2.0/tsdbinfo-v0.2.0-mac-amd64
mv tsdbinfo-v0.2.0-mac-amd64 tsdbinfo
chmod +x tsdbinfo
```

## Basic usage

**First make a copy of the data in your Prometheus `--storage.tsdb.path` path.**

Running `tsdbinfo` on the same path - in parallel - with your production Prometheus may cause race conditions with unpredictable results. More on [Copying the Prometheus data folder](#Copying-the-Prometheus-data-folder)


#### List all the blocks

```bash
  ➜  tsdbinfo blocks --storage.tsdb.path.copy=/my/prometheus/path/data --no-prom-logs
  ID                            FROM                         UNTIL                        STATS
  01CZWK46GK8BVHQCRNNS763NS3    2018-12-22T13:00:00+01:00    2018-12-29T07:00:00+01:00    {"numSamples":3167899784,"numSeries":3070548,"numChunks":29336192,"numBytes":4419004512}
  01D1EFWJ44G9WGN7AQ9398G2W2    2019-01-11T01:00:00+01:00    2019-01-11T19:00:00+01:00    {"numBytes":8634}
  01D1EFWJRQ35VYNT2M4YYEJV3R    2019-01-16T07:00:00+01:00    2019-01-17T01:00:00+01:00    {"numBytes":8634}
```

#### Identify the largest metrics

```bash
  ➜  tsdbinfo metrics --storage.tsdb.path.copy=/my/prometheus/path/data --block=01CZWK46GK8BVHQCRNNS763NS3 --no-bar  --no-prom-logs --top=3
  METRICSAMPLES        SERIES     LABELS
  solr_metrics_core_errors_total                              164,291,959    4,229      core: 99, handler: 32, collection: 16, replica: 9, instance: 5
  solr_metrics_core_time_seconds_total                        164,291,959    4,229      core: 99, handler: 32, collection: 16, replica: 9, instance: 5
  solr_metrics_core_timeouts_total                            164,291,959    4,229      core: 99, handler: 32, collection: 16, replica: 9, instance: 5
```

#### Investigate label explosion

```bash
  ➜  tsdbinfo metric --storage.tsdb.path.copy=/my/prometheus/path/data --block=01CZWK46GK8BVHQCRNNS763NS3 --metric=http_server_requests_total --no-prom-logs
  Metric        http_server_requests_total
  Samples       5,365,975
  TimeSeries    582
  Label         path                    172
  Label         kubernetes_pod_name     14
  Label         instance                14
  Label         code                    10
  Label         pod_template_hash       7
  Label         method                  6
  Label         version                 5
  Label         app                     5
  Label         feature                 2
  Label         kubernetes_namespace    2
  Label         k8scluster              2
  Label         __name__                1
  Label         job                     1
  LabelValue    pod_template_hash       4120792344
  LabelValue    pod_template_hash       1980135293
  LabelValue    pod_template_hash       102934907
  LabelValue    pod_template_hash       1006272702
  LabelValue    pod_template_hash       3602012261
  LabelValue    pod_template_hash       3571852123
  LabelValue    pod_template_hash       3513057117
  LabelValue    code                    400
  LabelValue    code                    401
  LabelValue    code                    406
  LabelValue    code                    200
  LabelValue    code                    301
  LabelValue    code                    302
  LabelValue    code                    304
  LabelValue    code                    404
  LabelValue    code                    422
  LabelValue    code                    500
  LabelValue    method                  get
  LabelValue    method                  post
  LabelValue    method                  head
  LabelValue    method                  put
  LabelValue    method                  patch
...
```

## Uncover the sources of cardinality explosion in Prometheus

`tsdbinfo` is best used to understand what labels you store and spot cardinality explosion that is bad for your Prometheus: https://prometheus.io/docs/practices/naming/#labels

Remember that every unique combination of key-value label pairs represents a new time series, which can dramatically increase the amount of data stored. Do not use labels to store dimensions with high cardinality (many different label values), such as user IDs, email addresses, or other unbounded sets of values.

## Design

It follows the Prometheus timeseries database (tsdb) notions: blocks, series, labels and samples. Data is organized into blocks that hold all samples of a time period. Each metric has a number of series associated with it and the series has samples. Each label combination defines a new timeseries. You can find more on blocks here: https://fabxc.org/tsdb/

## Copying the Prometheus data folder

If you don't have the space to duplicate your Prometheus data files, just take a sample first of the folders as each subfolder of the `--storage.tsdb.path.copy` is self containing.


This is what I do:

```
$ ls -l /data/prometheus
total 88
drwxr-xr-x 3 prometheus prometheus 4096 Nov  5 09:02 01CVHHMPR9E7FANBKZZX7WVN2Y
drwxr-xr-x 3 prometheus prometheus 4096 Nov 12 03:02 01CW2XTTVMBR0PYDZ0KZ7C8HMG
drwxr-xr-x 3 prometheus prometheus 4096 Nov 18 21:02 01CWMA0JF9PZABN8JQ3HAR7K22
drwxr-xr-x 3 prometheus prometheus 4096 Nov 25 15:02 01CX5P6DY7RK8XFWY37PB4T1SY
drwxr-xr-x 3 prometheus prometheus 4096 Dec  2 09:01 01CXQ2C1YDJTRZATV8XYDG3W5N
drwxr-xr-x 3 prometheus prometheus 4096 Dec  9 03:01 01CY8EJ3T6670KC3Q9EQ8PV6DB
drwxr-xr-x 3 prometheus prometheus 4096 Dec 15 21:02 01CYSTR2E1M23DEAXM0C808BSR
drwxr-xr-x 3 prometheus prometheus 4096 Dec 22 15:02 01CZB6Y6922TDXYTR6GVRW3ADV
drwxr-xr-x 3 prometheus prometheus 4096 Dec 29 09:02 01CZWK46GK8BVHQCRNNS763NS3
drwxr-xr-x 3 prometheus prometheus 4096 Jan  5 03:02 01D0DZ9Q7JP7VHA1T00PA3ZJ71
drwxr-xr-x 3 prometheus prometheus 4096 Jan 11 21:01 01D0ZBF8C2DD5QC8WD8V765SEX
drwxr-xr-x 3 prometheus prometheus 4096 Jan 18 15:01 01D1GQNBDP2NXH6BVN1E022VPS
drwxr-xr-x 3 prometheus prometheus 4096 Jan 25 09:02 01D223V9N0YPN673FZPTXPDME8
drwxr-xr-x 3 prometheus prometheus 4096 Jan 27 15:00 01D27X71GWJE0SAB2P6H98FE4W
drwxr-xr-x 3 prometheus prometheus 4096 Jan 28 09:00 01D29V089NE4VF321Q2TA8FF3T
drwxr-xr-x 3 prometheus prometheus 4096 Jan 29 03:00 01D2BRSQCV29G3SFTB5GAXMXV1
drwxr-xr-x 3 prometheus prometheus 4096 Jan 29 09:00 01D2CDCPWR58D42P8MVGKT0B1S
drwxr-xr-x 3 prometheus prometheus 4096 Jan 29 15:00 01D2D1ZJ7F2ZZP2H2PREEFFEFX
drwxr-xr-x 3 prometheus prometheus 4096 Jan 29 15:00 01D2D2006A2HK3V8A4SSZR92XS
drwxr-xr-x 3 prometheus prometheus 4096 Jan 29 17:00 01D2D8V9E2YVRW1BE884YKKW27
-rw-r--r-- 1 prometheus prometheus    0 Aug  9 12:37 lock
drwxr-xr-x 2 prometheus prometheus 4096 Jan 29 17:43 wal
```

Then copying recent larger chunks:

```
$ mkdir /data/prom-analysis
$ cd /data/prom-analysis
$ cp -r ../prometheus/01D*/ .
```

Then running `tsdbinfo`:

```
$ tsdbinfo blocks --storage.tsdb.path.copy=/data/prom-analysis --no-prom-logs
```