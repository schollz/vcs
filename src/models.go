package vcs

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/schollz/messagebox/keypair"
	"github.com/schollz/reldel"
)

// VCSystem is the structure for the version-controlled system
type VCSystem struct {
	BaseFolder    string   `json:"base_folder"`
	EncryptedKeys []string `json:"keys"`
	CurrentBranch string   `json:"current_branch"`
	key           keypair.KeyPair
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
func Init() (vc *VCSystem, err error) {
	vc = new(VCSystem)
	vc.key, err = keypair.New()
	if err != nil {
		return
	}
	bVCKey, err := json.Marshal(vc.key)
	fmt.Println(string(bVCKey))
	encryptedKey, err := identity.Encrypt(bVCKey, identity.Public)
	if err != nil {
		return vc, err
	}
	encryptedKeyBase64 := base64.StdEncoding.EncodeToString(encryptedKey)
	vc.EncryptedKeys = []string{encryptedKeyBase64}
	vc.CurrentBranch = "master"
	vc.BaseFolder, err = os.Getwd()
	if err != nil {
		return vc, err
	}
	vc.BaseFolder, err = filepath.Abs(vc.BaseFolder)
	if err != nil {
		return vc, err
	}
	os.Mkdir(".vcs", 0755)
	// load all files in directory
	return vc, nil
}

// Commit will write the current commit to a file
func (vc *VCSystem) Commit(filename, text string) (err error) {
	// TODO add file locking

	_, filename = filepath.Split(filename)
	filenameExt := filepath.Ext(filename)
	filename = strings.TrimSuffix(filename, filenameExt)

	fileList := []string{}
	err = filepath.Walk(path.Join(vc.BaseFolder, ".vcs"), func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})
	if err != nil {
		return
	}

	for _, f := range fileList {
		var fname string
		fname, err = filepath.Abs(f)
		if err != nil {
			return
		}
		basename := strings.TrimPrefix(fname, vc.BaseFolder)
		fmt.Println(basename)
	}

	currentText := ""
	n := Block{
		Branch: vc.CurrentBranch,
		Patch:  reldel.GetPatch([]byte(currentText), []byte(text)),
	}
	h := sha256.New()
	h.Write([]byte("vcf")) // salt
	h.Write([]byte(text))
	h.Write([]byte(n.Patch.Time.String()))
	bN, err := json.Marshal(n)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s.%x%s", filename, h.Sum(nil), filenameExt), bN, 0755)
	return
}

// func (vc *VCFile) writeToFile() (err error) {
// 	f, err := os.Create(vc.Filename + ".json.gz")
// 	if err != nil {
// 		return
// 	}
// 	defer f.Close()

// 	w := gzip.NewWriter(f)
// 	defer w.Close()
// 	enc := json.NewEncoder(w)
// 	enc.SetIndent("", " ")
// 	return enc.Encode(vc)
// }

// func (vc *VCFile) readFromFile() (err error) {
// 	fi, err := os.Open(vc.Filename)
// 	if err != nil {
// 		return
// 	}
// 	defer fi.Close()

// 	fz, err := gzip.NewReader(fi)
// 	if err != nil {
// 		return
// 	}
// 	defer fz.Close()

// 	s, err := ioutil.ReadAll(fz)
// 	if err != nil {
// 		return
// 	}

// 	err = json.Unmarshal(s, &vc)
// 	return
// }
