package gen

type Sources struct {
	Generated  []TextFile
	Scaffolded []TextFile
}

func NewSources() *Sources {
	return &Sources{[]TextFile{}, []TextFile{}}
}

func (sources *Sources) AddScaffolded(files ...*TextFile) {
	for _, file := range files {
		sources.Scaffolded = append(sources.Scaffolded, *file)
	}
}

func (sources *Sources) AddScaffoldedAll(files []TextFile) {
	for _, file := range files {
		sources.Scaffolded = append(sources.Scaffolded, file)
	}
}

func (sources *Sources) AddGenerated(files ...*TextFile) {
	for _, file := range files {
		sources.Generated = append(sources.Generated, *file)
	}
}

func (sources *Sources) AddGeneratedAll(files []TextFile) {
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
