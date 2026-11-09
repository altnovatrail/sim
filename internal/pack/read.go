// Package pack reads published packages of modules as the bench needs them: the server part,
// server.json, by the types of spec's canonical view and without module-kit, checked through the
// chain of digests from a pinned module index (docs/architecture.md, sections 1, 3.3 and 3.9; R62,
// R72, R77 spec; the text of the format, sections 25 and 26). It translates the canonical view into
// what learner-model and engine take, as api would (R14 in docs/decisions.md).
package pack

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strconv"

	"github.com/altnovatrail/spec/types/formatv0"
)

// formatVersions are the versions of the module format the bench reads: while the format is v0,
// only the current one, 0.3 since spec v0.9.0 (R44, R101 spec), the window module-kit reads too.
var formatVersions = []formatv0.ModuleFormatVersion{"0.3"}

// ErrDigest is a file whose SHA-256 is not the one the chain of digests names.
var ErrDigest = errors.New("pack: a digest does not match")

// ErrFormat is a package of a version of the format the bench does not read.
var ErrFormat = errors.New("pack: a version of the module format outside the window")

// ErrWithdrawn is a version the index marks withdrawn from publication: api begins no attempt on
// it (R63 spec), and the bench reads none.
var ErrWithdrawn = errors.New("pack: the version is withdrawn from publication")

// ErrIdentity is a package that is not the module and version its entry names.
var ErrIdentity = errors.New("pack: the package is not the version its entry names")

// ReadIndex reads a module index whose SHA-256, in hex, must be digest: the digest the gate file
// pins, as the configuration of api pins it (R63 spec; docs/architecture.md, section 3.9). It reads
// the index's members by their places alone, whatever version of the format the index is of: a
// reader skips an entry whose format version it does not read, rather than refuse the index (the
// text of the format, section 2), and Entry does that.
func ReadIndex(data []byte, digest string) (*formatv0.ModuleIndex, error) {
	if hexDigest(data) != digest {
		return nil, fmt.Errorf("%w: the module index is %s, not %s", ErrDigest, hexDigest(data), digest)
	}
	var idx formatv0.ModuleIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, fmt.Errorf("pack: the module index: %w", err)
	}
	return &idx, nil
}

// Entry gives the entry of a version of a module in an index, if the bench may read it: one the
// index names, of a version of the format in the window, and not withdrawn.
func Entry(idx *formatv0.ModuleIndex, id formatv0.ModuleID, version formatv0.ModuleVersion) (formatv0.Entry, error) {
	e, ok := idx.Modules[id][version]
	if !ok {
		return e, fmt.Errorf("pack: the index names no %s %s", id, version)
	}
	if !Readable(e.ModuleFormatVersion) {
		return e, fmt.Errorf("%w: %s %s is %q", ErrFormat, id, version, e.ModuleFormatVersion)
	}
	if e.Withdrawn != nil {
		return e, fmt.Errorf("%w: %s %s, %s", ErrWithdrawn, id, version, *e.Withdrawn)
	}
	return e, nil
}

// ReadPackage reads the server part of the archive of a version of a module whose entry in the
// index is e: the SHA-256 of its inventory.json must be the entry's package digest, and the digest
// and size of server.json the inventory's; the inventory and the manifest must name the module,
// the version and the version of the format the entry names; server.json is read by the types of
// the canonical view (the text of the format, sections 25 and 26).
func ReadPackage(archive []byte, id formatv0.ModuleID, version formatv0.ModuleVersion, e formatv0.Entry) (*formatv0.CanonicalView, error) {
	files, err := Untar(archive)
	if err != nil {
		return nil, err
	}
	inventory, server := files["inventory.json"], files["server.json"]
	if inventory == nil || server == nil {
		return nil, errors.New("pack: the archive has no inventory.json or no server.json")
	}
	cv, inv, err := ReadFiles(inventory, server, &e)
	if err != nil {
		return nil, err
	}
	if inv.ID != id || inv.Version != version {
		return nil, fmt.Errorf("%w: the inventory is of %s %s; the entry names %s %s", ErrIdentity, inv.ID, inv.Version, id, version)
	}
	return cv, nil
}

