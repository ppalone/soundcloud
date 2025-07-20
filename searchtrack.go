package soundcloud

type SearchTracksResults struct {
	Total  int
	Tracks []Track
}

type searchTracksAPIResponse struct {
	Collection   []trackAPIResponse `json:"collection"`
	TotalResults int                `json:"total_results"`
}

func (r *searchTracksAPIResponse) toResults() SearchTracksResults {
	tracks := make([]Track, 0)
	for _, t := range r.Collection {
		tracks = append(tracks, t.toTrack())
	}

	return SearchTracksResults{
		Total:  r.TotalResults,
		Tracks: tracks,
	}
}
