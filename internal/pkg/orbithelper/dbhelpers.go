package orbithelper

import (
	"context"

	"github.com/LilRooks/five/internal/pkg/config"
	"github.com/LilRooks/five/internal/pkg/ipfshelper"
	"go.uber.org/zap"

	orbit "berty.tech/go-orbit-db"
	ifacedb "berty.tech/go-orbit-db/iface"
)

// Connect connects the EventLogStore with configuration
func Connect(logger *zap.Logger, conf *config.ConfSet, isLocal bool) (ifacedb.OrbitDB, error) {
	ipfs, err := ipfshelper.NewIPFS(conf, isLocal)
	if err != nil {
		return nil, err
	}
	dbConfig := &orbit.NewOrbitDBOptions{}

	db, err := orbit.NewOrbitDB(context.Background(), ipfs, dbConfig)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// OpenRemoteLog Opens the remote store (log)
func OpenRemoteLog(db ifacedb.OrbitDB, conf *config.ConfSet) (ifacedb.EventLogStore, error) {
	local := false
	cdb := ifacedb.CreateDBOptions{Create: &conf.MkStore, LocalOnly: &local}

	return db.Log(context.Background(), conf.Address, &cdb)
}

// OpenLocalDocs opens the local document store
func OpenLocalDocs(db ifacedb.OrbitDB, conf *config.ConfSet) (ifacedb.DocumentStore, error) {
	local := true
	cdb := ifacedb.CreateDBOptions{Create: &conf.MkLocal, LocalOnly: &local}

	return db.Docs(context.Background(), conf.AddrDoc, &cdb)
}
