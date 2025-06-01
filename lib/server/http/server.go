package http

import (
	"LukeWinikates/january-twenty-five/lib/schedule"
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
const PUT_DEVICE_SETTINGS_ROUTE_PATTERN = "PUT /api/schedules/{schedule_id}/device_settings/{device_id}"
const POST_DEVICE_SETTINGS_ROUTE_PATTERN = "POST /api/schedules/{schedule_id}/device_settings"

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
	scheduleStore, err := schedule.NewDBStore(s.db, true)
	if err != nil {
		return err
	}
	deviceStore := schedule.NewDBDeviceStore(s.db)
	mux.HandleFunc("/", indexPage(scheduleStore, deviceStore))
	mux.HandleFunc(PUT_SCHEDULES_ROUTE_PATTERN, api.SchedulePUTHandler(scheduleStore))
	mux.HandleFunc(POST_SCHEDULES_ROUTE_PATTERN, api.SchedulePOSTHandler(scheduleStore))
	mux.HandleFunc(PUT_DEVICE_SETTINGS_ROUTE_PATTERN, api.ScheduleDevicePUTHandler(scheduleStore))
	mux.HandleFunc(POST_DEVICE_SETTINGS_ROUTE_PATTERN, api.ScheduleDevicePOSTHandler(scheduleStore, deviceStore))
	server := &http.Server{Addr: addr, Handler: mux}
	s.server = server
	return server.ListenAndServe()
}
