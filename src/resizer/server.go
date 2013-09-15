package resizer

import (
  "errors"
  "io"
  "net"
  "net/http"
  "os"
  "path/filepath"
  "strings"
)

type Server struct {
  root     string       // root direcoty
  store    *Store       // the store
  listener net.Listener // network listener
  stop     chan bool    // stop channel
  Log      chan string  // log channel
}

// Returns a new server.
func New(root, addr string) (s *Server, err error) {
  stat, err := os.Stat(root)

  if err != nil {
    return
  }

  if !stat.IsDir() {
    return nil, errors.New("This is not a directory: " + root)
  }
  listener, err := listen(addr)

  if err != nil {
    return
  }

  s = &Server{
    root:     root,
    store:    NewStore(filepath.Join(root, "images")),
    listener: listener,
    stop:     make(chan bool),
    Log:      make(chan string, 100),
  }
  return
}

// Starts the server.
func (s *Server) Run() {
  go http.Serve(s.listener, s)
}

// Stops the server.
func (s *Server) Stop() {
  s.listener.Close()
  close(s.Log)
}

// Handles a HTTP request.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  err := s.serveHTTP(w, r)

  if err == nil {
    return
  }
  s.Log <- err.Error()

  switch {
  case os.IsPermission(err):
    http.Error(w, "Permission denied.", 403)
  case os.IsNotExist(err):
    http.Error(w, "Page not found.", 404)
  default:
    http.Error(w, "Something went wrong.", 500)
  }
}

// Dispatches requets.
func (s *Server) serveHTTP(w http.ResponseWriter, r *http.Request) (err error) {
  if r.Method == "POST" {
    return s.handlePost(w, r)
  }
  return s.handleGet(w, r)
}

// Handles GET requests.
func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) (err error) {
  abs := filepath.Join(s.root, r.URL.Path)
  stat, err := os.Stat(abs)

  if err == nil {
    if stat.IsDir() {
      http.ServeFile(w, r, filepath.Join(abs, "index.html"))
    } else {
      http.ServeFile(w, r, abs)
    }
    return
  }

  if !os.IsNotExist(err) {
    return
  }
  parts := strings.Split(strings.TrimLeft(r.URL.Path, "/"), "/")

  if len(parts) < 1 {
    return
  }

  if parts[0] != "images" {
    if s.store.Exist(parts[0] + ".png") {
      http.ServeFile(w, r, filepath.Join(s.root, "index.html"))
      return nil
    }
    return os.ErrNotExist
  }

  if len(parts) < 3 {
    return os.ErrNotExist
  }
  abs, err = s.store.Get(parts[len(parts)-1], parts[1], parts[2])

  if err != nil {
    return
  }
  http.ServeFile(w, r, abs)
  return
}

// Handles POST requests.
func (s *Server) handlePost(w http.ResponseWriter, r *http.Request) (err error) {
  abs, err := s.store.Store(r.Body)

  if err != nil {
    return
  }
  rel, err := filepath.Rel(s.root, abs)

  if err != nil {
    return
  }
  _, err = io.WriteString(w, rel)
  return err
}

// Returns a new nework listener listening to addr.
func listen(addr string) (net.Listener, error) {
  if len(addr) > 0 && addr[0] == '/' {
    return net.Listen("unix", addr)
  }
  return net.Listen("tcp", addr)
}
