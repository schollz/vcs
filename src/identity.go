package vcs

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/cihub/seelog"
	"github.com/schollz/messagebox/keypair"
)

var identity keypair.KeyPair

func init() {
	b, err := ioutil.ReadFile("identity.json")
	if err != nil {
		identity, _ = keypair.New()
		b, _ = json.Marshal(identity)
		ioutil.WriteFile("identity.json", b, 0644)
	} else {
		json.Unmarshal(b, &identity)
		identity, _ = keypair.New(identity)
	}
	log.Debugf("using identity: %s", identity.Public)
}
