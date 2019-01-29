package common

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/storage/tsdb"
	promTsdb "github.com/prometheus/tsdb"
)

func options() *tsdb.Options {
	options := tsdb.Options{}

	duration := new(model.Duration)
	duration.Set("2h")
	options.MinBlockDuration = *duration

	maxBlockDuration := new(model.Duration)
	maxBlockDuration.Set("<default>")
	options.MaxBlockDuration = *maxBlockDuration

	return &options
}

func Open(storagePath string, noPromLogs bool) (*promTsdb.DB, error) {
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
		options(),
	)

	return db, err
}
