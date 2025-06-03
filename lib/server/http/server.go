package http

import (
	"LukeWinikates/january-twenty-five/lib/database"
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
}

func NewServer(db *gorm.DB) Server {
	return &realServer{
		db: db,
	}
}

func (s *realServer) Stop() error {
	return s.server.Shutdown(context.TODO())
}

func (s *realServer) Serve(addr string) error {
	mux := http.NewServeMux()
	fs := os.DirFS("./public")
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.FS(fs))))
	scheduleStore, err := database.NewDBStore(s.db, false)
	if err != nil {
		return err
	}
	deviceStore := database.NewDBDeviceStore(s.db)
	mux.HandleFunc("/", indexPage(scheduleStore, deviceStore))
	mux.HandleFunc(PUT_SCHEDULES_ROUTE_PATTERN, api.SchedulePUTHandler(scheduleStore, s.db))
	mux.HandleFunc(POST_SCHEDULES_ROUTE_PATTERN, api.SchedulePOSTHandler(s.db))
	server := &http.Server{Addr: addr, Handler: mux}
	s.server = server
	return server.ListenAndServe()
}
