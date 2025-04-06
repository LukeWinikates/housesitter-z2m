package http

import (
	"LukeWinikates/january-twenty-five/lib/schedule"
	"LukeWinikates/january-twenty-five/lib/server/http/api"
	"context"
	"html/template"
	"net/http"
	"os"
)

// cool idea - a router object passed in to all the templates
//const ROUTES

const PUT_SCHEDULES_ROUTE_PATTERN = "PUT /api/schedules/{schedule_id}"
const PUT_DEVICE_SETTINGS_ROUTE_PATTERN = "PUT /api/schedules/{schedule_id}/device_settings/{device_id}"

var homepageTemplate *template.Template

func init() {
	homepageTemplate = template.Must(template.ParseFiles("lib/server/http/index.gohtml"))
}

type Server interface {
	Serve(addr string) error
	Stop() error
}

type realServer struct {
	server *http.Server
}

func NewServer() Server {
	return &realServer{}
}

func (s *realServer) Stop() error {
	return s.server.Shutdown(context.TODO())
}

func (s *realServer) Serve(addr string) error {
	mux := http.NewServeMux()
	fs := os.DirFS("./public")
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.FS(fs))))
	scheduleStore := schedule.NewStore()
	mux.HandleFunc("/", indexPage(scheduleStore))
	mux.HandleFunc(PUT_SCHEDULES_ROUTE_PATTERN, api.SchedulePutHandler(scheduleStore))
	mux.HandleFunc(PUT_DEVICE_SETTINGS_ROUTE_PATTERN, api.ScheduleDevicePutHandler(scheduleStore))
	server := &http.Server{Addr: addr, Handler: mux}
	s.server = server
	return server.ListenAndServe()
}
