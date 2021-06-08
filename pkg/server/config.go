package server

type (
	Config struct {
		Host        string
		Port        int
		World       WorldConf
		Compression CompressionConf
	}

	WorldConf struct {
		Schematic      string
		RenderDistance int
	}

	CompressionConf struct {
		Threshold int
		Level     int
	}
)
