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
	"sort"
	"text/tabwriter"

	"github.com/laszlocph/tsdbinfo/pkg/common"
	promTsdb "github.com/prometheus/tsdb"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var metric string

// metricCmd represents the metric command
var metricCmd = &cobra.Command{
	Use:   "metric",
	Short: "To dig deep on a single metric",
	Long: `
Digs deep on a given metric. It's best used to understand what labels you store and spot cardinality explosion that is bad for your Prometheus: https://prometheus.io/docs/practices/naming/#labels

Remember that every unique combination of key-value label pairs represents a new time series, which can dramatically increase the amount of data stored. Do not use labels to store dimensions with high cardinality (many different label values), such as user IDs, email addresses, or other unbounded sets of values.

Example usage:

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
	Run: func(cmd *cobra.Command, args []string) {
		if storagePath == "" {
			fmt.Fprintln(os.Stderr, "error: set --storage.tsdb.path")
			os.Exit(1)
		}

		if blockId == "" {
			fmt.Fprintln(os.Stderr, "error: set --block")
			os.Exit(2)
		}

		db, err := common.Open(storagePath)
		if err != nil {
			fmt.Printf("opening storage failed: %s", err)
		}

		var block *promTsdb.Block
		for _, b := range db.Blocks() {
			if b.Meta().ULID.String() == blockId {
				block = b
				break
			}
		}

		if block == nil {
			fmt.Fprintln(os.Stderr, "error: can't find block with id %s", blockId)
			os.Exit(2)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.TabIndent)
		p := message.NewPrinter(language.English)

		stat := numSamples(metric, db, block, false)

		fmt.Fprintf(w, "%s\t%v\n", "Metric", p.Sprint(stat.Metric))
		fmt.Fprintf(w, "%s\t%v\n", "Samples", p.Sprint(stat.Samples))
		fmt.Fprintf(w, "%s\t%v\n", "TimeSeries", p.Sprint(stat.Series))

		lstats := labelStats(metric, block)
		sort.Slice(lstats, func(i, j int) bool {
			return lstats[i].Occurrences > lstats[j].Occurrences
		})
		for _, s := range lstats {
			fmt.Fprintf(w, "Label\t%s\t%v\n", s.Label, p.Sprint(s.Occurrences))
		}

		labelStats := rawLabelStats(metric, block)
		for label, values := range labelStats {
			for v := range values {
				fmt.Fprintf(w, "LabelValue\t%s\t%v\n", label, v)
			}
		}

		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(metricCmd)
	metricCmd.PersistentFlags().StringVar(&blockId, "block", "", "verbose output")
	metricCmd.PersistentFlags().StringVar(&metric, "metric", "", "verbose output")
}