// ReadFiles reads a package from its inventory.json and server.json, as a folder holds them: the
// digest and size of server.json must be the inventory's, whatever the version of the format, since
// the inventory keeps server_part in every version (the text of the format, section 2); that
// version must be one the bench reads, and the manifest must name the module, the version and the
// version of the format the inventory names; server.json is read by the types of the canonical
// view. With an entry of an index, the SHA-256 of inventory.json must be the entry's package digest
// and its version of the format the entry's: the chain of digests from the index. With none, the
// package is one a build wrote, which no index names yet. The inventory comes back whenever it could
// be read, so that a caller can tell a package of a version of the format outside the window.
func ReadFiles(inventory, server []byte, e *formatv0.Entry) (*formatv0.CanonicalView, *formatv0.Inventory, error) {
	if e != nil && Digest(inventory) != e.PackageDigest {
		return nil, nil, fmt.Errorf("%w: inventory.json is %s, the index names %s", ErrDigest, Digest(inventory), e.PackageDigest)
	}
	var inv formatv0.Inventory
	if err := json.Unmarshal(inventory, &inv); err != nil {
		return nil, nil, fmt.Errorf("pack: inventory.json: %w", err)
	}
	if e != nil && inv.ModuleFormatVersion != e.ModuleFormatVersion {
		return nil, &inv, fmt.Errorf("%w: the inventory is of format %s; the entry names format %s", ErrIdentity, inv.ModuleFormatVersion, e.ModuleFormatVersion)
	}
	if got := Digest(server); got != inv.ServerPart.Digest || int64(len(server)) != inv.ServerPart.Size {
		return nil, &inv, fmt.Errorf("%w: server.json is %s of %d bytes, the inventory names %s of %d", ErrDigest, got, len(server), inv.ServerPart.Digest, inv.ServerPart.Size)
	}
	if !Readable(inv.ModuleFormatVersion) {
		return nil, &inv, fmt.Errorf("%w: %s %s is %q", ErrFormat, inv.ID, inv.Version, inv.ModuleFormatVersion)
	}
	var cv formatv0.CanonicalView
	if err := json.Unmarshal(server, &cv); err != nil {
		return nil, &inv, fmt.Errorf("pack: server.json: %w", err)
	}
	if cv.Manifest.ID != inv.ID || cv.Manifest.Version != inv.Version || cv.Manifest.ModuleFormatVersion != inv.ModuleFormatVersion {
		return nil, &inv, fmt.Errorf("%w: the manifest is of %s %s, format %s", ErrIdentity, cv.Manifest.ID, cv.Manifest.Version, cv.Manifest.ModuleFormatVersion)
	}
	return &cv, &inv, nil
}

// Digest gives the SHA-256 of a file as the format writes a digest: the package digest, when the
// file is a package's inventory.json.
func Digest(data []byte) formatv0.Digest {
	return formatv0.Digest("sha256:" + hexDigest(data))
}

// Readable tells whether the bench reads a version of the module format.
func Readable(v formatv0.ModuleFormatVersion) bool {
	return slices.Contains(formatVersions, v)
}

func hexDigest(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// Untar gives the regular files of a package's archive, by name. The archive is a gzip of a tar of
// 512-byte blocks: a header — the name in its first 100 bytes, the size in octal at 124, the type
// at 156 — then the file's blocks; the packages name their files short, with no ustar prefix. It
// is read here, as engine's test reads it, since archive/tar would pull os/user, whose cgo the
// rules of code cannot type-check for every target.
func Untar(archive []byte) (map[string][]byte, error) {
	zr, err := gzip.NewReader(bytes.NewReader(archive))
	if err != nil {
		return nil, fmt.Errorf("pack: the archive: %w", err)
	}
	var b bytes.Buffer
	if _, err := b.ReadFrom(zr); err != nil {
		return nil, fmt.Errorf("pack: the archive: %w", err)
	}
	data := b.Bytes()
	files := map[string][]byte{}
	for len(data) >= 512 && data[0] != 0 {
		h := data[:512]
		name := string(bytes.TrimRight(h[:100], "\x00"))
		size, err := strconv.ParseInt(string(bytes.Trim(h[124:136], " \x00")), 8, 64)
		if err != nil || size < 0 {
			return nil, fmt.Errorf("pack: a broken header at %q", name)
		}
		next := 512 + (size+511)/512*512
		if int64(len(data)) < next {
			return nil, fmt.Errorf("pack: a short block at %q", name)
		}
		if h[156] == '0' || h[156] == 0 {
			files[name] = data[512 : 512+size]
		}
		data = data[next:]
	}
	return files, nil
}
