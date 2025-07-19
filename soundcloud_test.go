package soundcloud_test

import (
	"testing"

	"github.com/ppalone/soundcloud"
	"github.com/stretchr/testify/assert"
)

func Test_NewClient(t *testing.T) {
	c := soundcloud.NewClient()
	assert.NotNil(t, c)
}
