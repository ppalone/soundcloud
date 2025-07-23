package soundcloud

type Protocol string

const (
	HLS         Protocol = "hls"
	PROGRESSIVE Protocol = "progressive"
)

func (p Protocol) String() string {
	return string(p)
}
