package server

type (
	Config struct {
		Host        string
		Port        int
		Compression Compression
	}

	Compression struct {
		Threshold int
		Level     int
	}
)
