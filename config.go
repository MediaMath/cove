package main

type GoCovConfig struct {
	OutputDir     OutputFiles
	ShortTests    bool
	OpenInBrowser bool
	KeepProfile   bool
}

func Config(outputDir string, short bool, keep bool) (*GoCovConfig, error) {
	if outputDir == "" {
		return BrowserConfig(short, keep)
	} else {
		return OutputConfig(outputDir, short, keep)
	}
}

func BrowserConfig(short bool, keep bool) (*GoCovConfig, error) {
	out, outerr := TempOutputFiles()
	if outerr != nil {
		return nil, outerr
	}

	return &GoCovConfig{out, short, true, keep}, nil
}

func OutputConfig(outputDir string, short bool, keep bool) (*GoCovConfig, error) {
	out, outerr := ExistingDirOutputFiles(outputDir)
	if outerr != nil {
		return nil, outerr
	}

	return &GoCovConfig{out, short, false, keep}, nil
}
