// Package bazellogs provides functionality for reading bazel workspace logs.
//
// This Go package is based on the Jave implementation
// https://cs.opensource.google/bazel/bazel/+/master:src/tools/workspacelog/src/main/java/com/google/devtools/build/workspacelog/WorkspaceLogParser.java;l=35?q=parser%20workspacelog&ss=bazel%2Fbazel
// which is distributed under http://www.apache.org/licenses/LICENSE-2.0.
package bazellogs

import (
	"errors"
	"fmt"
	"io"

	"github.com/matttproud/golang_protobuf_extensions/pbutil"

	pb "github.com/gonzojive/bazelgopackagesdriver/proto/bazelworkspacelogpb"
)

// ReadWorkspaceEvents reads workspace events from a stream and calls callback with each event.
func ReadWorkspaceEvents(r io.Reader, callback func(ev *pb.WorkspaceEvent)) error {
	eventIndex := 0
	next := func() (*pb.WorkspaceEvent, error) {
		defer func() { eventIndex++ }()
		evt := &pb.WorkspaceEvent{}
		if _, err := pbutil.ReadDelimited(r, evt); err != nil {
			return nil, fmt.Errorf("failed to read event[%d] from workspace log: %w", eventIndex, err)
		}
		return evt, nil
	}
	for {
		ev, err := next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		callback(ev)
	}
}
