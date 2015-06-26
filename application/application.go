package application

import (
	"github.com/Sirupsen/logrus"
	"github.com/carbocation/interpose"
	"github.com/didip/mcrouter-hub/handlers"
	"github.com/didip/mcrouter-hub/middlewares"
	"github.com/go-fsnotify/fsnotify"
	"github.com/gorilla/mux"
)

func New(mcRouterConfigFile, mcRouterHubCentralURL string) (*Application, error) {
	app := &Application{}
	app.McRounterConfigFile = mcRouterConfigFile
	app.McRouterHubCentralURL = mcRouterHubCentralURL

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	app.fileWatcher = watcher

	return app, nil
}

type Application struct {
	McRounterConfigFile   string
	McRouterHubCentralURL string
	fileWatcher           *fsnotify.Watcher
}

func (app *Application) WatchMcRounterConfigFile() {
	defer app.fileWatcher.Close()

	go func() {
		for {
			select {
			case event := <-app.fileWatcher.Events:
				if (event.Op&fsnotify.Create == fsnotify.Create) || (event.Op&fsnotify.Remove == fsnotify.Remove) || (event.Op&fsnotify.Write == fsnotify.Write) || (event.Op&fsnotify.Rename == fsnotify.Rename) {
					logrus.WithFields(logrus.Fields{
						"Event": event.String(),
						"File":  event.Name,
					}).Info("Config file changed")
				}
			case err := <-app.fileWatcher.Errors:
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"Error": err.Error(),
					}).Error("Error while watching config file")
				}
			}
		}
	}()

	err := app.fileWatcher.Add(app.McRounterConfigFile)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err.Error(),
			"File":  app.McRounterConfigFile,
		}).Error("Error while adding config file to watch")
	}
}

func (app *Application) MiddlewareStruct() (*interpose.Middleware, error) {
	middle := interpose.New()
	middle.Use(middlewares.SetMcRounterConfigFile(app.McRounterConfigFile))
	middle.UseHandler(app.mux())

	return middle, nil
}

func (app *Application) mux() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", handlers.GetRoot).Methods("GET")
	router.HandleFunc("/config", handlers.GetConfig).Methods("GET")
	router.HandleFunc("/config", handlers.PostConfig).Methods("POST")

	return router
}
