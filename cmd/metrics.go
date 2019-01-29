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
	"strings"
	"text/tabwriter"

	"github.com/gosuri/uiprogress"
	"github.com/laszlocph/tsdbinfo/pkg/common"
	promTsdb "github.com/prometheus/tsdb"
	"github.com/prometheus/tsdb/chunks"
	promTsdbLabels "github.com/prometheus/tsdb/labels"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var blockId string
var top int
var top_labels int
var no_bar bool

type metricStat struct {
	Metric  string
	Series  int
	Samples int
}

type labelStat struct {
	Label       string
	Occurrences int
}

func metrics(indexReader promTsdb.IndexReader) []string {
	values, _ := indexReader.LabelValues("__name__")

	var metrics []string
	for i := 0; i < values.Len(); i++ {
		ts, _ := values.At(i)
		for _, t := range ts {
			metrics = append(metrics, t)
		}
	}

	return metrics
}

func numSamples(metric string, db *promTsdb.DB, block *promTsdb.Block, debug bool) metricStat {
	var totalSamples int
	var totalTimeseries int
	meta := block.Meta()
	querier, _ := db.Querier(meta.MinTime, meta.MaxTime)
	seriesSet, err := querier.Select(promTsdbLabels.NewEqualMatcher("__name__", metric))
	if err != nil {
		fmt.Println(err)
	} else {
		for seriesSet.Next() {
			totalTimeseries++
			series := seriesSet.At()
			var numSamples int
			it := series.Iterator()
			for it.Next() {
				numSamples++
			}
			totalSamples = totalSamples + numSamples
			if debug {
				fmt.Printf("\t%v - %v samples\n", series.Labels(), numSamples)
			}
		}
	}

	return metricStat{metric, totalTimeseries, totalSamples}
}

func rawLabelStats(metric string, block *promTsdb.Block) map[string]map[string]bool {
	indexReader, _ := block.Index()
	p, _ := promTsdb.PostingsForMatchers(indexReader, promTsdbLabels.NewEqualMatcher("__name__", metric))

	var lset promTsdbLabels.Labels
	var chks []chunks.Meta

	labelStats := make(map[string]map[string]bool)
	for p.Next() {
		indexReader.Series(p.At(), &lset, &chks)

		for _, l := range lset {
			if labelStats[l.Name] == nil {
				labelStats[l.Name] = make(map[string]bool)
			}
			labelStats[l.Name][l.Value] = true
			// fmt.Printf("%v\n", l)
		}
	}

	return labelStats
}

func labelStats(metric string, block *promTsdb.Block) []labelStat {
	labelStats := rawLabelStats(metric, block)

	var stat []labelStat
	for label, values := range labelStats {
		stat = append(stat, labelStat{label, len(values)})
	}

	return stat
}

// metricsCmd represents the metrics command
var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "To identify the largest metrics in a given block",
	Long: `
Identifies the largest metrics in a given block. You can get block IDs with the "tsdb blocks" command.

NOTE: It does a sequencial scan on the given block so it may take a long time

Example usage:

  ➜  tsdbinfo metrics --storage.tsdb.path=/my/prometheus/path/data --block=01CZWK46GK8BVHQCRNNS763NS3 --no-bar --top=3
  METRICSAMPLES        SERIES     LABELS
  solr_metrics_core_errors_total                              164,291,959    4,229      core: 99, handler: 32, collection: 16, replica: 9, instance: 5
  solr_metrics_core_time_seconds_total                        164,291,959    4,229      core: 99, handler: 32, collection: 16, replica: 9, instance: 5
  solr_metrics_core_timeouts_total                            164,291,959    4,229      core: 99, handler: 32, collection: 16, replica: 9, instance: 5

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

		db, err := common.Open(storagePath, noPromLogs)
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

		indexReader, _ := block.Index()
		metrics := metrics(indexReader)

		uiprogress.Start()
		var bar *uiprogress.Bar
		if !no_bar {
			bar = uiprogress.AddBar(len(metrics))
			bar.AppendCompleted()
			bar.PrependElapsed()
		}

		var stat []metricStat
		for _, metric := range metrics {
			if !no_bar {
				bar.Incr()
			}
			stat = append(stat, numSamples(metric, db, block, false))
		}

		uiprogress.Stop()

		// metrics with most samples
		sort.Slice(stat, func(i, j int) bool {
			return stat[i].Samples > stat[j].Samples
		})

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.TabIndent)
		fmt.Fprintln(w, "METRIC\tSAMPLES\tSERIES\tLABELS")
		p := message.NewPrinter(language.English)

		if top > len(stat) {
			top = len(stat) - 1
		}
		for _, values := range stat[:top] {
			var statStrings []string
			lstats := labelStats(values.Metric, block)
			sort.Slice(lstats, func(i, j int) bool {
				return lstats[i].Occurrences > lstats[j].Occurrences
			})
			if top_labels > len(lstats) {
				top_labels = len(lstats) - 1
			}
			for _, s := range lstats[:top_labels] {
				statStrings = append(statStrings, p.Sprintf("%s: %d", s.Label, s.Occurrences))
			}

			fmt.Fprintf(w, "%s\t%v\t%v\t%s\n",
				values.Metric,
				p.Sprint(values.Samples),
				p.Sprint(values.Series),
				strings.Join(statStrings, ", "),
			)

		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(metricsCmd)
	metricsCmd.PersistentFlags().StringVar(&blockId, "block", "", "The ID of the TSDB block to inspect.")
	metricsCmd.PersistentFlags().IntVar(&top, "top", 100, "To control the length of the resultset. Default: 100")
	metricsCmd.PersistentFlags().IntVar(&top_labels, "top-labels", 5, "Number of labels to display. Default: 5")
	metricsCmd.PersistentFlags().BoolVar(&no_bar, "no-bar", false, "To hide the progressbar. In case you want to process the results.")
}
