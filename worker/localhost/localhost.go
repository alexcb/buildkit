package localhost

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/containerd/containerd/snapshots"
	//"github.com/containerd/containerd"
	"github.com/containerd/containerd/gc"
	//"github.com/containerd/containerd/leases"
	//"github.com/moby/buildkit/cache"
	"github.com/moby/buildkit/cache/metadata"
	"github.com/moby/buildkit/executor/localhostexecutor"
	//"github.com/moby/buildkit/executor/oci"
	"github.com/moby/buildkit/snapshot"
	//containerdsnapshot "github.com/moby/buildkit/snapshot/containerd"
	//"github.com/moby/buildkit/util/leaseutil"
	//"github.com/moby/buildkit/util/network/netproviders"
	//"github.com/moby/buildkit/util/winlayers"
	"github.com/moby/buildkit/worker/base"
	//specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/containerd/containerd/leases"
	"github.com/docker/docker/pkg/idtools"
	"github.com/pkg/errors"
)

// NewWorkerOpt creates a WorkerOpt.
func NewWorkerOpt(root, snapshotterName string) (base.WorkerOpt, error) {
	if strings.Contains(snapshotterName, "/") {
		return base.WorkerOpt{}, errors.Errorf("bad snapshotter name: %q", snapshotterName)
	}
	name := "localhost-" + snapshotterName
	root = filepath.Join(root, name)
	if err := os.MkdirAll(root, 0700); err != nil {
		return base.WorkerOpt{}, errors.Wrapf(err, "failed to create %s", root)
	}

	//df := client.DiffService()
	// TODO: should use containerd daemon instance ID (containerd/containerd#1862)?
	id, err := base.ID(root)
	if err != nil {
		return base.WorkerOpt{}, err
	}
	xlabels := base.Labels("localhost", snapshotterName)
	//for k, v := range labels {
	//	xlabels[k] = v
	//}

	//lm := leaseutil.WithNamespace(client.LeasesService(), "ns-todo")

	gc := func(ctx context.Context) (gc.Stats, error) {
		return nil, nil // TODO
	}
	//	l, err := lm.Create(ctx)
	//	if err != nil {
	//		return nil, nil
	//	}
	//	return nil, lm.Delete(ctx, leases.Lease{ID: l.ID}, leases.SynchronousDelete)
	//}

	//cs := containerdsnapshot.NewContentStore(client.ContentStore(), ns)

	//resp, err := client.IntrospectionService().Plugins(context.TODO(), []string{"type==io.containerd.runtime.v1", "type==io.containerd.runtime.v2"})
	//if err != nil {
	//	return base.WorkerOpt{}, errors.Wrap(err, "failed to list runtime plugin")
	//}
	//if len(resp.Plugins) == 0 {
	//	return base.WorkerOpt{}, errors.New("failed to find any runtime plugins")
	//}

	//var platforms []specs.Platform
	//for _, plugin := range resp.Plugins {
	//	for _, p := range plugin.Platforms {
	//		platforms = append(platforms, specs.Platform{
	//			OS:           p.OS,
	//			Architecture: p.Architecture,
	//			Variant:      p.Variant,
	//		})
	//	}
	//}

	//np, err := netproviders.Providers(nopt)
	//if err != nil {
	//	return base.WorkerOpt{}, err
	//}

	//snap := containerdsnapshot.NewSnapshotter(snapshotterName, client.SnapshotService(snapshotterName), ns, nil)

	//if err := cache.MigrateV2(context.TODO(), filepath.Join(root, "metadata.db"), filepath.Join(root, "metadata_v2.db"), cs, snap, lm); err != nil {
	//	return base.WorkerOpt{}, err
	//}

	md, err := metadata.NewStore(filepath.Join(root, "metadata_v2.db"))
	if err != nil {
		return base.WorkerOpt{}, err
	}

	exec, err := localhostexecutor.New()
	if err != nil {
		return base.WorkerOpt{}, err
	}

	opt := base.WorkerOpt{
		ID:            id,
		Labels:        xlabels,
		MetadataStore: md,
		Executor:      exec,
		Snapshotter:   newSnapStub(),
		//ContentStore: cs,
		//Applier:        winlayers.NewFileSystemApplierWithWindows(cs, df),
		//Differ:         winlayers.NewWalkingDiffWithWindows(cs, df),
		//ImageStore:     client.ImageService(),
		//Platforms:      platforms,
		LeaseManager:   NewStubLeaseManager(),
		GarbageCollect: gc,
	}
	return opt, nil
}

func NewStubLeaseManager() leases.Manager {
	return &sLM{}
}

type sLM struct {
}

func (l *sLM) Create(ctx context.Context, opts ...leases.Opt) (leases.Lease, error) {
	return leases.Lease{}, nil
}

func (l *sLM) Delete(ctx context.Context, lease leases.Lease, opts ...leases.DeleteOpt) error {
	return nil
}

func (l *sLM) List(ctx context.Context, filters ...string) ([]leases.Lease, error) {
	return nil, nil
}

func (l *sLM) AddResource(ctx context.Context, lease leases.Lease, resource leases.Resource) error {
	return nil
}

func (l *sLM) DeleteResource(ctx context.Context, lease leases.Lease, resource leases.Resource) error {
	return nil
}

func (l *sLM) ListResources(ctx context.Context, lease leases.Lease) ([]leases.Resource, error) {
	return nil, nil
}

type snapStub struct {
}

func newSnapStub() snapshot.Snapshotter {
	return &snapStub{}
}

func (s *snapStub) Name() string {
	return "snapstub"
}
func (s *snapStub) Mounts(ctx context.Context, key string) (snapshot.Mountable, error) {
	return nil, nil
}
func (s *snapStub) Prepare(ctx context.Context, key, parent string, opts ...snapshots.Opt) error {
	return nil
}
func (s *snapStub) View(ctx context.Context, key, parent string, opts ...snapshots.Opt) (snapshot.Mountable, error) {
	return nil, nil
}
func (s *snapStub) Stat(ctx context.Context, key string) (snapshots.Info, error) {
	return snapshots.Info{}, nil
}
func (s *snapStub) Update(ctx context.Context, info snapshots.Info, fieldpaths ...string) (snapshots.Info, error) {
	return snapshots.Info{}, nil
}
func (s *snapStub) Usage(ctx context.Context, key string) (snapshots.Usage, error) {
	return snapshots.Usage{}, nil
}
func (s *snapStub) Commit(ctx context.Context, name, key string, opts ...snapshots.Opt) error {
	return nil
}
func (s *snapStub) Remove(ctx context.Context, key string) error { return nil }
func (s *snapStub) Walk(ctx context.Context, fn snapshots.WalkFunc, filters ...string) error {
	return nil
}
func (s *snapStub) Close() error                              { return nil }
func (s *snapStub) IdentityMapping() *idtools.IdentityMapping { return nil }
