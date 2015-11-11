package updater

// modified version of https://github.com/sanbornm/go-selfupdate/blob/master/selfupdate/selfupdate.go
// 463b28194bdc57bd431b638b80fcbb20eeb0790a

// Changes 9/10/15:
//     strip all space from time read in from the cktime file
//     removed all log statements
// Changes 9/11/15:
//     changed all usages of time to use validTime (this tells the program how long to wait before updating)
//     added a ForcedUpgrade method to rewrite the valid cktime and do a BackgroundRun

// Update protocol:
//
//   GET hk.heroku.com/hk/linux-amd64.json
//
//   200 ok
//   {
//       "Version": "2",
//       "Sha256": "..." // base64
//   }
//
// then
//
//   GET hkpatch.s3.amazonaws.com/hk/1/2/linux-amd64
//
//   200 ok
//   [bsdiff data]
//
// or
//
//   GET hkdist.s3.amazonaws.com/hk/2/linux-amd64.gz
//
//   200 ok
//   [gzipped executable data]
//
//

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/catalyzeio/catalyze/config"
	"github.com/kardianos/osext"
	"github.com/kr/binarydist"
	"gopkg.in/inconshreveable/go-update.v0"
)

const (
	upcktimePath = "cktime"
	plat         = runtime.GOOS + "-" + runtime.GOARCH
)

const validTime = 1 * 24 * time.Hour

// CLI auto updater
var AutoUpdater = &Updater{
	CurrentVersion: config.VERSION,
	APIURL:         "https://s3.amazonaws.com/cli-autoupdates/",
	BinURL:         "https://s3.amazonaws.com/cli-autoupdates/",
	DiffURL:        "https://s3.amazonaws.com/cli-autoupdates/",
	Dir:            ".catalyze_update",
	CmdName:        "catalyze",
}

// mismatch hash error
var ErrHashMismatch = errors.New("new file hash mismatch after patch")
var up = update.New()

// Updater is the configuration and runtime data for doing an update.
//
// Note that ApiURL, BinURL and DiffURL should have the same value if all files are available at the same location.
//
// Example:
//
//  updater := &selfupdate.Updater{
//  	CurrentVersion: version,
//  	ApiURL:         "http://updates.yourdomain.com/",
//  	BinURL:         "http://updates.yourdownmain.com/",
//  	DiffURL:        "http://updates.yourdomain.com/",
//  	Dir:            "update/",
//  	CmdName:        "myapp", // app name
//  }
//  if updater != nil {
//  	go updater.BackgroundRun()
//  }
type Updater struct {
	CurrentVersion string // Currently running version.
	APIURL         string // Base URL for API requests (json files).
	CmdName        string // Command name is appended to the ApiURL like http://apiurl/CmdName/. This represents one binary.
	BinURL         string // Base URL for full binary downloads.
	DiffURL        string // Base URL for diff downloads.
	Dir            string // Directory to store selfupdate state.
	Info           struct {
		Version string
		Sha256  []byte
	}
}

func (u *Updater) getExecRelativeDir(dir string) string {
	filename, _ := osext.Executable()
	path := filepath.Join(filepath.Dir(filename), dir)
	return path
}

// BackgroundRun starts the update check and apply cycle.
func (u *Updater) BackgroundRun() error {
	os.MkdirAll(u.getExecRelativeDir(u.Dir), 0755)
	if u.wantUpdate() {
		if err := up.CanUpdate(); err != nil {
			// fail
			return err
		}
		//self, err := osext.Executable()
		//if err != nil {
		// fail update, couldn't figure out path to self
		//return
		//}
		// TODO(bgentry): logger isn't on Windows. Replace w/ proper error reports.
		if err := u.update(); err != nil {
			return err
		}
	}
	return nil
}

// ForcedUpgrade writes a time in the past to the cktime file and then triggers
// the normal update process. This is useful when an update is required for
// the program to continue functioning normally.
func (u *Updater) ForcedUpgrade() error {
	path := u.getExecRelativeDir(filepath.Join(u.Dir, upcktimePath))
	writeTime(path, time.Now().Add(-1*validTime))
	return u.BackgroundRun()
}

func (u *Updater) wantUpdate() bool {
	path := u.getExecRelativeDir(filepath.Join(u.Dir, upcktimePath))
	if u.CurrentVersion == "dev" || readTime(path).After(time.Now()) {
		return false
	}
	//wait := 24*time.Hour + randDuration(24*time.Hour)
	return writeTime(path, time.Now().Add(validTime))
}

