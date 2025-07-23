package soundcloud_test

import (
	"context"
	"strings"
	"testing"

	"github.com/ppalone/soundcloud"
	"github.com/stretchr/testify/assert"
)

func Test_NewClient(t *testing.T) {
	c, err := soundcloud.NewClient()
	assert.NoError(t, err)
	assert.NotNil(t, c)
}

func Test_SearchTracks(t *testing.T) {
	c, err := soundcloud.NewClient()
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
		res, err := c.SearchTracks(context.Background(), q, soundcloud.WithLimit(limit))
		assert.NoError(t, err)
		assert.Len(t, res.Tracks, limit)

		for _, track := range res.Tracks {
			assert.NotEmpty(t, track.Transcodings)
		}
	})

	t.Run("with client id option", func(t *testing.T) {
		c2, err := soundcloud.NewClient(soundcloud.WithClientID(c.ClientId()))
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
	c, err := soundcloud.NewClient()
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

func Test_Stream(t *testing.T) {
	c, err := soundcloud.NewClient()
	assert.NoError(t, err)

	q := "Martin Garrix Animals"
	limit := 1
	res, err := c.SearchTracks(context.Background(), q, soundcloud.WithLimit(limit))
	assert.NoError(t, err)
	assert.Len(t, res.Tracks, limit)

	track := res.Tracks[0]

	t.Run("with valid hls transcoding", func(t *testing.T) {
		var transcoding soundcloud.Transcoding
		for _, t := range track.Transcodings {
			if strings.HasPrefix(t.Preset, soundcloud.MP3.String()) && t.Format.Protocol == soundcloud.HLS.String() {
				transcoding = t
				break
			}
		}
		assert.NotNil(t, transcoding)

		stream, err := c.GetStream(context.Background(), transcoding)
		assert.NoError(t, err)
		assert.NotNil(t, stream)
		stream.Close()
	})

	t.Run("with valid progressive transcoding", func(t *testing.T) {
		var transcoding soundcloud.Transcoding
		for _, t := range track.Transcodings {
			if strings.HasPrefix(t.Preset, soundcloud.MP3.String()) && t.Format.Protocol == soundcloud.PROGRESSIVE.String() {
				transcoding = t
				break
			}
		}
		assert.NotNil(t, transcoding)

		stream, err := c.GetStream(context.Background(), transcoding)
		assert.NoError(t, err)
		assert.NotNil(t, stream)
		stream.Close()
	})

	t.Run("with invalid progressive transcoding", func(t *testing.T) {
		transcoding := soundcloud.Transcoding{
			URL:      "https://api-v2.soundcloud.com/media/soundcloud:tracks:98081145/d27e0e2b-16e6-4fb2-b2a2-49bb6ad7e489/stream/hls", // invalid url
			Preset:   "mp3_0_1",
			Duration: 304300,
			Snipped:  false,
			Format: struct {
				Protocol string
				MimeType string
			}{
				Protocol: "hls",
				MimeType: "audio/mpeg",
			},
			Quality:             "sq",
			IsLegacyTranscoding: true,
		}
		stream, err := c.GetStream(context.Background(), transcoding)
		assert.ErrorContains(t, err, "404")
		assert.Nil(t, stream)
	})
}

func Test_GetStreamById(t *testing.T) {
	c, err := soundcloud.NewClient()
	assert.NoError(t, err)

	t.Run("with valid id and available stream options", func(t *testing.T) {
		id := 98081145
		stream, err := c.GetStreamById(context.Background(), id, soundcloud.WithPreset(soundcloud.MP3), soundcloud.WithProtocol(soundcloud.PROGRESSIVE))
		assert.NoError(t, err)
		assert.NotNil(t, stream)
	})

	t.Run("with valid id and unavailable stream options", func(t *testing.T) {
		id := 98081145
		stream, err := c.GetStreamById(context.Background(), id, soundcloud.WithPreset(soundcloud.AAC), soundcloud.WithProtocol(soundcloud.PROGRESSIVE))
		assert.ErrorContains(t, err, "transcoding with preset aac and protocol progressive not found for track")
		assert.Nil(t, stream)
	})
}
