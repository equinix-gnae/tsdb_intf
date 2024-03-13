package ts_test

import (
	"context"
	"testing"

	"github.com/equinix-gnae/tsdb_intf/pkg/ts"
	"github.com/kr/pretty"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var v = viper.GetViper()

func TestConfig(t *testing.T) {
	config := ts.ConfigFromViper(v)
	assert.Equal(t, "mimir", config.Name, "make sure to run the test with TSDB_NAME=mimir env set")
}

func TestTSDBObjUsingConfig(t *testing.T) {
	config := ts.ConfigFromViper(v)
	tsdb, err := ts.NewTSDBClient(&config)

	if err != nil {
		t.Errorf("got an error: %v", err)
	}

	if result, err := tsdb.Query(context.Background(), baseQuery); err != nil {
		t.Errorf("got an error: %v", err)
	} else {
		pretty.Print(result)
	}
}

func init() {
	ts.MustViperFlags(v, &pflag.FlagSet{})
}
