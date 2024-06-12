package larago

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/go-chi/chi/v5"
	"github.com/igorsouza-dev/larago/render"
	"github.com/joho/godotenv"
)

type Larago struct {
	AppName  string
	Debug    bool
	Version  string
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	Routes   *chi.Mux
	Render   *render.Render
	RootPath string
	JetViews *jet.Set
	config   config
}

type config struct {
	port     string
	renderer string
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

	err = l.checkDotEnv(rootPath)

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
	l.Routes = l.routes().(*chi.Mux)
	l.RootPath = rootPath

	l.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
	}

	views := jet.NewSet(
		jet.NewOSFileSystemLoader(fmt.Sprintf("%s/views", rootPath)),
		jet.InDevelopmentMode(),
	)

	l.JetViews = views

	l.createRenderer()

	return nil
}

func (l *Larago) Init(p initPaths) error {
	root := p.rootPath

	for _, path := range p.folderNames {
		err := l.createDirIfNotExist(root + "/" + path)
		if err != nil {
			return err
		}
	}
	return nil
}
func (l *Larago) ListenAndServe() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", l.config.port),
		Handler:      l.Routes,
		ErrorLog:     l.ErrorLog,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}
	l.InfoLog.Printf("Starting server on port %s", l.config.port)
	err := srv.ListenAndServe()

	if err != nil {
		return err
	}

	return nil
}
func (l *Larago) checkDotEnv(path string) error {
	err := l.createFileIfNotExists(fmt.Sprintf("%s/.env", path))
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

func (l *Larago) createRenderer() {
	renderer := render.Render{
		Renderer: l.config.renderer,
		Port:     l.config.port,
		RootPath: l.RootPath,
		JetViews: l.JetViews,
	}

	l.Render = &renderer
}
