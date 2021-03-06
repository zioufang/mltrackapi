package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/zioufang/mltrackapi/pkg/api/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Server is the struct for the server
type Server struct {
	DB     *gorm.DB
	Router *chi.Mux
}

// Init initialize the server with database and router
// TODO adds postgres, mysql support
func (s *Server) Init(DbDriver, DbName string) {
	var err error
	switch DbDriver {
	case "sqlite3":
		s.DB, err = gorm.Open(sqlite.Open(DbName), &gorm.Config{})
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal(fmt.Errorf("%s is not a supported database", DbDriver))
	}
	// TODO add foreign key & index when necessary
	s.DB.AutoMigrate(
		&model.Project{},
		&model.Model{},
		&model.ModelRun{},
		&model.RunNumAttr{},
		&model.RunTag{},
	)
	s.Router = chi.NewRouter()
	// A good base middleware stack
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	// set routes
	s.SetRoutes()
}

// SetRoutes sets the routs for the server
func (s *Server) SetRoutes() {
	r := s.Router
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello Mltrack\n")
	})

	//project endpoints
	r.Route("/projects", func(r chi.Router) {
		r.Post("/", s.CreateProject)
		r.Get("/", s.GetProjectByParam)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", s.GetProjectByID)
			r.Put("/", s.UpdateProjectByID)
			r.Delete("/", s.DeleteProjectByID)
		})
		r.Get("/all", s.GetAllProjects)
	})

	// model endpoints
	r.Route("/models", func(r chi.Router) {
		r.Post("/", s.CreateModel)
		r.Get("/", s.GetModelByParam)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", s.GetModelByID)
			r.Put("/", s.UpdateModelByID)
			r.Delete("/", s.DeleteModelByID)
		})
		r.Get("/list", s.GetModelListByParam)
		r.Get("/all", s.GetAllModels)
	})

	// model run endpoints
	r.Route("/runs", func(r chi.Router) {
		r.Post("/", s.CreateModelRun)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", s.GetModelRunByID)
			r.Delete("/", s.DeleteModelRunByID)
		})
		r.Get("/list", s.GetModelRunListByParam)
	})

	// run attribute endpoints
	r.Route("/num_attrs", func(r chi.Router) {
		r.Post("/", s.CreateRunNumAttr)
		r.Get("/list", s.GetRunNumAttrListByParam)
	})

	// run tag enpoints
	r.Route("/tags", func(r chi.Router) {
		r.Post("/", s.CreateRunTag)
		r.Get("/list", s.GetRunTagListByParam)
	})

}

// Run runs the server
func (s *Server) Run(port uint) {
	fmt.Printf("Listening to port %d\n", port)
	addr := ":" + fmt.Sprint(port)
	log.Fatal(http.ListenAndServe(addr, s.Router))
}
