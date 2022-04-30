package node

import (
	"os"
	"path/filepath"

	"github.com/p2p-games/wordle/libs/fslock"
)

// Init initializes the Node FileSystem Store for the given Node Type 'tp' in the directory under 'path' with
// default Config. Options are applied over default Config and persisted on disk.
func Init(path string, tp Type) error {
	cfg := DefaultConfig(tp)

	path, err := storePath(path)
	if err != nil {
		return err
	}
	log.Infof("Initializing %s Node Store over '%s'", tp, path)

	err = initRoot(path)
	if err != nil {
		return err
	}

	flock, err := fslock.Lock(lockPath(path))
	if err != nil {
		if err == fslock.ErrLocked {
			return ErrOpened
		}
		return err
	}
	defer flock.Unlock() // nolint: errcheck

	err = initDir(keysPath(path))
	if err != nil {
		return err
	}

	err = initDir(dataPath(path))
	if err != nil {
		return err
	}

	cfgPath := configPath(path)
	if !exists(cfgPath) {
		err = SaveConfig(cfgPath, cfg)
		if err != nil {
			return err
		}
		log.Infow("Saving config", "path", cfgPath)
	} else {
		log.Infow("Config already exists", "path", cfgPath)
	}

	log.Info("Node Store initialized")
	return nil
}

// IsInit checks whether FileSystem Store was setup under given 'path'.
// If any required file/subdirectory does not exist, then false is reported.
func IsInit(path string) bool {
	path, err := storePath(path)
	if err != nil {
		log.Errorw("parsing store path", "path", path, "err", err)
		return false
	}

	_, err = LoadConfig(configPath(path)) // load the Config and implicitly check for its existence
	if err != nil {
		log.Errorw("loading config", "path", path, "err", err)
		return false
	}

	if exists(keysPath(path)) &&
		exists(dataPath(path)) {
		return true
	}

	return false
}

const perms = 0755

// initRoot initializes(creates) directory if not created and check if it is writable
func initRoot(path string) error {
	err := initDir(path)
	if err != nil {
		return err
	}

	// check for writing permissions
	f, err := os.Create(filepath.Join(path, ".check"))
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return os.Remove(f.Name())
}

// initDir creates a dir if not exist
func initDir(path string) error {
	if exists(path) {
		return nil
	}
	return os.Mkdir(path, perms)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
