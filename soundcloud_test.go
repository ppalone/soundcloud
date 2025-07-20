package soundcloud

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewClient(t *testing.T) {
	c, err := NewClient()
	assert.NoError(t, err)
	assert.NotNil(t, c)
}

func Test_SearchTracks(t *testing.T) {
	c, err := NewClient()
	assert.NoError(t, err)

	t.Run("without options", func(t *testing.T) {
		q := "nujabes"
		res, err := c.SearchTracks(context.Background(), q)
		assert.NoError(t, err)
		assert.NotEmpty(t, res.Tracks)
		assert.Equal(t, 20, len(res.Tracks))

		for _, track := range res.Tracks {
			assert.NotEmpty(t, track.Transcodings)
		}
	})

	t.Run("with limit", func(t *testing.T) {
		q := "nujabes"
		limit := 100
		res, err := c.SearchTracks(context.Background(), q, WithLimit(limit))
		assert.NoError(t, err)
		assert.Len(t, res.Tracks, limit)

		for _, track := range res.Tracks {
			assert.NotEmpty(t, track.Transcodings)
		}
	})

}
