package app

import (
	"context"
	"flag"
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/LilRooks/five/internal/pkg/commands"
	"github.com/LilRooks/five/internal/pkg/config"
	"github.com/LilRooks/five/internal/pkg/ipfshelper"
	"github.com/LilRooks/five/internal/pkg/orbithelper"
)

var (
	address string
	cfgPath string
	mkStore bool
	mkLocal bool
)

const (
	errorNone = iota
	errorGeneric
)

// Run is the actual code for the command
func Run(args []string) (int, error) {
	// Parse the flags
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.StringVar(&cfgPath, "conf", "", "path to the configuration file to use")
	flags.StringVar(&cfgPath, "addr", "", "address to the remote store to use")
	flags.BoolVar(&mkLocal, "mkdocs", false, "make the document store")
	flags.BoolVar(&mkStore, "mklogs", false, "make the remote log store")

	err := flags.Parse(args[1:])
	if err != nil {

	}
	nonFlagArgs := flags.Args()

	// Read configuration into configs variable
	configs, err := config.Read(cfgPath)
	if err != nil {
		return errorGeneric, err
	}
	// Set defaults
	flagConfigs := config.ConfSet{
		MkStore: mkStore,
		MkLocal: mkLocal,
	}
	err = configs.Defaults(&flagConfigs)
	if err != nil {
		return errorGeneric, err
	}
	if val, ok := os.LookupEnv("TOKEN"); ok {
		configs.IDToken = val
	}

	// Create logger
	rLogger, err := zap.NewDevelopment()
	if err != nil {
		return 1, err
	}
	defer rLogger.Sync()
	lLogger, err := zap.NewDevelopment()
	if err != nil {
		return 1, err
	}
	defer rLogger.Sync()

	if err := ipfshelper.SetupPlugins(); err != nil {
		return errorGeneric, err
	}

	// Connect database
	rdb, err := orbithelper.Connect(rLogger, configs, false)
	if err != nil {
		return errorGeneric, err
	}
	defer rdb.Close()

	ldb, err := orbithelper.Connect(lLogger, configs, true)
	if err != nil {
		return errorGeneric, err
	}
	fmt.Printf("Opened store at db\n")

	// Open store on configs
	rStore, err := orbithelper.OpenRemoteLog(rdb, configs)
	if err != nil {
		return errorGeneric, err
	}
	defer rStore.Close()
	fmt.Printf("Opened store at %s\n", rStore.Address().String())
	if err := rStore.Load(context.Background(), 100); err != nil {
		return errorGeneric, err
	}

	// Open local store on configs
	lStore, err := orbithelper.OpenLocalDocs(ldb, configs)
	if err != nil {
		return errorGeneric, err
	}
	defer lStore.Close()

	commands.Parse(rStore, lStore, nonFlagArgs, configs)

	return errorNone, nil
}
