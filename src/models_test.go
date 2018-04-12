package vcf

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	vc, err := Init("test.txt")
	assert.Nil(t, err)
	fmt.Println(vc)
}
