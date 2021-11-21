package orbithelper

import (
	"context"
	"fmt"

	"github.com/LilRooks/five/internal/pkg/config"
	"github.com/LilRooks/five/internal/pkg/ipfshelper"
	"go.uber.org/zap"

	orbit "berty.tech/go-orbit-db"
	"berty.tech/go-orbit-db/accesscontroller"
	"berty.tech/go-orbit-db/cache"
	"berty.tech/go-orbit-db/cache/cacheleveldown"
	ifacedb "berty.tech/go-orbit-db/iface"
)

const greplicate bool = true

// Connect connects the EventLogStore with configuration
func Connect(logger *zap.Logger, conf *config.ConfSet, isLocal bool) (ifacedb.OrbitDB, error) {
	ipfs, err := ipfshelper.NewIPFS(conf, isLocal)
	if err != nil {
		return nil, err
	}
	cacheSettings := cache.Options{Logger: logger}
	dbConfig := &orbit.NewOrbitDBOptions{Cache: cacheleveldown.New(&cacheSettings)}

	db, err := orbit.NewOrbitDB(context.Background(), ipfs, dbConfig)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// OpenRemoteLog Opens the remote store (log)
func OpenRemoteLog(db ifacedb.OrbitDB, conf *config.ConfSet) (ifacedb.EventLogStore, error) {
	local := false
	replicate := greplicate
	ac := accesscontroller.NewEmptyManifestParams()
	ac.SetAccess("write", []string{"*"})

	var cdb ifacedb.CreateDBOptions
	if conf.MkStore {
		cdb = ifacedb.CreateDBOptions{Create: &conf.MkStore, LocalOnly: &local, Replicate: &replicate}
		cdb.AccessController = ac
	}

	fmt.Printf("attept open store log\n")
	st, err := db.Log(context.Background(), conf.Address, &cdb)
	fmt.Printf("A open store log\n")
	return st, err
}

// OpenLocalDocs opens the local document store
func OpenLocalDocs(db ifacedb.OrbitDB, conf *config.ConfSet) (ifacedb.DocumentStore, error) {
	local := true
	replicate := false
	var cdb ifacedb.CreateDBOptions
	if conf.MkLocal {
		cdb = ifacedb.CreateDBOptions{Create: &conf.MkLocal, LocalOnly: &local, Replicate: &replicate}
	}

	return db.Docs(context.Background(), conf.AddrDoc, &cdb)
}
