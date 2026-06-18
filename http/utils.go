package fbhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/filebrowser/filebrowser/v2/users"
	libErrors "github.com/filebrowser/filebrowser/v2/errors"
	imgErrors "github.com/filebrowser/filebrowser/v2/img"
	"github.com/spf13/afero"
)

// Global atomic counter for rename sequence numbers
var renameCounter uint64

func renderJSON(w http.ResponseWriter, _ *http.Request, data interface{}) (int, error) {
	marsh, err := json.Marshal(data)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if _, err := w.Write(marsh); err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}

func errToStatus(err error) int {
	switch {
	case err == nil:
		return http.StatusOK
	case os.IsPermission(err):
		return http.StatusForbidden
	case os.IsNotExist(err), errors.Is(err, libErrors.ErrNotExist):
		return http.StatusNotFound
	case os.IsExist(err), errors.Is(err, libErrors.ErrExist):
		return http.StatusConflict
	case errors.Is(err, libErrors.ErrPermissionDenied):
		return http.StatusForbidden
	case errors.Is(err, libErrors.ErrInvalidRequestParams):
		return http.StatusBadRequest
	case errors.Is(err, libErrors.ErrRootUserDeletion):
		return http.StatusForbidden
	case errors.Is(err, imgErrors.ErrImageTooLarge):
		return http.StatusRequestEntityTooLarge
	default:
		return http.StatusInternalServerError
	}
}

// This is an adaptation if http.StripPrefix in which we don't
// return 404 if the page doesn't have the needed prefix.
func stripPrefix(prefix string, h http.Handler) http.Handler {
	if prefix == "" || prefix == "/" {
		return h
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(r.URL.Path, prefix)
		rp := strings.TrimPrefix(r.URL.RawPath, prefix)

		// If the path is exactly the prefix (no trailing slash), redirect to
		// the prefix with a trailing slash so the router receives "/" instead
		// of "", which would otherwise cause a redirect to the site root.
		if p == "" {
			http.Redirect(w, r, prefix+"/", http.StatusMovedPermanently)
			return
		}

		r2 := new(http.Request)
		*r2 = *r
		r2.URL = new(url.URL)
		*r2.URL = *r.URL
		r2.URL.Path = p
		r2.URL.RawPath = rp
		h.ServeHTTP(w, r2)
	})
}

// applyAutoRename applies the user's auto-rename pattern to a filename.
// The pattern supports the following placeholders:
//   {timestamp} - current time formatted as YYYYMMDDHHmmss
//   {date} - current date formatted as YYYYMMDD
//   {time} - current time formatted as HHmmss
//   {year} - current year (YYYY)
//   {month} - current month (MM)
//   {day} - current day (DD)
//   {hour} - current hour (HH)
//   {minute} - current minute (mm)
//   {second} - current second (ss)
//   {name} - original filename without extension
//   {ext} - file extension (including dot)
//   {n} - auto-increment number (1, 2, 3, ...)
//
// If pattern is empty, defaults to "{timestamp}{ext}".
func applyAutoRename(pattern, originalName string, counter int) string {
	if pattern == "" {
		pattern = "{timestamp}{ext}"
	}

	now := time.Now()
	ext := filepath.Ext(originalName)
	name := strings.TrimSuffix(originalName, ext)

	replacements := map[string]string{
		"{timestamp}": now.Format("20060102150405"),
		"{date}":      now.Format("20060102"),
		"{time}":      now.Format("150405"),
		"{year}":      now.Format("2006"),
		"{month}":     now.Format("01"),
		"{day}":       now.Format("02"),
		"{hour}":      now.Format("15"),
		"{minute}":    now.Format("04"),
		"{second}":    now.Format("05"),
		"{name}":      name,
		"{ext}":       ext,
		"{n}":         fmt.Sprintf("%d", counter),
	}

	result := pattern
	for placeholder, value := range replacements {
		result = strings.ReplaceAll(result, placeholder, value)
	}

	// If result doesn't have an extension, add the original one
	if filepath.Ext(result) == "" && ext != "" {
		result += ext
	}

	return result
}

// applyAutoRenameToPath applies auto-rename to a file path based on user settings.
// Returns the new path with the renamed filename.
func applyAutoRenameToPath(filePath string, user *users.User) string {
	if !user.AutoRename {
		return filePath
	}

	dir := path.Dir(filePath)
	filename := path.Base(filePath)
	
	// Use atomic counter to ensure unique sequence numbers for batch uploads
	counter := int(atomic.AddUint64(&renameCounter, 1))
	
	newFilename := applyAutoRename(user.RenamePattern, filename, counter)
	
	if dir == "." || dir == "/" {
		return "/" + newFilename
	}
	return dir + "/" + newFilename
}

// applyAutoRenameToPathWithFs applies auto-rename and checks for existing files.
// If the generated filename already exists, it adds a suffix to make it unique.
func applyAutoRenameToPathWithFs(filePath string, user *users.User, fs afero.Fs) string {
	if !user.AutoRename {
		return filePath
	}

	dir := path.Dir(filePath)
	filename := path.Base(filePath)
	
	// Use atomic counter to ensure unique sequence numbers for batch uploads
	counter := int(atomic.AddUint64(&renameCounter, 1))
	
	newFilename := applyAutoRename(user.RenamePattern, filename, counter)
	
	// Check if file exists and add suffix if needed
	ext := filepath.Ext(newFilename)
	nameWithoutExt := strings.TrimSuffix(newFilename, ext)
	
	finalPath := filepath.Join(dir, newFilename)
	if dir == "." || dir == "/" {
		finalPath = "/" + newFilename
	}
	
	// If file exists, try adding suffixes until we find a unique name
	suffix := 1
	for {
		exists, _ := afero.Exists(fs, finalPath)
		if !exists {
			break
		}
		suffix++
		newFilename = fmt.Sprintf("%s_%d%s", nameWithoutExt, suffix, ext)
		if dir == "." || dir == "/" {
			finalPath = "/" + newFilename
		} else {
			finalPath = filepath.Join(dir, newFilename)
		}
	}
	
	return finalPath
}
