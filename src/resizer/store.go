package resizer

import (
  "crypto/sha1"
  "encoding/hex"
  "errors"
  "github.com/nfnt/resize"
  "image"
  _ "image/gif"
  _ "image/jpeg"
  "image/png"
  "io"
  "os"
  "path/filepath"
  "strconv"
)

// Map of interpolation methods.
var interpolations = map[string]resize.InterpolationFunction{
  "bicubic":          resize.Bicubic,
  "bilinear":         resize.Bilinear,
  "lanczos2":         resize.Lanczos2,
  "lanczos3":         resize.Lanczos3,
  "mitchelnetravali": resize.MitchellNetravali,
  "nearestneighbor":  resize.NearestNeighbor,
}

type Store struct {
  root string // root directory
}

// Returns a new store.
func NewStore(root string) *Store {
  return &Store{root: root}
}

// Stores a new image and returns the absolute path to the stored file.
func (s *Store) Store(r io.Reader) (string, error) {
  hsh := sha1.New()
  img, _, err := image.Decode(io.TeeReader(r, hsh))

  if err != nil {
    return "", err
  }
  return s.store("orig", hex.EncodeToString(hsh.Sum(nil)), img)
}

// Returns the absolute path to the image in the store.
func (s *Store) Get(name, size, interp string) (path string, err error) {
  path = s.hashedPath(name, filepath.Join(size, interp))
  _, err = os.Stat(path)

  if err == nil || !os.IsNotExist(err) {
    return
  }
  return s.generate(name, size, interp)
}

// Checks if an image available.
func (s Store) Exist(name string) bool {
  _, err := os.Stat(s.hashedPath(name, "orig"))
  return err == nil
}

// Generates a resized image image from the original.
func (s *Store) generate(name, size, interp string) (path string, err error) {
  uSize, err := strconv.ParseUint(size, 10, 0)

  if err != nil {
    return
  }
  interpolation, ok := interpolations[interp]

  if !ok {

    return "", errors.New("Unknown interpolation: " + interp)
  }
  f, err := os.Open(s.hashedPath(name, "orig"))

  if err != nil {
    return
  }
  defer f.Close()

  img, _, err := image.Decode(f)

  if err != nil {
    return
  }
  newImg := resize.Resize(uint(uSize), 0, img, interpolation)

  return s.store(filepath.Join(size, interp), name, newImg)
}

// Stores an image.
func (s *Store) store(dir, name string, img image.Image) (path string, err error) {
  if len(filepath.Ext(name)) == 0 {
    name += ".png"
  }
  path = s.hashedPath(name, dir)

  if err = os.MkdirAll(filepath.Dir(path), 0755); err != nil {
    return
  }
  f, err := os.Create(path)

  if err != nil {
    return
  }
  defer f.Close()

  return path, png.Encode(f, img)
}

// Returns an abs path for an image file.
func (s *Store) hashedPath(name, dir string) string {
  return filepath.Join(s.root, dir, string(name[0]), string(name[1]), name)
}
