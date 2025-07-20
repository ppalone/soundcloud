package soundcloud

type Transcoding struct {
	URL      string
	Preset   string
	Duration int
	Snipped  bool
	Format   struct {
		Protocol string
		MimeType string
	}
	Quality             string
	IsLegacyTranscoding bool
}

type transcodingAPIResponse struct {
	URL      string `json:"url"`
	Preset   string `json:"preset"`
	Duration int    `json:"duration"`
	Snipped  bool   `json:"snipped"`
	Format   struct {
		Protocol string `json:"protocol"`
		MimeType string `json:"mime_type"`
	} `json:"format"`
	Quality             string `json:"quality"`
	IsLegacyTranscoding bool   `json:"is_legacy_transcoding"`
}

func (r *transcodingAPIResponse) toTranscoding() Transcoding {
	return Transcoding{
		URL:      r.URL,
		Preset:   r.Preset,
		Duration: r.Duration,
		Snipped:  r.Snipped,
		Format: struct {
			Protocol string
			MimeType string
		}(r.Format),
		Quality:             r.Quality,
		IsLegacyTranscoding: r.IsLegacyTranscoding,
	}
}
