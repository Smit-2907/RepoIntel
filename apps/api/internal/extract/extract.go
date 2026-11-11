// Package extract safely unpacks repository tarballs (design v1.1 §3 step 4, §11).
// It never executes repository code; it only reads and filters files.
package extract

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"path"
	"strings"
)

type Limits struct {
	MaxFiles      int
	MaxFileBytes  int
	MaxTotalBytes int64
}

type File struct {
	Path     string
	Language string
	Content  []byte
}

var excludedDirs = map[string]bool{
	"node_modules": true, "vendor": true, "dist": true, "build": true, ".git": true,
	"target": true, "__pycache__": true, ".venv": true, "venv": true, "coverage": true, ".next": true,
}

var excludedFiles = map[string]bool{
	"package-lock.json": true, "yarn.lock": true, "pnpm-lock.yaml": true, "go.sum": true,
	"poetry.lock": true, "Pipfile.lock": true, "Cargo.lock": true,
}

func LanguageFor(p string) string {
	switch strings.ToLower(path.Ext(p)) {
	case ".py":
		return "python"
	case ".go":
		return "go"
	case ".js", ".jsx", ".mjs", ".cjs":
		return "javascript"
	case ".ts", ".tsx", ".mts", ".cts":
		return "typescript"
	}
	return "other"
}

func isExcluded(p string) bool {
	parts := strings.Split(p, "/")
	for _, part := range parts[:len(parts)-1] {
		if excludedDirs[part] {
			return true
		}
	}
	base := path.Base(p)
	if excludedFiles[base] {
		return true
	}
	return strings.HasSuffix(base, ".min.js") || strings.HasSuffix(base, ".min.css")
}

func isBinary(b []byte) bool {
	n := len(b)
	if n > 512 {
		n = 512
	}
	return bytes.IndexByte(b[:n], 0) >= 0
}

// Tarball filters an archive into indexable source files, recording skip reasons.
func Tarball(data []byte, lim Limits) ([]File, map[string]int, error) {
	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, nil, err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	files := make([]File, 0)
	skipped := map[string]int{}
	var total int64

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}
		name := path.Clean(strings.TrimPrefix(hdr.Name, "./"))
		if hdr.Typeflag == tar.TypeDir {
			continue
		}
		if hdr.Typeflag != tar.TypeReg {
			skipped["non_regular"]++
			continue
		}
		if strings.HasPrefix(name, "/") || strings.Contains(name, "..") {
			skipped["unsafe_path"]++
			continue
		}
		// GitHub archives nest everything under "{owner}-{repo}-{sha}/"; strip it.
		i := strings.IndexByte(name, '/')
		if i <= 0 {
			continue
		}
		name = name[i+1:]
		if isExcluded(name) {
			skipped["excluded"]++
			continue
		}
		if len(files) >= lim.MaxFiles {
			skipped["max_files"]++
			continue
		}
		if hdr.Size > int64(lim.MaxFileBytes) {
			skipped["too_large"]++
			continue
		}
		content, err := io.ReadAll(io.LimitReader(tr, int64(lim.MaxFileBytes)+1))
		if err != nil {
			return nil, nil, err
		}
		if len(content) > lim.MaxFileBytes {
			skipped["too_large"]++
			continue
		}
		if isBinary(content) {
			skipped["binary"]++
			continue
		}
		lang := LanguageFor(name)
		if lang == "other" {
			skipped["unsupported"]++
			continue
		}
		total += int64(len(content))
		if total > lim.MaxTotalBytes {
			return nil, nil, fmt.Errorf("extracted content exceeds %d bytes", lim.MaxTotalBytes)
		}
		files = append(files, File{Path: name, Language: lang, Content: content})
	}
	return files, skipped, nil
}
