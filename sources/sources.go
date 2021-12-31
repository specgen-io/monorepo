package sources

type Sources struct {
	Generated  []CodeFile
	Scaffolded []CodeFile
}

func NewSources() *Sources {
	return &Sources{[]CodeFile{}, []CodeFile{}}
}

func (sources *Sources) AddScaffolded(files ...*CodeFile) {
	for _, file := range files {
		if file != nil {
			sources.Scaffolded = append(sources.Scaffolded, *file)
		}
	}
}

func (sources *Sources) AddScaffoldedAll(files []CodeFile) {
	for _, file := range files {
		sources.Scaffolded = append(sources.Scaffolded, file)
	}
}

func (sources *Sources) AddGenerated(files ...*CodeFile) {
	for _, file := range files {
		if file != nil {
			sources.Generated = append(sources.Generated, *file)
		}
	}
}

func (sources *Sources) AddGeneratedAll(files []CodeFile) {
	for _, file := range files {
		sources.Generated = append(sources.Generated, file)
	}
}

func (sources *Sources) Write(overwriteAll bool) error {
	err := WriteFiles(sources.Scaffolded, false || overwriteAll)
	if err != nil {
		return err
	}

	err = WriteFiles(sources.Generated, true)
	if err != nil {
		return err
	}
	return nil
}
