package ts

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.infratographer.com/x/viperx"
)

// Config is used to configure a new tsdb server
type Config struct {
	// name for tsdb, name is used to figure which client to create
	Name string

	// url/endpoint for tsdb
	URL string

	// Name of the database/database like construct e.g. bucket for influxDB
	Database string

	// org is a common construct used in multi tenant DBs e.g. X-ORG-ID for mimir and ORG-ID for influxDB
	Org string

	// ID can be use for secret
	ID string

	// Secret can be use to store auth related data e.g. token for influxDB
	Secret string
}

// MustViperFlags returns the cobra flags and wires them up with viper to prevent code duplication
func MustViperFlags(v *viper.Viper, flags *pflag.FlagSet) {
	flags.String("tsdb-name", "prometheus", "name for the tsdb")
	viperx.MustBindFlag(v, "tsdb.name", flags.Lookup("tsdb-name"))

	flags.String("tsdb-url", "http://prometheus:9090", "address to connect on")
	viperx.MustBindFlag(v, "tsdb.url", flags.Lookup("tsdb-url"))

	flags.String("tsdb-database", "", "database name")
	viperx.MustBindFlag(v, "tsdb.database", flags.Lookup("tsdb-database"))

	flags.String("tsdb-org", "", "org or tenant name")
	viperx.MustBindFlag(v, "tsdb.org", flags.Lookup("tsdb-org"))

	flags.String("tsdb-id", "", "id/username")
	viperx.MustBindFlag(v, "tsdb.id", flags.Lookup("tsdb-id"))

	flags.String("tsdb-secret", "", "secret/token for authentication")
	viperx.MustBindFlag(v, "tsdb.secret", flags.Lookup("tsdb-secret"))
}

// ConfigFromViper builds a new Config from viper.
func ConfigFromViper(v *viper.Viper) Config {
	return Config{
		Name:     v.GetString("tsdb.name"),
		URL:      v.GetString("tsdb.url"),
		Database: v.GetString("tsdb.database"),
		Org:      v.GetString("tsdb.org"),
		ID:       v.GetString("tsdb.id"),
		Secret:   v.GetString("tsdb.secret"),
	}
}
