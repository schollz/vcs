package vcs

import (
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/schollz/reldel"
)

// VCSystem is the structure for the version-controlled system
type VCSystem struct {
	Owners        []string          `json:"owners"`
	CurrentBranch string            `json:"current_branch"`
	Files         map[string]VCFile `json:"files"`
}

// VCFile is the structure for a version-controlled file
type VCFile struct {
	StartingBlock string           `json:"starting_block"`
	BlockMap      map[string]Block `json:"blocks"` // organized by hash
}

// Block contains the hash data and patch
type Block struct {
	Branch       string       `json:"branch"`
	PreviousHash string       `json:"previous"`
	Patch        reldel.Patch `json:"patch"`
}

// Init returns a new version controlled file or loads
// one from a file
func Init() (*VCSystem, error) {
	vc := new(VCSystem)
	vc.Owners = []string{identity.Public}
	vc.CurrentBranch = "master"
	vc.Files = make(map[string]VCFile)
	return vc, nil
}

// Commit will write the current commit to a file
func (vc *VCSystem) Commit(text string) (err error) {
	// TODO add file locking
	// need to determine current hash by re-capitulating file

	n := Block{
		Branch:       vc.CurrentBranch,
		PreviousHash: currentHash,
		Patch:        reldel.GetPatch(vc.CurrentText, text),
	}
	h := sha256.New()
	h.Write([]byte("vcf")) // salt
	h.Write([]byte(text))
	h.Write([]byte(n.Patch.Time.String()))
	blockHash = fmt.Sprintf("%x", h.Sum(nil))
	vc.StartingBlock = blockHash
	vc.BlockMap[blockHash] = n

	return
}

func (vc *VCFile) writeToFile() (err error) {
	f, err := os.Create(vc.Filename + ".json.gz")
	if err != nil {
		return
	}
	defer f.Close()

	w := gzip.NewWriter(f)
	defer w.Close()
	enc := json.NewEncoder(w)
	enc.SetIndent("", " ")
	return enc.Encode(vc)
}

func (vc *VCFile) readFromFile() (err error) {
	fi, err := os.Open(vc.Filename)
	if err != nil {
		return
	}
	defer fi.Close()

	fz, err := gzip.NewReader(fi)
	if err != nil {
		return
	}
	defer fz.Close()

	s, err := ioutil.ReadAll(fz)
	if err != nil {
		return
	}

	err = json.Unmarshal(s, &vc)
	return
}
