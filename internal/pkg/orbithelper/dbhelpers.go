package orbithelper

import (
	"context"

	"go.uber.org/zap"

	orbit "berty.tech/go-orbit-db"
	"berty.tech/go-orbit-db/iface"

	"github.com/LilRooks/five/internal/pkg/config"
	"github.com/LilRooks/five/internal/pkg/ipfshelper"
)

// Connect connects the EventLogStore with configuration
func Connect(logger *zap.Logger, conf *config.ConfSet) (iface.EventLogStore, error) {
	ipfs := ipfshelper.NewIPFS(conf.RepoDir)
	dbConfig := &orbit.NewOrbitDBOptions{
		Logger: logger}

	db, err := orbit.NewOrbitDB(context.Background(), ipfs, dbConfig)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