func (u *Updater) update() error {
	path, err := osext.Executable()
	if err != nil {
		return err
	}
	old, err := os.Open(path)
	if err != nil {
		return err
	}
	defer old.Close()

	err = u.FetchInfo()
	if err != nil {
		return err
	}
	if u.Info.Version == u.CurrentVersion {
		return nil
	}
	bin, err := u.fetchAndVerifyPatch(old)
	if err != nil {
		if err == ErrHashMismatch {
			log.Println("update: hash mismatch from patched binary - attempting full replacement...")
		} else {
			if u.DiffURL != "" {
				//log.Println("update: error patching binary,", err)
				// dont print this because a full binary replacement will be attempted.
				// if that fails, it will print out an error
				// if it succeeds, nothing needs to be printed
			}
		}

		bin, err = u.fetchAndVerifyFullBin()
		if err != nil {
			if err == ErrHashMismatch {
				log.Println("update: hash mismatch from full binary")
			} else {
				log.Println("update: error fetching full binary,", err)
			}
			log.Println("update: please upgrade your CLI manually")
			return err
		}
	}

	// close the old binary before installing because on windows
	// it can't be renamed if a handle to the file is still open
	old.Close()

	err, errRecover := up.FromStream(bytes.NewBuffer(bin))
	if errRecover != nil {
		return fmt.Errorf("update and recovery errors: %q %q", err, errRecover)
	}
	if err != nil {
		return err
	}
	log.Println("update: your CLI has been successfully upgraded!")
	return nil
}

// FetchInfo fetches and updates the info for latest CLI version available.
func (u *Updater) FetchInfo() error {
	r, err := fetch(u.APIURL + u.CmdName + "/" + plat + ".json")
	if err != nil {
		return err
	}
	defer r.Close()
	err = json.NewDecoder(r).Decode(&u.Info)
	if err != nil {
		return err
	}
	if len(u.Info.Sha256) != sha256.Size {
		return errors.New("bad cmd hash in info")
	}
	return nil
}

func (u *Updater) fetchAndVerifyPatch(old io.Reader) ([]byte, error) {
	bin, err := u.fetchAndApplyPatch(old)
	if err != nil {
		return nil, err
	}
	if !verifySha(bin, u.Info.Sha256) {
		return nil, ErrHashMismatch
	}
	return bin, nil
}

func (u *Updater) fetchAndApplyPatch(old io.Reader) ([]byte, error) {
	r, err := fetch(u.DiffURL + u.CmdName + "/" + u.CurrentVersion + "/" + u.Info.Version + "/" + plat)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var buf bytes.Buffer
	err = binarydist.Patch(old, &buf, r)
	return buf.Bytes(), err
}

func (u *Updater) fetchAndVerifyFullBin() ([]byte, error) {
	bin, err := u.fetchBin()
	if err != nil {
		return nil, err
	}
	verified := verifySha(bin, u.Info.Sha256)
	if !verified {
		return nil, ErrHashMismatch
	}
	return bin, nil
}

func (u *Updater) fetchBin() ([]byte, error) {
	r, err := fetch(u.BinURL + u.CmdName + "/" + u.Info.Version + "/" + plat + ".gz")
	if err != nil {
		return nil, err
	}
	defer r.Close()
	buf := new(bytes.Buffer)
	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(buf, gz); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// returns a random duration in [0,n).
// func randDuration(n time.Duration) time.Duration {
// 	return time.Duration(rand.Int63n(int64(n)))
// }

func fetch(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad http status from %s: %v", url, resp.Status)
	}
	return resp.Body, nil
}

func readTime(path string) time.Time {
	p, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		return time.Time{}
	}
	if err != nil {
		//return time.Now().Add(1000 * time.Hour)
		return time.Now().Add(validTime)
	}
	t, err := time.Parse(time.RFC3339, strings.TrimSpace(string(p)))
	if err != nil {
		//return time.Now().Add(1000 * time.Hour)
		return time.Now().Add(validTime)
	}
	return t
}

func verifySha(bin []byte, sha []byte) bool {
	h := sha256.New()
	h.Write(bin)
	return bytes.Equal(h.Sum(nil), sha)
}

func writeTime(path string, t time.Time) bool {
	return ioutil.WriteFile(path, []byte(t.Format(time.RFC3339)), 0644) == nil
}
