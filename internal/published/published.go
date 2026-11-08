// Package published holds the published packages of modules the bench is pinned to, as their
// release has them, and reads them through the chain of digests (docs/architecture.md, sections
// 3.3 and 3.9; R14 in docs/decisions.md). The files came in only after release-verify accepted
// each of them (RULES.md, Hard rules; R77 spec): from then on the digest of the module index,
// pinned here, holds them, through the package digests it names and the server parts' digests
// their inventories name. A new release of modules is a new pin: a pull request of sim with new
// bands (docs/architecture.md, section 3.7).
package published

import (
	"embed"
	"fmt"
	"maps"
	"slices"
	"strings"
	"sync"

	"github.com/altnovatrail/sim/internal/pack"
	"github.com/altnovatrail/spec/types/formatv0"
)

// The release of modules the bench is pinned to and the SHA-256 of its module index. The archives
// lie each in the folder of the release its entry names, as files of that release.
const (
	Release     = "v0.1.0"
	IndexDigest = "ed4d8bb325fdcb3ea54da439e08c58d927f773ba0f097338f7767ad76be8c90e"
)

//go:embed v0.1.0
var files embed.FS

// versions are the modules of the first slice (R41 spec), by full ID, and the version of each the
// bench reads. They are a pin, as the release is: unexported, so nothing outside changes them.
var versions = map[formatv0.ModuleID]formatv0.ModuleVersion{
	"novatrail/olympiad-math":  "0.1.0",
	"novatrail/english-essays": "0.1.0",
	"novatrail/grade-3-math":   "0.1.0",
}

// read reads the pinned packages once: the bytes are fixed, and so is what they give.
var read = sync.OnceValues(func() (map[string]*pack.Module, error) {
	return readPin(Release, IndexDigest, versions)
})

// Modules gives the pinned packages as the bench's modules, by full ID, read once through every
// digest of the chain from the pinned index. The modules are shared, and read only.
func Modules() (map[string]*pack.Module, error) {
	return read()
}

// index reads the index of a release, held to its digest.
func index(release, indexDigest string) (*formatv0.ModuleIndex, error) {
	data, err := files.ReadFile(release + "/module-index.json")
	if err != nil {
		return nil, err
	}
	return pack.ReadIndex(data, indexDigest)
}

func readPin(release, indexDigest string, versions map[formatv0.ModuleID]formatv0.ModuleVersion) (map[string]*pack.Module, error) {
	idx, err := index(release, indexDigest)
	if err != nil {
		return nil, err
	}
	out := map[string]*pack.Module{}
	for _, id := range slices.Sorted(maps.Keys(versions)) {
		e, err := pack.Entry(idx, id, versions[id])
		if err != nil {
			return nil, fmt.Errorf("published: %w", err)
		}
		hex, ok := strings.CutPrefix(string(e.PackageDigest), "sha256:")
		if !ok {
			return nil, fmt.Errorf("published: %s has the package digest %q, not sha256:", id, e.PackageDigest)
		}
		archive, err := files.ReadFile(string(e.Release) + "/sha256-" + hex + ".tar.gz")
		if err != nil {
			return nil, err
		}
		cv, err := pack.ReadPackage(archive, id, versions[id], e)
		if err != nil {
			return nil, fmt.Errorf("published: %s: %w", id, err)
		}
		m, err := pack.Translate(cv)
		if err != nil {
			return nil, err
		}
		out[string(id)] = m
	}
	return out, nil
}

// Package is a pinned package: its module's full ID, its version, and its package digest, as the
// pinned index names them.
type Package struct {
	Module  string `json:"module"`
	Version string `json:"version"`
	Digest  string `json:"digest"`
}

// Pinned gives the pinned packages in the order of their IDs, read from the pinned index, so that a
// gate file that names the bench's packages is held to the bytes the bench reads.
func Pinned() ([]Package, error) {
	idx, err := index(Release, IndexDigest)
	if err != nil {
		return nil, err
	}
	var out []Package
	for _, id := range slices.Sorted(maps.Keys(versions)) {
		e, err := pack.Entry(idx, id, versions[id])
		if err != nil {
			return nil, fmt.Errorf("published: %w", err)
		}
		out = append(out, Package{Module: string(id), Version: string(versions[id]), Digest: string(e.PackageDigest)})
	}
	return out, nil
}
