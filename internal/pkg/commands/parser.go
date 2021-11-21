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
func Parse(store orbitdb.EventLogStore, cache orbitdb.EventLogStore, args []string, conf *config.ConfSet) error {
	amount := max
	cStream := ifacedb.StreamOptions{Amount: &amount} // Default up to 100 results, TODO Configurable
	switch cmd := args[0]; cmd {
	case "list":
		resultChan := make(chan operation.Operation)
		go store.Stream(context.Background(), resultChan, &cStream)
		for range resultChan {
			fmt.Fprintf(os.Stdout, "%s\n", "thing")
		}
	case "push":
		for _, word := range args[1:] {
			_, err := store.Add(context.Background(), []byte(word))
			if err != nil {
				fmt.Fprintf(os.Stdout, "errored: %s\n", err)
				return err
			}
		}

		//case "count":
		ops, err := store.List(context.Background(), &cStream)
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "%d\n", len(ops))
		for _, oop := range ops {
			fmt.Fprintf(os.Stdout, "%s\n", oop.GetValue())
		}
	default:
		fmt.Fprintf(os.Stdout, "should be done then!")
	}
	return nil
}

func search(store orbitdb.EventLogStore, conf *config.ConfSet) {

}
