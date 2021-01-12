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
	// stdout and stderr get closed in execOp.Exec()

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

	var exitCodeSet bool
	var exitCode int
	for {
		var msg OutputMessage
		err := stream.RecvMsg(&msg)
		if err != nil {
			if err == io.EOF {
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
