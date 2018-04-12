package vcs

import (
	"fmt"
	"testing"

	log "github.com/cihub/seelog"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	defer log.Flush()
	vc, err := Init()
	assert.Nil(t, err)
	fmt.Println(vc)
}
