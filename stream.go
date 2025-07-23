package soundcloud

type streamOptions struct {
	preset   Preset
	protocol Protocol
}

type StreamOption func(o *streamOptions)

func defaultStreamOptions() *streamOptions {
	return &streamOptions{
		preset:   AAC,
		protocol: HLS,
	}
}

func WithPreset(p Preset) StreamOption {
	return func(o *streamOptions) {
		o.preset = p
	}
}

func WithProtocol(p Protocol) StreamOption {
	return func(o *streamOptions) {
		o.protocol = p
	}
}
