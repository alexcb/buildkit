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
	//runc             *runc.Runc
	//root             string
	//cgroupParent     string
	//rootless         bool
	//networkProviders map[pb.NetMode]network.Provider
	//processMode      oci.ProcessMode
	//idmap            *idtools.IdentityMapping
	//noPivot          bool
	//dns              *oci.DNSConfig
	//oomScoreAdj      *int
	//running          map[string]chan error
	//mu               sync.Mutex
	sm *session.Manager
	g  session.Group
}

func New() (executor.Executor, error) {
	fmt.Printf("creating new LocalhostExecutor\n")
	return &LocalhostExecutor{}, nil
	//cmds := opt.CommandCandidates
	//if cmds == nil {
	//	cmds = defaultCommandCandidates
	//}

	//var cmd string
	//var found bool
	//for _, cmd = range cmds {
	//	if _, err := exec.LookPath(cmd); err == nil {
	//		found = true
	//		break
	//	}
	//}
	//if !found {
	//	return nil, errors.Errorf("failed to find %s binary", cmd)
	//}

	//root := opt.Root

	//if err := os.MkdirAll(root, 0711); err != nil {
	//	return nil, errors.Wrapf(err, "failed to create %s", root)
	//}

	//root, err := filepath.Abs(root)
	//if err != nil {
	//	return nil, err
	//}
	//root, err = filepath.EvalSymlinks(root)
	//if err != nil {
	//	return nil, err
	//}

	//// clean up old hosts/resolv.conf file. ignore errors
	//os.RemoveAll(filepath.Join(root, "hosts"))
	//os.RemoveAll(filepath.Join(root, "resolv.conf"))

	//runtime := &runc.Runc{
	//	Command:   cmd,
	//	Log:       filepath.Join(root, "runc-log.json"),
	//	LogFormat: runc.JSON,
	//	Setpgid:   true,
	//	// we don't execute runc with --rootless=(true|false) explicitly,
	//	// so as to support non-runc runtimes
	//}

	//updateRuncFieldsForHostOS(runtime)

	//w := &runcExecutor{
	//	runc:             runtime,
	//	root:             root,
	//	cgroupParent:     opt.DefaultCgroupParent,
	//	rootless:         opt.Rootless,
	//	networkProviders: networkProviders,
	//	processMode:      opt.ProcessMode,
	//	idmap:            opt.IdentityMapping,
	//	noPivot:          opt.NoPivot,
	//	dns:              opt.DNS,
	//	oomScoreAdj:      opt.OOMScoreAdj,
	//	running:          make(map[string]chan error),
	//}
	//return w, nil
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
	meta := process.Meta
	fmt.Printf("entered LocalhostExecutor.Exec with %v\n", meta)
	process.Stdout.Write([]byte(fmt.Sprintf("TODO: perform the exec call for %v\n", meta)))
	process.Stdout.Close()
	return nil
}

func (w *LocalhostExecutor) SetSessionManager(sm *session.Manager) {
	fmt.Printf("set the session manager to %v from \n%s\n", sm, debug.Stack())
	w.sm = sm
}

func (w *LocalhostExecutor) SetSessionGroup(g session.Group) {
	fmt.Printf("set the group to %v from \n%s\n", g, debug.Stack())
	w.g = g
}
