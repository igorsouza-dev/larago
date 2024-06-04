package larago

type Larago struct {
	AppName string
	Debug   bool
	Version string
}

func (l *Larago) New(rootPath string) error {
	pathConfig := initPaths{
		rootPath:    rootPath,
		folderNames: []string{"handlers", "migrations", "views", "data", "public", "tmp", "logs", "middlewares"},
	}

	err := l.Init(pathConfig)
	if err != nil {
		return err
	}

	return nil
}

func (l *Larago) Init(p initPaths) error {
	root := p.rootPath

	for _, path := range p.folderNames {
		err := l.CreateDirIfNotExist(root + "/" + path)
		if err != nil {
			return err
		}
	}
	return nil
}
