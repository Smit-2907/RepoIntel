package httpapi

import (
	"path"
	"regexp"
	"strconv"
	"strings"

	"codebaseintel/internal/store"
)

type stackFrame struct {
	File string
	Line int
	Func string
}

var (
	rePyFrame = regexp.MustCompile(`File "([^"]+)", line (\d+)(?:, in ([A-Za-z0-9_<>.]+))?`)
	reJSFrame = regexp.MustCompile(`at (?:([^\s(]+) )?\(?([^\s():]+):(\d+):\d+\)?`)
	reGoFrame = regexp.MustCompile(`([A-Za-z0-9_./-]+\.go):(\d+)`)
)

func parseStackFrames(trace string) []stackFrame {
	seen := map[string]bool{}
	out := make([]stackFrame, 0, 8)
	add := func(f stackFrame) {
		if f.File == "" {
			return
		}
		key := f.File + ":" + strconv.Itoa(f.Line)
		if seen[key] {
			return
		}
		seen[key] = true
		out = append(out, f)
	}
	for _, m := range rePyFrame.FindAllStringSubmatch(trace, -1) {
		line, _ := strconv.Atoi(m[2])
		add(stackFrame{File: m[1], Line: line, Func: m[3]})
	}
	for _, m := range reJSFrame.FindAllStringSubmatch(trace, -1) {
		line, _ := strconv.Atoi(m[3])
		add(stackFrame{File: m[2], Line: line, Func: m[1]})
	}
	for _, m := range reGoFrame.FindAllStringSubmatch(trace, -1) {
		line, _ := strconv.Atoi(m[2])
		add(stackFrame{File: m[1], Line: line})
	}
	return out
}

func matchFile(files []store.FileRef, want string) (string, bool) {
	want = strings.TrimPrefix(want, "./")
	for _, f := range files {
		if f.Path == want {
			return f.Path, true
		}
	}
	base := path.Base(want)
	for _, f := range files {
		if path.Base(f.Path) == base {
			return f.Path, true
		}
	}
	for _, f := range files {
		if strings.HasSuffix(f.Path, "/"+want) || strings.HasSuffix(f.Path, want) {
			return f.Path, true
		}
	}
	return "", false
}

func dedupeChunks(in []store.RetrievedChunk, limit int) []store.RetrievedChunk {
	seen := map[string]bool{}
	out := make([]store.RetrievedChunk, 0, limit)
	for _, c := range in {
		if seen[c.ChunkID] {
			continue
		}
		seen[c.ChunkID] = true
		out = append(out, c)
		if len(out) >= limit {
			break
		}
	}
	return out
}
