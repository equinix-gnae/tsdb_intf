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

const (
	url string = "http://mgmtsrv1.sv11.edn.equinix.com:32090"
)

func TestConfig(t *testing.T) {
	config := ts.ConfigFromViper(v)
	assert.Equal(t, url, config.URL)
}

func TestTSDBObjUsingConfig(t *testing.T) {
	config := ts.ConfigFromViper(v)
	tsdb, err := ts.NewTSDBClient(&config)

	if err != nil {
		t.Errorf("got an error: %v", err)
	}

	if result, err := tsdb.Query(context.Background(), BaseQuery); err != nil {
		t.Errorf("got an error: %v", err)
	} else {
		pretty.Print(result)
	}
}

func init() {
	v.SetDefault("tsdb.url", url)
	v.SetDefault("tsdb.id", "")
	v.SetDefault("tsdb.secret", "")
	ts.MustViperFlags(v, &pflag.FlagSet{})
}
