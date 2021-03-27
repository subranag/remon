package remon

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"

	"github.com/gorilla/websocket"
)

const (
	statsWriteFrequency = time.Millisecond * 500
	writeWait           = time.Second * 2
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Server interface {
	Start() error
}

// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type spaHandler struct {
	staticPath string
	indexPath  string
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Println("Iam here")
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

type muxServer struct {
	router *mux.Router
	addr   string
}

func (s *muxServer) Start() error {
	srv := &http.Server{
		Handler: s.router,
		Addr:    s.addr,
	}
	return srv.ListenAndServe()
}

func serveWsCpu(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			fmt.Println(err)
		}
		return
	}

	go writer(ws)
	reader(ws)
}

func reader(ws *websocket.Conn) {
	ws.SetReadDeadline(time.Time{})
	for {
		if _, _, err := ws.ReadMessage(); err != nil {
			ws.Close()
			break
		}
	}
}

func writer(ws *websocket.Conn) {
	writeTicker := time.NewTicker(statsWriteFrequency)

	defer func() {
		writeTicker.Stop()
	}()

	stats := make(CpuStats)
	prevStats := make(CpuStats)
	cpuStats, err := NewCpuStatsReader()
	cpuUtil := make(map[string]float64)

	if err != nil {
		closeMsg := websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error())
		ws.WriteMessage(websocket.CloseMessage, closeMsg)
	}

	for {
		select {
		case <-writeTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			cpuStats.Read(stats)
			if len(prevStats) > 0 {
				for k, v := range stats {
					cpuUtil[k] = v.Utilization(prevStats[k])
				}
				//write map to websocket
				ws.WriteJSON(cpuUtil)
			}
			stats.CopyTo(prevStats)
		}
	}
}

func NewServer(addr, staticPath, indexPath string) Server {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	})

	r.HandleFunc("/api/ws/cpu", serveWsCpu)

	spa := spaHandler{staticPath: staticPath, indexPath: indexPath}
	r.PathPrefix("/").Handler(spa)
	return &muxServer{router: r, addr: addr}
}
