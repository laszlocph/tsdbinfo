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
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/laszlocph/tsdbinfo/pkg/common"
	"github.com/spf13/cobra"
)

// blocksCmd represents the blocks command
var blocksCmd = &cobra.Command{
	Use:   "blocks",
	Short: "Lists all the blocks under the Prometheus TSDB path. With metadata",
	Long: `
Lists all the blocks under the Prometheus TSDB path. With metadata

Example usage:

  ➜  tsdbinfo blocks --storage.tsdb.path.copy=/my/prometheus/path/data
  ID                            FROM                         UNTIL                        STATS
  01CZWK46GK8BVHQCRNNS763NS3    2018-12-22T13:00:00+01:00    2018-12-29T07:00:00+01:00    {"numSamples":3167899784,"numSeries":3070548,"numChunks":29336192,"numBytes":4419004512}
  01D1EFWJ44G9WGN7AQ9398G2W2    2019-01-11T01:00:00+01:00    2019-01-11T19:00:00+01:00    {"numBytes":8634}
  01D1EFWJRQ35VYNT2M4YYEJV3R    2019-01-16T07:00:00+01:00    2019-01-17T01:00:00+01:00    {"numBytes":8634}

`,
	Run: func(cmd *cobra.Command, args []string) {

		if storagePath == "" {
			fmt.Fprintln(os.Stderr, "error: set --storage.tsdb.path.copy")
			os.Exit(1)
		}

		db, err := common.Open(storagePath, noPromLogs)
		if err != nil {
			fmt.Printf("opening storage failed: %s", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.TabIndent)
		fmt.Fprintln(w, "ID\tFROM\tUNTIL\tSTATS")

		for _, block := range db.Blocks() {
			meta := block.Meta()
			stats, _ := json.Marshal(meta.Stats)
			fmt.Fprintf(w, "%s\t%v\t%v\t%s\n",
				meta.ULID,
				time.Unix(meta.MinTime/1000, 0).Format(time.RFC3339),
				time.Unix(meta.MaxTime/1000, 0).Format(time.RFC3339),
				string(stats),
			)
		}
		w.Flush()

	},
}

func init() {
	rootCmd.AddCommand(blocksCmd)
}
