package modules

const defaultResponseMaxSize = 1024 * 1024 * 1024

type maxSizer struct {
	MaxSize int64
}

func (m *maxSizer) maxSize() int64 {
	if m.MaxSize != 0 {
		return m.MaxSize
	} else {
		return defaultResponseMaxSize
	}
}
