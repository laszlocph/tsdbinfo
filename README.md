# tsdbinfo - Understand the series and labels you store in Prometheus

[![Build Status](https://cloud.drone.io/api/badges/laszlocph/tsdbinfo/status.svg)](https://cloud.drone.io/laszlocph/tsdbinfo)

`tsdbinfo` is a tool that looks into the data folder of Prometheus and gives you a basic understanding of what is stored there.

## Install

Linux
```
curl -LO https://github.com/laszlocph/tsdbinfo/releases/download/v0.1.5/tsdbinfo-v0.1.5-linux-amd64
mv tsdbinfo-v0.1.5-linux-amd64 tsdbinfo
chmod +x tsdbinfo
```

Mac
```
curl -LO https://github.com/laszlocph/tsdbinfo/releases/download/v0.1.5/tsdbinfo-v0.1.5-mac-amd64
mv tsdbinfo-v0.1.5-mac-amd64 tsdbinfo
chmod +x tsdbinfo
```

## Basic usage

#### List all the blocks

```bash
  ➜  tsdbinfo blocks --storage.tsdb.path=/my/prometheus/path/data
  ID                            FROM                         UNTIL                        STATS
  01CZWK46GK8BVHQCRNNS763NS3    2018-12-22T13:00:00+01:00    2018-12-29T07:00:00+01:00    {"numSamples":3167899784,"numSeries":3070548,"numChunks":29336192,"numBytes":4419004512}
  01D1EFWJ44G9WGN7AQ9398G2W2    2019-01-11T01:00:00+01:00    2019-01-11T19:00:00+01:00    {"numBytes":8634}
  01D1EFWJRQ35VYNT2M4YYEJV3R    2019-01-16T07:00:00+01:00    2019-01-17T01:00:00+01:00    {"numBytes":8634}
```

#### Identify the largest metrics

```bash
  ➜  tsdbinfo metrics --storage.tsdb.path=/my/prometheus/path/data --block=01CZWK46GK8BVHQCRNNS763NS3 --no-bar --top=3
  METRICSAMPLES        SERIES     LABELS
  solr_metrics_core_errors_total                              164,291,959    4,229      core: 99, handler: 32, collection: 16, replica: 9, instance: 5
  solr_metrics_core_time_seconds_total                        164,291,959    4,229      core: 99, handler: 32, collection: 16, replica: 9, instance: 5
  solr_metrics_core_timeouts_total                            164,291,959    4,229      core: 99, handler: 32, collection: 16, replica: 9, instance: 5
```

#### Investigate label explosion

```bash
  ➜  tsdbinfo metric --storage.tsdb.path=/my/prometheus/path/data --block=01CZWK46GK8BVHQCRNNS763NS3 --metric=http_server_requests_total
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
