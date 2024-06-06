package larago

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Larago struct {
	AppName  string
	Debug    bool
	Version  string
	ErrorLog *log.Logger
	InfoLog  *log.Logger
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

	err = l.CheckDotEnv(rootPath)

	if err != nil {
		return err
	}
	err = godotenv.Load(rootPath + "/.env")

	if err != nil {
		return err
	}

	// create logs
	infoLog, errorLog := l.startLoggers()

	l.InfoLog = infoLog
	l.ErrorLog = errorLog
	l.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	l.Version = os.Getenv("VERSION")
	l.AppName = os.Getenv("APP_NAME")

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

func (l *Larago) CheckDotEnv(path string) error {
	err := l.CreateFileIfNotExists(fmt.Sprintf("%s/.env", path))
	if err != nil {
		return err
	}
	return nil
}

func (l *Larago) startLoggers() (*log.Logger, *log.Logger) {
	var infoLog *log.Logger
	var errorLog *log.Logger

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}
