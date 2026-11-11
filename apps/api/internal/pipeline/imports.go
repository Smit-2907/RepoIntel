package pipeline

import (
	"context"
	"path"
	"strings"
)

// resolveImports links import specifiers to in-repo files (heuristic, per design v1.1).
func (p *Pipeline) resolveImports(ctx context.Context, analysisID string) error {
	files, err := p.Store.ListFiles(ctx, analysisID)
	if err != nil {
		return err
	}
	byPath := make(map[string]string, len(files))
	paths := make([]string, 0, len(files))
	for _, f := range files {
		byPath[f.Path] = f.ID
		paths = append(paths, f.Path)
	}
	edges, err := p.Store.ListImportEdges(ctx, analysisID)
	if err != nil {
		return err
	}
	for _, e := range edges {
		target, kind := resolveSpecifier(e.Raw, e.FromPath, byPath, paths)
		var toID *string
		if target != "" {
			id := byPath[target]
			toID = &id
		}
		if err := p.Store.SetImportResolution(ctx, e.ID, toID, kind); err != nil {
			return err
		}
	}
	return nil
}

func resolveSpecifier(raw, fromPath string, byPath map[string]string, paths []string) (string, string) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", "unknown"
	}

	// Relative imports (JS/TS/Python)
	if strings.HasPrefix(raw, ".") {
		base := path.Join(path.Dir(fromPath), raw)
		for _, ext := range []string{".ts", ".tsx", ".js", ".jsx", ".mjs", ".cjs", ".py"} {
			if _, ok := byPath[base+ext]; ok {
				return base + ext, "internal"
			}
		}
		for _, ext := range []string{".ts", ".tsx", ".js", ".jsx", ".py"} {
			for _, suffix := range []string{"/index" + ext, "/__init__" + ext} {
				if _, ok := byPath[base+suffix]; ok {
					return base + suffix, "internal"
				}
			}
		}
		if _, ok := byPath[base]; ok {
			return base, "internal"
		}
	}

	// Python dotted module (absolute)
	if !strings.Contains(raw, "/") && strings.Contains(raw, ".") {
		mod := strings.ReplaceAll(raw, ".", "/")
		for _, cand := range []string{mod + ".py", mod + "/__init__.py"} {
			if _, ok := byPath[cand]; ok {
				return cand, "internal"
			}
		}
	}

	// Go package path: match by trailing directory name
	segs := strings.Split(raw, "/")
	last := segs[len(segs)-1]
	if last != "" {
		for _, p := range paths {
			if path.Base(path.Dir(p)) == last {
				return p, "internal"
			}
		}
	}

	// Fallback: suffix match
	for _, p := range paths {
		if strings.HasSuffix(p, raw) {
			return p, "internal"
		}
	}

	// Unresolved: domain-like or path-like specifiers are external packages.
	if strings.Contains(raw, "/") || strings.Contains(raw, ".") {
		return "", "external"
	}
	return "", "unknown"
}
