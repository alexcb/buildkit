package localhost

import (
	"context"
	//"encoding/json"
	"fmt"
	//"time"
	"runtime/debug"

	"github.com/docker/docker/pkg/idtools"
	"github.com/moby/buildkit/cache"
	"github.com/moby/buildkit/cache/metadata"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/snapshot"
	"github.com/moby/buildkit/solver"
	"github.com/moby/buildkit/source"
	"github.com/moby/buildkit/util/compression"
	"github.com/pkg/errors"
	//"github.com/sirupsen/logrus"
	//"github.com/tonistiigi/fsutil"
	//fstypes "github.com/tonistiigi/fsutil/types"
	//bolt "go.etcd.io/bbolt"
	//"google.golang.org/grpc/codes"
	//"google.golang.org/grpc/status"
)

type Opt struct {
	CacheAccessor cache.Accessor
	MetadataStore *metadata.Store
}

func NewSource(opt Opt) (source.Source, error) {
	ls := &localHostSource{
		cm: opt.CacheAccessor,
		md: opt.MetadataStore,
	}
	return ls, nil
}

type localHostSource struct {
	cm cache.Accessor
	md *metadata.Store
}

func (ls *localHostSource) ID() string {
	return source.LocalHostScheme
}

func (ls *localHostSource) Resolve(ctx context.Context, id source.Identifier, sm *session.Manager, _ solver.Vertex) (source.SourceInstance, error) {
	//panic("unhandled localhost resolve")
	localHostIdentifier, ok := id.(*source.LocalHostIdentifier)
	if !ok {
		return nil, errors.Errorf("invalid localhost identifier %v", id)
	}

	return &localHostSourceHandler{
		src:             *localHostIdentifier,
		sm:              sm,
		localHostSource: ls,
	}, nil
}

type localHostSourceHandler struct {
	src source.LocalHostIdentifier
	sm  *session.Manager
	*localHostSource
}

func (ls *localHostSourceHandler) CacheKey(ctx context.Context, g session.Group, index int) (string, solver.CacheOpts, bool, error) {
	// TODO always return not cached
	//panic("unimplemented localhost CacheKey")
	//sessionID := ls.src.SessionID

	//if sessionID == "" {
	//	id := g.SessionIterator().NextSession()
	//	if id == "" {
	//		return "", nil, false, errors.New("could not access local files without session")
	//	}
	//	sessionID = id
	//}
	//dt, err := json.Marshal(struct {
	//	SessionID       string
	//	IncludePatterns []string
	//	ExcludePatterns []string
	//	FollowPaths     []string
	//}{SessionID: sessionID, IncludePatterns: ls.src.IncludePatterns, ExcludePatterns: ls.src.ExcludePatterns, FollowPaths: ls.src.FollowPaths})
	//if err != nil {
	//	return "", nil, false, err
	//}
	//return "session:" + ls.src.Name + ":" + digest.FromBytes(dt).String(), nil, true, nil

	return "localhost-do-not-cache", nil, true, nil // TODO inject this somewhere to prevent caching
}

func (ls *localHostSourceHandler) Snapshot(ctx context.Context, g session.Group) (cache.ImmutableRef, error) {
	//panic("unimplemented snapshot")
	//var ref cache.ImmutableRef
	//err := ls.sm.Any(ctx, g, func(ctx context.Context, _ string, c session.Caller) error {
	//	r, err := ls.snapshot(ctx, g, c)
	//	if err != nil {
	//		return err
	//	}
	//	ref = r
	//	return nil
	//})
	//if err != nil {
	//	return nil, err
	//}

	_ = debug.Stack

	fmt.Printf("returning nil snapshot\n%s\n", debug.Stack())
	return nil, nil

	//ref := &LocalHostRef{
	//	id: 0,
	//}
	//fmt.Printf("creating fake snapshot under %v\n", ref)
	//return ref, nil
}

type LocalHostRef struct {
	id int
}

func (lhr *LocalHostRef) ID() string                                      { panic("not implemented") }
func (lhr *LocalHostRef) Release(context.Context) error                   { panic("not implemented") }
func (lhr *LocalHostRef) Size(ctx context.Context) (int64, error)         { panic("not implemented") }
func (lhr *LocalHostRef) Metadata() *metadata.StorageItem                 { panic("not implemented") }
func (lhr *LocalHostRef) IdentityMapping() *idtools.IdentityMapping       { panic("not implemented") }
func (lhr *LocalHostRef) Parent() cache.ImmutableRef                      { panic("not implemented") }
func (lhr *LocalHostRef) Finalize(ctx context.Context, commit bool) error { panic("not implemented") }
func (lhr *LocalHostRef) Clone() cache.ImmutableRef                       { panic("not implemented") }
func (lhr *LocalHostRef) Info() cache.RefInfo                             { panic("not implemented") }
func (lhr *LocalHostRef) Extract(ctx context.Context, s session.Group) error {
	panic("not implemented")
}
func (lhr *LocalHostRef) GetRemote(ctx context.Context, createIfNeeded bool, compressionType compression.Type, s session.Group) (*solver.Remote, error) {
	panic("not implemented")
}
func (lhr *LocalHostRef) Mount(ctx context.Context, readonly bool, s session.Group) (snapshot.Mountable, error) {
	panic("not implemented")
}
