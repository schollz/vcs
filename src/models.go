package vcf

import (
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/schollz/reldel"
)

// VCFile is the structure for a version-controlled file
type VCFile struct {
	Filename      string           `json:"filename"`
	CurrentText   string           `json:"current_text"`
	CurrentHash   string           `json:"current_hash"`
	CurrentBranch string           `json:"current_branch"`
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
func Init(filename string) (*VCFile, error) {
	vc := new(VCFile)
	vc.Filename = filename
	errOpen := vc.readFromFile()
	if errOpen == nil {
		return vc, nil
	}

	vc.CurrentBranch = "master"
	vc.CurrentText = ""
	vc.StartingBlock = ""
	vc.BlockMap = make(map[string]Block)
	_, err := vc.Commit("")
	vc.StartingBlock = vc.CurrentHash
	return vc, err
}

// Commit will write the current commit to a file
func (vc *VCFile) Commit(text string) (blockHash string, err error) {
	// TODO add file locking
	n := Block{
		Branch:       vc.CurrentBranch,
		PreviousHash: vc.CurrentHash,
		Patch:        reldel.GetPatch(vc.CurrentText, text),
	}
	h := sha256.New()
	h.Write([]byte("vcf")) // salt
	h.Write([]byte(text))
	h.Write([]byte(n.Patch.Time.String()))
	blockHash = fmt.Sprintf("%x", h.Sum(nil))
	vc.StartingBlock = blockHash
	vc.CurrentHash = blockHash
	vc.CurrentText = text
	vc.BlockMap[vc.CurrentHash] = n

	err = vc.writeToFile()
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
