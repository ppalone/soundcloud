package soundcloud

type Preset string

const (
	MP3  Preset = "mp3"
	AAC  Preset = "aac"
	OPUS Preset = "opus"
	ABR  Preset = "abr" // Adapative Bitrate
)

func (p Preset) String() string {
	return string(p)
}
