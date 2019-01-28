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
	Short: "A brief description of your command",
	Long: `
`,
	Run: func(cmd *cobra.Command, args []string) {

		if storagePath == "" {
			fmt.Fprintln(os.Stderr, "error: set --storage.tsdb.path")
			os.Exit(1)
		}

		db, err := common.Open(storagePath)
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
