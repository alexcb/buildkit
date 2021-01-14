package localhost

import (
	"context"
	"fmt"
	"io"

	"github.com/moby/buildkit/executor"
	"github.com/moby/buildkit/session"
)

// LocalhostExec is called by buildkitd; it connects to the user's client to request the client execute a command localy.
func LocalhostExec(ctx context.Context, c session.Caller, process executor.ProcessInfo) error {
	defer process.Stdout.Close()
	defer process.Stderr.Close()

	client := NewLocalhostClient(c.Conn())
	stream, err := client.Exec(ctx)
	if err != nil {
		return err
	}

	req := InputMessage{
		Command: process.Meta.Args,
	}
	if err := stream.SendMsg(&req); err != nil {
		return err
	}
	fmt.Printf("sent command, now waiting for response\n")

	var exitCodeSet bool
	var exitCode int
	for {
		var msg OutputMessage
		err := stream.RecvMsg(&msg)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("got EOF\n")
				break
			}
			return err
		}
		process.Stdout.Write(msg.Stdout)
		process.Stderr.Write(msg.Stderr)
		switch msg.Status {
		case RUNNING:
			break
		case DONE:
			if exitCodeSet {
				panic("received multiple DONE messages (shouldn't happen")
			}
			exitCode = int(msg.ExitCode)
		default:
			return fmt.Errorf("unhandled exit status: %d", msg.Status)
		}
	}

	if exitCode != 0 {
		return fmt.Errorf("exit code: %d", exitCode)
	}

	return nil
}

//// DefaultID is the default ssh ID
//const DefaultID = "default"
//
//const KeySSHID = "buildkit.ssh.id"
//
//type server struct {
//	caller session.Caller
//}
//
//func (s *server) run(ctx context.Context, l net.Listener, id string) error {
//	eg, ctx := errgroup.WithContext(ctx)
//
//	eg.Go(func() error {
//		<-ctx.Done()
//		return ctx.Err()
//	})
//
//	eg.Go(func() error {
//		for {
//			conn, err := l.Accept()
//			if err != nil {
//				return err
//			}
//
//			client := NewSSHClient(s.caller.Conn())
//
//			opts := make(map[string][]string)
//			opts[KeySSHID] = []string{id}
//			ctx = metadata.NewOutgoingContext(ctx, opts)
//
//			stream, err := client.ForwardAgent(ctx)
//			if err != nil {
//				conn.Close()
//				return err
//			}
//
//			go Copy(ctx, conn, stream, stream.CloseSend)
//		}
//	})
//
//	return eg.Wait()
//}
//
//type SocketOpt struct {
//	ID   string
//	UID  int
//	GID  int
//	Mode int
//}
//
//func MountSSHSocket(ctx context.Context, c session.Caller, opt SocketOpt) (sockPath string, closer func() error, err error) {
//	dir, err := ioutil.TempDir("", ".buildkit-ssh-sock")
//	if err != nil {
//		return "", nil, errors.WithStack(err)
//	}
//
//	defer func() {
//		if err != nil {
//			os.RemoveAll(dir)
//		}
//	}()
//
//	if err := os.Chmod(dir, 0711); err != nil {
//		return "", nil, errors.WithStack(err)
//	}
//
//	sockPath = filepath.Join(dir, "ssh_auth_sock")
//
//	l, err := net.Listen("unix", sockPath)
//	if err != nil {
//		return "", nil, errors.WithStack(err)
//	}
//
//	if err := os.Chown(sockPath, opt.UID, opt.GID); err != nil {
//		l.Close()
//		return "", nil, errors.WithStack(err)
//	}
//	if err := os.Chmod(sockPath, os.FileMode(opt.Mode)); err != nil {
//		l.Close()
//		return "", nil, errors.WithStack(err)
//	}
//
//	s := &server{caller: c}
//
//	id := opt.ID
//	if id == "" {
//		id = DefaultID
//	}
//
//	go s.run(ctx, l, id) // erroring per connection allowed
//
//	return sockPath, func() error {
//		err := l.Close()
//		os.RemoveAll(sockPath)
//		return errors.WithStack(err)
//	}, nil
//}
//
//func CheckSSHID(ctx context.Context, c session.Caller, id string) error {
//	client := NewSSHClient(c.Conn())
//	_, err := client.CheckAgent(ctx, &CheckAgentRequest{ID: id})
//	return errors.WithStack(err)
//}
