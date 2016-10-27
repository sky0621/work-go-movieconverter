package movieconverter

import "log"

// ConvertInfo ...
type ConvertInfo struct {
	InputDir  string
	OutputDir string
	DeployDir string
	Filename  string
	Scale     string
}

func (c *ConvertInfo) inputPath() string {
	s := c.InputDir + "/" + c.Filename
	// s := filepath.Join(c.InputDir, c.Filename)
	log.Println("[inputPath]", s)
	return s
}

func (c *ConvertInfo) outputPath() string {
	s := c.OutputDir + "/" + c.Filename
	// s := filepath.Join(c.OutputDir, c.Filename)
	log.Println("[outputPath]", s)
	return s
}

func (c *ConvertInfo) cmdConvertVideo() string {
	s := "cp " + c.inputPath() + " " + c.outputPath()
	// s := "ffmpeg -i " + c.inputPath() + " -vf scale=" + c.Scale + ":-1 " + c.outputPath()
	log.Println("[cmdConvertVideo]", s)
	return s
}

func (c *ConvertInfo) cmdCreateThumbnail() string {
	s := "ls"
	// s := "ffmpeg -i " + c.outputPath() + " -ss 1 -t 1 -r 1 -f image2 " + c.outputPath() + ".jpg"
	log.Println("[cmdCreateThumbnail]", s)
	return s
}

func (c *ConvertInfo) cmdRotateThumbnail() string {
	s := "ls"
	// s := "convert -rotate 90 " + c.outputPath() + ".jpg " + c.outputPath() + "r.jpg"
	log.Println("[cmdRotateThumbnail]", s)
	return s
}

func joinPath(paths ...string) string {
	var s string
	for _, path := range paths {
		if s == "" {
			s = path
		} else {
			s = s + "/" + path
		}
	}
	// s = filepath.Join(paths...)
	log.Println("[joinPath]", s)
	return s
}
