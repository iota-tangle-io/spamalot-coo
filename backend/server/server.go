package server

import (
	"github.com/iota-tangle-io/spamalot-coo/backend/controllers"
	"github.com/iota-tangle-io/spamalot-coo/backend/routers"
	"github.com/iota-tangle-io/spamalot-coo/backend/utilities"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/facebookgo/inject"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/globalsign/mgo"
	"html/template"
	"io"
	"os"
	"time"
	"github.com/iota-tangle-io/spamalot-coo/api"
	"github.com/iota-tangle-io/spamalot-coo/backend/server/config"
)

type TemplateRendered struct {
	templates *template.Template
}

func (t *TemplateRendered) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type Server struct {
	Config    *config.Configuration
	WebEngine *echo.Echo
	Mongo     *mgo.Session
}

func (server *Server) Start() {
	start := time.Now().UnixNano()

	// load config
	configuration := config.LoadConfig()
	server.Config = configuration
	appConfig := server.Config.App
	httpConfig := server.Config.Net.HTTP

	// init logger
	utilities.Debug = appConfig.Verbose
	logger, err := utilities.GetLogger("app")
	if err != nil {
		panic(err)
	}
	logger.Info("booting up app...")

	// connect to mongo
	mongo := &mgo.Session{}
	var mongoConnErr error
	mongoConfig := server.Config.Net.Database.Mongo
	if mongoConfig.Use {
		mongo, mongoConnErr = connectMongoDB(server.Config.Net.Database.Mongo)
		if mongoConnErr != nil {
			panic(mongoConnErr)
		}
		if err = mongo.Ping(); err != nil {
			panic(err)
		}
		server.Mongo = mongo
		logger.Info("connection to MongoDB established")
	} else {
		logger.Info("MongoDB connection disabled")
	}

	// init web server
	e := echo.New()
	e.HideBanner = true
	server.WebEngine = e
	if httpConfig.LogRequests {
		requestLogFile, err := os.Create(fmt.Sprintf("./logs/requests.log"))
		if err != nil {
			panic(err)
		}
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Output: requestLogFile}))
		e.Logger.SetLevel(3)
	}

	// load html files
	e.Renderer = &TemplateRendered{
		templates: template.Must(template.ParseGlob(fmt.Sprintf("%s/*.html", httpConfig.Assets.HTML))),
	}

	// coordinator
	coordinator := &api.Coordinator{Config: configuration.Net.Coordinator}

	// asset paths
	e.Static("/assets", httpConfig.Assets.Static)
	e.File("/favicon.ico", httpConfig.Assets.Favicon)

	// create controllers
	appCtrl := &controllers.AppCtrl{}
	configCtrl := &controllers.ConfigCtrl{}
	instanceCtrl := &controllers.InstanceCtrl{}
	controllers := []controllers.Controller{
		appCtrl, configCtrl, instanceCtrl,
	}

	// create routers
	indexRouter := &routers.IndexRouter{}
	configRouter := &routers.ConfigRouter{}
	instanceRouter := &routers.InstanceRouter{}
	rters := []routers.Router{indexRouter, configRouter, instanceRouter}

	// create injection graph for automatic dependency injection
	g := inject.Graph{}

	// add various objects to the graph
	if err = g.Provide(
		&inject.Object{Value: e},
		&inject.Object{Value: coordinator},
		&inject.Object{Value: mongo},
		&inject.Object{Value: appConfig.Dev, Name: "dev"},
	); err != nil {
		panic(err)
	}

	// add controllers to graph
	for _, controller := range controllers {
		if err = g.Provide(&inject.Object{Value: controller}); err != nil {
			panic(err)
		}
	}

	// add routers to graph
	for _, router := range rters {
		if err = g.Provide(&inject.Object{Value: router}); err != nil {
			panic(err)
		}
	}

	// run dependency injection
	if err = g.Populate(); err != nil {
		panic(err)
	}

	// init controllers
	for _, controller := range controllers {
		if err = controller.Init(); err != nil {
			panic(err)
		}
	}
	logger.Info("initialised controllers")

	// init routers
	for _, router := range rters {
		router.Init()
	}
	logger.Info("initialised routers")

	// firing up coordinator
	coordinator.Run()

	// boot up server
	go e.Start(httpConfig.Address)

	// finish
	delta := (time.Now().UnixNano() - start) / 1000000
	logger.Info(fmt.Sprintf("SPA ready, listening on %s", httpConfig.Address), "startup", delta)

}

func (server *Server) Shutdown(timeout time.Duration) {
	select {
	case <-time.After(timeout):
	}
}

func connectMongoDB(config config.MongoDBConfig) (*mgo.Session, error) {
	var session *mgo.Session
	var err error
	if config.Auth {
		cred := &mgo.Credential{
			Username:  config.Username,
			Password:  config.Password,
			Mechanism: config.Mechanism,
			Source:    config.Source,
		}
		session, err = mgo.Dial(config.Address)
		if err = session.Login(cred); err != nil {
			return nil, err
		}
	} else {
		session, err = mgo.Dial(config.Address)
	}
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	session.SetSafe(&mgo.Safe{})
	return session, nil
}