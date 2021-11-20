package commands

import (
	"context"
	"fmt"
	"os"

	orbitdb "berty.tech/go-orbit-db"
	ifacedb "berty.tech/go-orbit-db/iface"
	"berty.tech/go-orbit-db/stores/operation"
	"github.com/LilRooks/five/internal/pkg/config"
)

const max int = 10

// Parse parses remaining args and logs accordingly
func Parse(store orbitdb.EventLogStore, cache orbitdb.DocumentStore, args []string, conf *config.ConfSet) error {
	amount := max
	cStream := ifacedb.StreamOptions{Amount: &amount} // Default up to 100 results, TODO Configurable
	switch cmd := args[0]; cmd {
	case "list":
		resultChan := make(chan operation.Operation)
		go store.Stream(context.Background(), resultChan, &cStream)
		for _ = range resultChan {
			fmt.Fprintf(os.Stdout, "%s\n", "fucko")
		}
	case "push":
		op, err := store.Add(context.Background(), []byte(args[1]))
		if err != nil {
			return err
		}
		strB, err := op.Marshal()
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "%s", string(strB))
	case "count":
		ops, err := store.List(context.Background(), &cStream)
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "%d", len(ops))
	default:
	}
	return nil
}

func search(store orbitdb.EventLogStore, conf *config.ConfSet) {

}
