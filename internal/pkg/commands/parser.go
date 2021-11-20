package commands

import (
	"context"

	orbitdb "berty.tech/go-orbit-db"
	ifacedb "berty.tech/go-orbit-db/iface"
	"berty.tech/go-orbit-db/stores/operation"
	"github.com/LilRooks/five/internal/pkg/config"
	"go.uber.org/zap"
)

const max int = 100

// Parse parses remaining args and logs accordingly
func Parse(store orbitdb.EventLogStore, cache orbitdb.DocumentStore, args []string, conf *config.ConfSet) error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	sugar := logger.Sugar()
	defer sugar.Sync()
	switch cmd := args[0]; cmd {
	case "list":
		resultChan := make(chan operation.Operation)
		amount := max
		cStream := ifacedb.StreamOptions{Amount: &amount} // Default up to 100 results, TODO Configurable
		go store.Stream(context.Background(), resultChan, &cStream)
		for op := range resultChan {
			sugar.Info("content", string(op.GetValue()))
		}
	}
	return nil
}

func search(store orbitdb.EventLogStore, conf *config.ConfSet) {

}
