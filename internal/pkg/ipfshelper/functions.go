package ipfshelper

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/LilRooks/five/internal/pkg/config"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/plugin/loader"

	fsconfig "github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	icore "github.com/ipfs/interface-go-ipfs-core"
)

// NewIPFS creates a new ipfs coreapi instance and returns it
func NewIPFS(conf *config.ConfSet, isLocal bool) (icore.CoreAPI, error) {
	var createRepo bool
	var repoPath string
	if isLocal {
		repoPath = conf.RepoLDB
		createRepo = conf.MkLocal
	} else {
		repoPath = conf.RepoDir
		createRepo = conf.MkStore
	}
	return createNode(context.Background(), repoPath, isLocal, createRepo)
}

func createNode(ctx context.Context, repoPath string, isLocal bool, createRepo bool) (icore.CoreAPI, error) {
	if createRepo {
		initRepo(repoPath, isLocal)
	}

	// Open the repo
	repo, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, err
	}

	// Construct the node

	nodeOptions := &core.BuildCfg{
		Online:    !isLocal,
		Permanent: true,             // It is temporary way to signify that node is permanent
		Routing:   libp2p.DHTOption, // This option sets the node to be a full DHT node (both fetching and storing DHT Records)
		// Routing: libp2p.DHTClientOption, // This option sets the node to be a client DHT node (only fetching records)
		Repo: repo,
		ExtraOpts: map[string]bool{
			"pubsub":  true,
			"namesys": true,
			"ipnsps":  true,
		},
	}

	node, err := core.NewNode(ctx, nodeOptions)
	if err != nil {
		return nil, err
	}

	// Attach the Core API to the constructed node
	return coreapi.NewCoreAPI(node)
}

func initRepo(repoPath string, isLocal bool) error {
	cfg, err := fsconfig.Init(ioutil.Discard, 2048)
	if err != nil {
		return err
	}

	return fsrepo.Init(repoPath, cfg)
}

// SetupPlugins sets up plugins once
func SetupPlugins() error {
	// Load any external plugins if available on externalPluginsPath
	plugins, err := loader.NewPluginLoader("plugins")
	if err != nil {
		return fmt.Errorf("error loading plugins: %s", err)
	}

	// Load preloaded and external plugins
	if err := plugins.Initialize(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	if err := plugins.Inject(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	return nil
}
