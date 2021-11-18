package app

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"

	"io"
	"io/fs"
	"io/ioutil"

	"github.com/LilRooks/ytdl-ipfs-archiver/internal/pkg/config"
	"github.com/LilRooks/ytdl-ipfs-archiver/internal/pkg/ipfs"
	"github.com/LilRooks/ytdl-ipfs-archiver/internal/pkg/table"
	"github.com/LilRooks/ytdl-ipfs-archiver/internal/pkg/ytdl"

	"github.com/ipfs/go-cid"
	"github.com/web3-storage/go-w3s-client"
)

var (
	address string
	cfgPath string
	keyPath string
)

const (
	errorNone = iota
	errorConfig
	errorOrbit
	errorIPFS
)

// Run is the actual code for the command
func Run(args []string, stdout io.Writer, stderr io.Writer) (int, error) {
	logger := log.New(stdout, "[base] ", log.Ltime|log.Lmsgprefix)

	// Parse the flags
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.StringVar(&cfgPath, "conf", "", "path to the configuration file to use")
	flags.StringVar(&sortsBy, "sort", "", "sorts results by ()") // TODO insert the options
	flags.BoolVar(&useAnon, "anon", false, "Whether to push as an anonymous user")
	flags.IntVar(&pullNth, "pull", 0, "Pulls the selected result")

	err := flags.Parse(args[1:])
	if err != nil {
		return errorConfig, err
	}
	nonFlagArgs := flags.Args()

	// Read configuration into configs variable
	configs, err = config.Read(cfgPath)
	if err != nil {
		return errorConfig, err
	}

	// Set defaults
	err = configs.Defaults(ConfSet{SortsBy: sortsBy, AnonIdentity: useAnon})
	if err != nil {
		return errorConfig, err
	}
	if val, ok := os.LookupEnv("TOKEN"); ok {
		configs.Token = val
	}

	if len(ytdlPath) == 0 {
		ytdlPath = configs.Binary
	}

	// Check binary is there
	ytdlPath, err = ytdl.ParsePath(ytdlPath)
	if err != nil {
		return errorYTDL, err
	}

	// Get the keys needed for the table
	var (
		filename string
		id       string
		format   string
		location string
	)
	_, errTableExist := os.Stat(tablPath)

	db, err := table.OpenDB(tablPath)
	if err != nil {
		return errorTable, err
	}
	defer db.Close()

	// Only initialized if the file did not originally exist
	if errors.Is(errTableExist, os.ErrNotExist) {
		err := table.InitializeTable(db)
		if err != nil {
			return errorTable, err
		}
	}

	id, format, err = ytdl.GetIdentifiers(logger, ytdlPath, ytdlOptions)
	if err != nil {
		return errorYTDL, err
	}

	location, err = table.Fetch(db, id, format)
	if err != nil {
		return errorTable, err
	}

	cid, _ := cid.Decode(location)

	c, err := w3s.NewClient(w3s.WithToken(configs.Token))
	if err != nil {
		return errorConfig, err
	}

	if len(location) == 0 {
		err := ytdl.Download(logger, ytdlPath, ytdlOptions)
		if err != nil {
			return errorYTDL, err
		}

		filename, err = ytdl.GetFilename(logger, ytdlPath, ytdlOptions)
		if err != nil {
			return errorYTDL, err
		}

		// only uploads first match, may have undefined behavior if file of same name exists
		// this is why file is stored with identifying information
		filenames, _ := filepath.Glob(filename + ".*")
		filename = filenames[0]

		location, err = ipfs.Store(c, filename)
		if err != nil {
			return errorIPFS, err
		}

		err = table.Store(db, id, format, location)
		if err != nil {
			return errorIPFS, err
		}
	} else {
		logger.SetPrefix("[w3s] ")
		logger.Printf("Getting %s\n", location)
		res, _ := c.Get(context.Background(), cid)

		// Download directory contents
		f, fsys, _ := res.Files()
		if d, ok := f.(fs.ReadDirFile); ok {
			ents, _ := d.ReadDir(0)
			for _, ent := range ents {
				filename = ent.Name()
				file, err := fsys.Open("/" + ent.Name())
				if err != nil {
					return errorIPFS, err
				}

				data, err := ioutil.ReadAll(file)
				if err != nil {
					return errorIPFS, err
				}

				err = ioutil.WriteFile(ent.Name(), data, 0755)
				if err != nil {
					return errorIPFS, err
				}
			}
		}
	}

	logger.SetPrefix("[base] ")
	logger.Printf("File is available locally at %s\n", filename)
	logger.Printf("File is also available at https://%s.ipfs.dweb.link\n", location)
	return errorNone, nil
}
