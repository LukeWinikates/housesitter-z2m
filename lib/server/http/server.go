package http

import (
	"LukeWinikates/january-twenty-five/lib/database"
	"LukeWinikates/january-twenty-five/lib/runtime"
	"LukeWinikates/january-twenty-five/lib/server/http/api"
	"context"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"os"
)

// cool idea - a router object passed in to all the templates
//const ROUTES

const PUT_SCHEDULES_ROUTE_PATTERN = "PUT /api/schedules/{schedule_id}"
const POST_SCHEDULES_ROUTE_PATTERN = "POST /api/schedules/"
const DELETE_SCHEDULES_ROUTE_PATTERN = "DELETE /api/schedules/{schedule_id}"
const POST_RUNNER_STATUS_PATTERN = "POST /api/runner"

var homepageTemplate *template.Template

func init() {
	homepageTemplate = template.Must(template.ParseFiles("lib/server/http/index.gohtml", "lib/server/http/device.gohtml"))
}

type Server interface {
	Serve(addr string) error
	Stop() error
}

type realServer struct {
	server *http.Server
	db     *gorm.DB
	runner runtime.Runner
}

func NewServer(db *gorm.DB, runner runtime.Runner) Server {
	return &realServer{
		db:     db,
		runner: runner,
	}
}

func (s *realServer) Stop() error {
	return s.server.Shutdown(context.TODO())
}

func (s *realServer) Serve(addr string) error {
	mux := http.NewServeMux()
	fs := os.DirFS("./public")
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.FS(fs))))
	scheduleStore := database.NewDBStore(s.db)
	deviceStore := database.NewDBDeviceStore(s.db)
	mux.HandleFunc("/", indexPage(scheduleStore, deviceStore, s.runner))
	mux.HandleFunc(PUT_SCHEDULES_ROUTE_PATTERN, api.SchedulePUTHandler(scheduleStore, s.db))
	mux.HandleFunc(POST_SCHEDULES_ROUTE_PATTERN, api.SchedulePOSTHandler(s.db))
	mux.HandleFunc(DELETE_SCHEDULES_ROUTE_PATTERN, api.ScheduleDELETEHandler(scheduleStore, s.db))
	mux.HandleFunc(POST_RUNNER_STATUS_PATTERN, api.RunnerStatePOSTHandler(s.runner))
	server := &http.Server{Addr: addr, Handler: mux}
	s.server = server
	return server.ListenAndServe()
}
