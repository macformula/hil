package results

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"go.uber.org/zap"

	"github.com/google/uuid"
	proto "github.com/macformula/hil/results/client/generated"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	_pythonUnix               = "python3"
	_pythonWin                = "python"
	_winOs                    = "windows"
	_loggerName               = "result_processor"
	_waitForFastFailErrorTime = 1 * time.Second
)

type ResultProcessor struct {
	l                   *zap.Logger
	addr                string
	conn                *grpc.ClientConn
	client              proto.TagTunnelClient
	pushReportsToGithub bool
	ra                  *ResultAccumulator
	serverAutoStart     bool
	configPath          string
	serverPath          string
	serverCmd           *exec.Cmd
}

type Option = func(*ResultProcessor)

// WithServerAutoStart will automatically start the result processor server. Server path should be the path to main.py.
func WithServerAutoStart(configPath, serverPath string) Option {
	return func(r *ResultProcessor) {
		r.serverAutoStart = true
		r.configPath = configPath
		r.serverPath = serverPath
	}
}

// WithPushReportsToGithub will push hil reports to the macfe-hil.github.io page.
func WithPushReportsToGithub() Option {
	return func(r *ResultProcessor) {
		r.pushReportsToGithub = true
	}
}

func NewResultProcessor(l *zap.Logger, address string, opts ...Option) *ResultProcessor {
	ret := &ResultProcessor{
		l:                   l.Named(_loggerName),
		addr:                address,
		pushReportsToGithub: false,
		serverAutoStart:     false,
	}
	for _, o := range opts {
		o(ret)
	}

	return ret
}

func (r *ResultProcessor) Open(ctx context.Context) error {
	var errCh = make(chan error)

	if r.serverAutoStart {
		go r.startServer(errCh)
	}

	select {
	case <-time.After(_waitForFastFailErrorTime):
	case err := <-errCh:
		return errors.Wrap(err, "start server")
	}

	conn, err := grpc.DialContext(ctx, r.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return errors.Wrap(err, "dial context")
	}

	r.conn = conn
	r.client = proto.NewTagTunnelClient(conn)

	return nil
}

func (r *ResultProcessor) SubmitTag(ctx context.Context, tag string, value any) (bool, error) {
	request, err := createRequest(tag, value) //checks if the tag is correct and error free
	r.ra.NewResultAccumulator()

	// fmt.Println("request ", request, " err ", err)
	// request  tag:"FW001" value_bool:true  err  <nil>

	if err != nil { //error check if create request worked? not sure how this works with an actual error because it is not showing
		return false, errors.Wrap(err, "create request")
	}

	reply, err := r.client.SubmitTag(ctx, request) //actually submitting the tag
	//fmt.Println("reply ", reply, " err ", err)
	// reply  success:true is_passing:true  err  <nil>
	if err != nil {
		return reply.IsPassing, errors.Wrap(err, "submit tag")
	}
	// if false means that the reply did not come through
	if !reply.Success {
		return false, errors.New(reply.Error)
	}
	return reply.IsPassing, nil //returns if the tag is passing or failing
}

func (r *ResultProcessor) CompleteTest(ctx context.Context, testId uuid.UUID, sequenceName string) (bool, error) {
	reply, err := r.client.CompleteTest(ctx, &proto.CompleteTestRequest{
		TestId:             testId.String(),
		SequenceName:       sequenceName,
		PushReportToGithub: r.pushReportsToGithub,
	})
	// runs ra.generate_and_run_tests and returns test passed
	// ra.generate_and_run_tests
	// returns a list of all of the tag ids that were run for a specific test (tag_ids = list(self.tag_submissions.keys())) but self.tag_submissions is a dict so no repeats
	// retuns the datetime for all of the tests that were performed (date_time = self.__generate_test_file(tag_ids)) all the tags have slightly different times as to when they were run
	// checks if the tag list has any errors (has_errors = len(self.error_submissions) > 0) as errors are caught they are counted in the error_submissions variable
	// (overall_pass_fail = self.all_tags_passing and (not has_errors)) redundancy check to see if all tags are true and error_submissions = 0
	// in the code, they add the report to the json here
	// then tag_submissions, error_submissions, and all_tags_passing are reset to defaul values
	if err != nil {
		return false, errors.Wrap(err, "complete test")
	}

	return reply.TestPassed, nil
}

func (r *ResultProcessor) SubmitError(ctx context.Context, err error) error {
	_, submitErr := r.client.SubmitError(ctx, &proto.SubmitErrorRequest{Error: err.Error()})
	if submitErr != nil {
		return errors.Wrap(err, "submit error")
	}
	// the err is appended to self.error_submissions variable then the length of the error list is returned (which should be one) then error checked
	return nil
}

func (r *ResultProcessor) Close() error {
	r.l.Info("closing result processor")

	if r.serverCmd != nil && r.serverCmd.Process != nil {
		r.l.Info("killing server process",
			zap.Int("pid", r.serverCmd.Process.Pid))

		err := r.serverCmd.Process.Kill()
		if err != nil {
			return errors.Wrap(err, "kill server process")
		}
	}
	return nil
}

func (r *ResultProcessor) startServer(errCh chan error) {
	configFlag := fmt.Sprintf("--config=%s", r.configPath)

	if runtime.GOOS == _winOs {
		r.serverCmd = exec.Command(_pythonWin, r.serverPath, configFlag)
	} else {
		r.serverCmd = exec.Command(_pythonUnix, r.serverPath, configFlag)
	}

	r.l.Info("starting results server", zap.String("command", r.serverCmd.String()))

	out, err := r.serverCmd.Output()
	if err != nil {
		errCh <- errors.Wrap(err, "run output: "+string(out))
	}
}

func createRequest(tag string, data any) (*proto.SubmitTagRequest, error) {
	request := &proto.SubmitTagRequest{Tag: tag}
	// sends in a tag and data, however it seems that data can either be true false, or it can be an int/float, or a string (which I assume in error)
	// this goes into the submit_tag function
	// before this code runs, when the code is initialized, it parses through the tags.yaml or tags_scheme.json and create a Tag dictionary.
	// What it does is that it goes through everything in tags.yaml and checks it again tags_schema.json to see if it is in the right format
	// then it creates a list of tags based off tags.yaml
	// uses the logic in tag.py (ispassing) to determine if that specific test is a pass or fail case
	// if it is a fail, all_tags_passing is set to false it is saved to another dict as a false then return false and no error
	// else it is saved to another dict as a true then return true and no error
	// else if you get a keyvalue error where value is not a bool, that means that it was an error case that is probably fatal meaning it returns false with an error
	// then we go up a level into grpc, where if there was an error success=False, error=f"unknown tag id ({str(e)})", is_passing=is_passing
	// else if there was not an error success=True, error="", is_passing=is_passing
	// remember ispassing = false and error are not the same thing

	switch data.(type) {
	case int32:
		request.Data = &proto.SubmitTagRequest_ValueInt{ValueInt: data.(int32)}
	case float32:
		request.Data = &proto.SubmitTagRequest_ValueFloat{ValueFloat: data.(float32)}
	case string:
		request.Data = &proto.SubmitTagRequest_ValueStr{ValueStr: data.(string)}
	case bool:
		request.Data = &proto.SubmitTagRequest_ValueBool{ValueBool: data.(bool)}
	default:
		return nil, errors.Errorf("unsupported data type for tag submission (%T)", data)
	}
	return request, nil
}
