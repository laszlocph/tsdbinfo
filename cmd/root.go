// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var storagePath string
var noPromLogs bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tsdbinfo",
	Short: "Understand the series and labels you store in Prometheus",
	Long: `tsdbinfo is a tool that looks into the data folder of Prometheus and gives you a basic understanding of what is stored there.

It follows the Prometheus timeseries database (tsdb) notions: blocks, series, labels and samples. Data is organized into blocks that hold all samples of a time period. Each metric has a number of series associated with it and the series has samples. Each label combination defines a new timeseries. You can find more on blocks here: https://fabxc.org/tsdb/

You can use tsdbinfo:

- to list all the blocks

	➜  tsdbinfo blocks --storage.tsdb.path=/my/prometheus/path/data
	ID                            FROM                         UNTIL                        STATS
	01CZWK46GK8BVHQCRNNS763NS3    2018-12-22T13:00:00+01:00    2018-12-29T07:00:00+01:00    {"numSamples":3167899784,"numSeries":3070548,"numChunks":29336192,"numBytes":4419004512}
	01D1EFWJ44G9WGN7AQ9398G2W2    2019-01-11T01:00:00+01:00    2019-01-11T19:00:00+01:00    {"numBytes":8634}
	01D1EFWJRQ35VYNT2M4YYEJV3R    2019-01-16T07:00:00+01:00    2019-01-17T01:00:00+01:00    {"numBytes":8634}

- to identify the largest metrics

	➜  tsdbinfo metrics --storage.tsdb.path=/my/prometheus/path/data --block=01CZWK46GK8BVHQCRNNS763NS3 --no-bar --top=3
	METRICSAMPLES        SERIES     LABELS
	solr_metrics_core_errors_total                              164,291,959    4,229      core: 99, handler: 32, collection: 16, replica: 9, instance: 5
	solr_metrics_core_time_seconds_total                        164,291,959    4,229      core: 99, handler: 32, collection: 16, replica: 9, instance: 5
	solr_metrics_core_timeouts_total                            164,291,959    4,229      core: 99, handler: 32, collection: 16, replica: 9, instance: 5

- to investigate label explosion

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

`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&storagePath, "storage.tsdb.path", "", "Prometheus TSDB path, same as in your Prometheus config")
	rootCmd.PersistentFlags().BoolVar(&noPromLogs, "no-prom-logs", false, "Hides Prometheus logs. Default false.")
}
