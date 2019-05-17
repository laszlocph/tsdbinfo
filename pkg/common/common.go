package common

import (
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/tsdb"
)

func Open(storagePath string, noPromLogs bool) (*tsdb.DB, error) {
	var w io.Writer
	if noPromLogs {
		w = log.NewSyncWriter(ioutil.Discard)
	} else {
		w = log.NewSyncWriter(os.Stderr)
	}
	logger := log.NewLogfmtLogger(w)

	db, err := tsdb.Open(
		storagePath,
		log.With(logger, "component", "tsdb"),
		prometheus.DefaultRegisterer,
		&tsdb.Options{
			BlockRanges: tsdb.ExponentialBlockRanges(int64(time.Hour*2/time.Millisecond), 10, 3),
		},
	)

	return db, err
}
