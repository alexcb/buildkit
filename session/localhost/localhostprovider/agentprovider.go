package localhostprovider

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/localhost"
	//"github.com/pkg/errors"
	//"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func NewLocalhostProvider() (session.Attachable, error) {
	//m := map[string]source{}
	//for _, conf := range confs {
	//	if len(conf.Paths) == 0 || len(conf.Paths) == 1 && conf.Paths[0] == "" {
	//		conf.Paths = []string{os.Getenv("SSH_AUTH_SOCK")}
	//	}

	//	if conf.Paths[0] == "" {
	//		return nil, errors.Errorf("invalid empty ssh agent socket, make sure SSH_AUTH_SOCK is set")
	//	}

	//	src, err := toAgentSource(conf.Paths)
	//	if err != nil {
	//		return nil, err
	//	}
	//	if conf.ID == "" {
	//		conf.ID = sshforward.DefaultID
	//	}
	//	if _, ok := m[conf.ID]; ok {
	//		return nil, errors.Errorf("invalid duplicate ID %s", conf.ID)
	//	}
	//	m[conf.ID] = src
	//}

	fmt.Printf("creating new localhost provider\n")
	return &localhostProvider{}, nil
}

type localhostProvider struct {
}

func (lp *localhostProvider) Register(server *grpc.Server) {
	fmt.Printf("registering localhost provider\n")
	localhost.RegisterLocalhostServer(server, lp)
}

func (lp *localhostProvider) Exec(stream localhost.Localhost_ExecServer) error {
	opts, _ := metadata.FromIncomingContext(stream.Context()) // if no metadata continue with empty object
	fmt.Printf("handling localhost exec call with opts: %v\n", opts)

	// first message must contain the command (and no stdin)
	var msg localhost.InputMessage
	err := stream.RecvMsg(&msg)
	if err != nil {
		return err
	}

	fmt.Printf("Running command: %v\n", msg.Command)

	if len(msg.Command) == 0 {
		return fmt.Errorf("command is empty")
	}
	cmdStr := msg.Command[0]
	args := msg.Command[1:]

	cmd := exec.Command(cmdStr, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}

	m := sync.Mutex{}
	var wg sync.WaitGroup
	wg.Add(2)
	const readSize = 8196
	var readErr []error
	go func() {
		defer wg.Done()
		for buf := make([]byte, readSize); ; {
			n, err := stdout.Read(buf)
			if n > 0 {
				m.Lock()
				resp := localhost.OutputMessage{
					Stdout: buf[:n],
				}
				fmt.Printf("sending message %v\n", resp)
				err := stream.SendMsg(&resp)
				if err != nil {
					readErr = append(readErr, err)
					m.Unlock()
					return
				}
				m.Unlock()
			}
			if err != nil {
				m.Lock()
				if err != io.EOF {
					readErr = append(readErr, err)
				}
				m.Unlock()
				return
			}
		}
	}()
	go func() {
		defer wg.Done()
		for buf := make([]byte, readSize); ; {
			n, err := stderr.Read(buf)
			if n > 0 {
				m.Lock()
				resp := localhost.OutputMessage{
					Stderr: buf[:n],
				}
				fmt.Printf("sending message %v\n", resp)
				err := stream.SendMsg(&resp)
				if err != nil {
					readErr = append(readErr, err)
					m.Unlock()
					return
				}
				m.Unlock()
			}
			if err != nil {
				m.Lock()
				if err != io.EOF {
					readErr = append(readErr, err)
				}
				m.Unlock()
				return
			}
		}
	}()

	wg.Wait()
	if len(readErr) != 0 {
		for _, err := range readErr {
			fmt.Fprintf(os.Stderr, "got error while reading from locally-run process: %v\n", err)
		}
		return readErr[0]
	}

	var exitCode int
	status := localhost.DONE
	err = cmd.Wait()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if waitStatus, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitCode = waitStatus.ExitStatus()
			} else {
				status = localhost.KILLED
			}
		} else {
			status = localhost.KILLED
		}
		return err
	}

	resp := localhost.OutputMessage{
		ExitCode: int32(exitCode),
		Status:   status,
	}
	if err := stream.SendMsg(&resp); err != nil {
		return err
	}

	return nil
}
