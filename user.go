package soundcloud

import "time"

type User struct {
	ID                 int
	Username           string
	AvatarURL          string
	Description        string
	FirstName          string
	LastName           string
	FullName           string
	FollowersCount     int
	FollowingsCount    int
	TrackCount         int
	LikesCount         int
	PlaylistLikesCount int
	PermalinkURL       string
	CreatedAt          time.Time
	Kind               string
}

type userAPIResponse struct {
	AvatarURL            string    `json:"avatar_url"`
	City                 string    `json:"city"`
	CommentsCount        int       `json:"comments_count"`
	CountryCode          string    `json:"country_code"`
	CreatedAt            time.Time `json:"created_at"`
	CreatorSubscriptions []struct {
		Product struct {
			ID string `json:"id"`
		} `json:"product"`
	} `json:"creator_subscriptions"`
	CreatorSubscription struct {
		Product struct {
			ID string `json:"id"`
		} `json:"product"`
	} `json:"creator_subscription"`
	Description        string    `json:"description"`
	FollowersCount     int       `json:"followers_count"`
	FollowingsCount    int       `json:"followings_count"`
	FirstName          string    `json:"first_name"`
	FullName           string    `json:"full_name"`
	GroupsCount        int       `json:"groups_count"`
	ID                 int       `json:"id"`
	Kind               string    `json:"kind"`
	LastModified       time.Time `json:"last_modified"`
	LastName           string    `json:"last_name"`
	LikesCount         int       `json:"likes_count"`
	PlaylistLikesCount int       `json:"playlist_likes_count"`
	Permalink          string    `json:"permalink"`
	PermalinkURL       string    `json:"permalink_url"`
	PlaylistCount      int       `json:"playlist_count"`
	TrackCount         int       `json:"track_count"`
	URI                string    `json:"uri"`
	Urn                string    `json:"urn"`
	Username           string    `json:"username"`
	Verified           bool      `json:"verified"`
	Visuals            struct {
		Urn     string `json:"urn"`
		Enabled bool   `json:"enabled"`
		Visuals []struct {
			Urn       string `json:"urn"`
			EntryTime int    `json:"entry_time"`
			VisualURL string `json:"visual_url"`
		} `json:"visuals"`
		Tracking any `json:"tracking"`
	} `json:"visuals"`
	Badges struct {
		Pro            bool `json:"pro"`
		CreatorMidTier bool `json:"creator_mid_tier"`
		ProUnlimited   bool `json:"pro_unlimited"`
		Verified       bool `json:"verified"`
	} `json:"badges"`
	StationUrn       string `json:"station_urn"`
	StationPermalink string `json:"station_permalink"`
	DateOfBirth      struct {
		Month int `json:"month"`
		Year  int `json:"year"`
		Day   int `json:"day"`
	} `json:"date_of_birth"`
}

func (r *userAPIResponse) toUser() User {
	return User{
		ID:                 r.ID,
		Username:           r.Username,
		AvatarURL:          r.AvatarURL,
		Description:        r.Description,
		FirstName:          r.FirstName,
		LastName:           r.LastName,
		FullName:           r.FullName,
		FollowersCount:     r.FollowersCount,
		FollowingsCount:    r.FollowingsCount,
		TrackCount:         r.TrackCount,
		LikesCount:         r.LikesCount,
		PlaylistLikesCount: r.PlaylistCount,
		PermalinkURL:       r.PermalinkURL,
		CreatedAt:          r.CreatedAt,
		Kind:               r.Kind,
	}
}
