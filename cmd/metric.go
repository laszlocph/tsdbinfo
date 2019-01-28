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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
