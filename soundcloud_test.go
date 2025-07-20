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

	t.Run("with client id option", func(t *testing.T) {
		c2, err := NewClient(WithClientID(c.clientId))
		assert.NoError(t, err)

		q := "monstercat"
		res, err := c2.SearchTracks(context.Background(), q)
		assert.NoError(t, err)
		assert.NotEmpty(t, res.Tracks)
		assert.Equal(t, 20, len(res.Tracks))

		for _, track := range res.Tracks {
			assert.NotEmpty(t, track.Transcodings)
		}
	})

}

func Test_GetTrackById(t *testing.T) {
	c, err := NewClient()
	assert.NoError(t, err)

	t.Run("with valid id", func(t *testing.T) {
		id := 98081145 // Martin Garrix - Animals
		res, err := c.GetTrackById(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, res.ID, id)
		assert.Contains(t, res.Title, "Animals")
		assert.NotEmpty(t, res.Transcodings)
	})

	t.Run("with invalid id", func(t *testing.T) {
		id := 0
		_, err := c.GetTrackById(context.Background(), id)
		assert.ErrorContains(t, err, "invalid id")
	})
}
