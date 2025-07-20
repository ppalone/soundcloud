package soundcloud

import "time"

type Track struct {
	ID                 int
	Title              string
	Description        string
	ArtworkURL         string
	Duration           int
	Genre              string
	CommentCount       int
	LikesCount         int
	TrackAuthorization string
	Transcodings       []Transcoding
	User               User
	Kind               string
}

type trackAPIResponse struct {
	ArtworkURL    string    `json:"artwork_url"`
	CommentCount  int       `json:"comment_count"`
	CreatedAt     time.Time `json:"created_at"`
	Description   string    `json:"description"`
	Duration      int       `json:"duration"`
	Genre         string    `json:"genre"`
	ID            int       `json:"id"`
	Kind          string    `json:"kind"`
	LabelName     string    `json:"label_name"`
	LikesCount    int       `json:"likes_count"`
	Permalink     string    `json:"permalink"`
	PermalinkURL  string    `json:"permalink_url"`
	PlaybackCount int       `json:"playback_count"`
	Public        bool      `json:"public"`
	RepostsCount  int       `json:"reposts_count"`
	Sharing       string    `json:"sharing"`
	Title         string    `json:"title"`
	URI           string    `json:"uri"`
	Urn           string    `json:"urn"`
	UserID        int       `json:"user_id"`
	WaveformURL   string    `json:"waveform_url"`
	Media         struct {
		Transcodings []transcodingAPIResponse `json:"transcodings"`
	} `json:"media"`
	TrackAuthorization string          `json:"track_authorization"`
	User               userAPIResponse `json:"user"`
}

func (r *trackAPIResponse) toTrack() Track {
	transcodings := make([]Transcoding, 0)
	for _, t := range r.Media.Transcodings {
		transcodings = append(transcodings, t.toTranscoding())
	}

	return Track{
		ID:                 r.ID,
		Title:              r.Title,
		Description:        r.Description,
		ArtworkURL:         r.ArtworkURL,
		Duration:           r.Duration,
		Genre:              r.Genre,
		CommentCount:       r.CommentCount,
		LikesCount:         r.LikesCount,
		TrackAuthorization: r.TrackAuthorization,
		Transcodings:       transcodings,
		User:               r.User.toUser(),
		Kind:               r.Kind,
	}
}
