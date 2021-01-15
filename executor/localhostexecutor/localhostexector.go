package localhostexecutor

import (
	"context"
	//"encoding/json"
	"fmt"
	//"io"
	"os"
	//"os/exec"
	//"path/filepath"
	"runtime/debug"
	//"strings"
	//"sync"
	//"syscall"
	//"time"

	"github.com/moby/buildkit/session"

	//"github.com/containerd/containerd/mount"
	//containerdoci "github.com/containerd/containerd/oci"
	//"github.com/containerd/continuity/fs"
	//runc "github.com/containerd/go-runc"
	"github.com/docker/docker/pkg/idtools"
	"github.com/moby/buildkit/executor"
	"github.com/moby/buildkit/executor/oci"
	"github.com/moby/buildkit/session/localhost"
	//"github.com/moby/buildkit/frontend/gateway/errdefs"
	//"github.com/moby/buildkit/identity"
	//"github.com/moby/buildkit/solver/pb"
	//"github.com/moby/buildkit/util/network"
	//rootlessspecconv "github.com/moby/buildkit/util/rootless/specconv"
	//"github.com/moby/buildkit/util/stack"
	//"github.com/opencontainers/runtime-spec/specs-go"
	//"github.com/pkg/errors"
	//"github.com/sirupsen/logrus"
)

type Opt struct {
	// root directory
	Root              string
	CommandCandidates []string
	// without root privileges (has nothing to do with Opt.Root directory)
	Rootless bool
	// DefaultCgroupParent is the cgroup-parent name for executor
	DefaultCgroupParent string
	// ProcessMode
	ProcessMode     oci.ProcessMode
	IdentityMapping *idtools.IdentityMapping
	// runc run --no-pivot (unrecommended)
	NoPivot     bool
	DNS         *oci.DNSConfig
	OOMScoreAdj *int
}

var defaultCommandCandidates = []string{"buildkit-runc", "runc"}

type LocalhostExecutor struct {
	sm *session.Manager
	g  session.Group
}

func New() (executor.Executor, error) {
	fmt.Printf("creating new LocalhostExecutor\n")
	return &LocalhostExecutor{}, nil
}

func (w *LocalhostExecutor) Run(ctx context.Context, id string, root executor.Mount, mounts []executor.Mount, process executor.ProcessInfo, started chan<- struct{}) error {
	meta := process.Meta
	fmt.Printf("entered LocalhostExecutor.Run with %v; root: %v, mounts: %v; pid: %v\n%s\n", meta, root, mounts, os.Getpid(), debug.Stack())
	if len(mounts) > 0 {
		return fmt.Errorf("LocalhostExecutor does not support mounts")
	}

	err := w.sm.Any(ctx, w.g, func(ctx context.Context, _ string, caller session.Caller) error {
		err := localhost.LocalhostExec(ctx, caller, process)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (w *LocalhostExecutor) Exec(ctx context.Context, id string, process executor.ProcessInfo) (err error) {
	panic("LocalhostExecutor.Exec is not implemented")
}

func (w *LocalhostExecutor) SetSessionManager(sm *session.Manager) {
	fmt.Printf("set the session manager to %v from \n%s\n", sm, debug.Stack())
	w.sm = sm
}

func (w *LocalhostExecutor) SetSessionGroup(g session.Group) {
	fmt.Printf("set the group to %v from \n%s\n", g, debug.Stack())
	w.g = g
}
