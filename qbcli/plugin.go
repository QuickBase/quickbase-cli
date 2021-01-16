package qbcli

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
)

// LoggerPlugin implements qbclient.Plugin and logs requests.
type LoggerPlugin struct {
	ctx    context.Context
	logger *cliutil.LeveledLogger
}

// NewLoggerPlugin returns a LoggerPlugin, which implements qbclient.Plugin.
func NewLoggerPlugin(ctx context.Context, logger *cliutil.LeveledLogger) qbclient.Plugin {
	return LoggerPlugin{ctx: ctx, logger: logger}
}

// PreRequest implements qbclient.Plugin.PreRequest.
func (p LoggerPlugin) PreRequest(req *http.Request) {
	ctx := p.ctx
	ctx = cliutil.ContextWithLogTag(ctx, "method", req.Method)
	ctx = cliutil.ContextWithLogTag(ctx, "url", req.URL.String())
	p.logger.Debug(ctx, "api request constructed")
}

// PostResponse implements qbclient.Plugin.PostResponse.
func (p LoggerPlugin) PostResponse(resp *http.Response) {
	if resp != nil {
		ctx := p.ctx
		ctx = cliutil.ContextWithLogTag(ctx, "method", resp.Request.Method)
		ctx = cliutil.ContextWithLogTag(ctx, "url", resp.Request.URL.String())
		ctx = cliutil.ContextWithLogTag(ctx, "status", resp.Status)
		p.logger.Info(ctx, "api response returned")
	}
}

// DumpPlugin implements qbclient.Plugin and dumps requests and responses to
// files in a directory.
type DumpPlugin struct {
	ctx       context.Context
	directory string
	logger    *cliutil.LeveledLogger
	transid   string
}

// NewDumpPlugin returns a DumpPlugin, which implements qbclient.Plugin.
func NewDumpPlugin(ctx context.Context, logger *cliutil.LeveledLogger, transid string, directory string) qbclient.Plugin {
	dir := strings.TrimRight(directory, string(os.PathSeparator))
	return DumpPlugin{ctx: ctx, logger: logger, transid: transid, directory: dir}
}

// PreRequest implements qbclient.Plugin.PreRequest.
func (p DumpPlugin) PreRequest(req *http.Request) {
	ctx, file, err := p.openDumpFile("request")
	if err != nil {
		return
	}
	defer file.Close()

	// Dump the request headers.
	headers, err := httputil.DumpRequestOut(req, false)
	if err != nil {
		p.logger.Error(ctx, "error dumping request headers", err)
		return
	}

	// Read the request body.
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		p.logger.Error(ctx, "error reading request body", err)
		return
	}

	// Build the request sent over the wire.
	buf := bytes.NewBuffer([]byte(``))
	buf.Write(headers)
	buf.Write(body)

	// Mask any user tokens.
	dump := qbclient.MaskUserToken(buf.Bytes())

	// Write the request to the dump file, and log the result.
	n, err := file.Write(dump)
	ctx = cliutil.ContextWithLogTag(ctx, "bytes", strconv.Itoa(n))
	if err == nil {
		p.logger.Debug(ctx, "wrote request to dump file")
	} else {
		p.logger.Error(ctx, "error writing request to dump file", err)
	}

	// Put the request body back so we can read it again.
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
}

// PostResponse implements qbclient.Plugin.PostResponse.
func (p DumpPlugin) PostResponse(resp *http.Response) {
	ctx, file, err := p.openDumpFile("response")
	if err != nil {
		return
	}
	defer file.Close()

	// Dump the response headers.
	headers, err := httputil.DumpResponse(resp, false)
	if err != nil {
		p.logger.Error(ctx, "error dumping response headers", err)
		return
	}

	// Read the response body.
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		p.logger.Error(ctx, "error reading response body", err)
		return
	}

	// Build the response returned over the wire.
	buf := bytes.NewBuffer([]byte(``))
	buf.Write(headers)
	buf.Write(body)

	// Mask any user tokens.
	dump := qbclient.MaskUserToken(buf.Bytes())

	// Write the response to the dump file, and log the result.
	n, err := file.Write(dump)
	ctx = cliutil.ContextWithLogTag(ctx, "bytes", strconv.Itoa(n))
	if err == nil {
		p.logger.Debug(ctx, "wrote response to dump file")
	} else {
		p.logger.Error(ctx, "error writing response to dump file", err)
	}

	// Put the response body back so we can read it again.
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
}

func (p DumpPlugin) openDumpFile(op string) (ctx context.Context, file *os.File, err error) {
	filename := fmt.Sprintf("%v-%s-%s.txt", time.Now().Unix(), p.transid, op)
	filepath := qbclient.Filepath(p.directory, filename)

	ctx = cliutil.ContextWithLogTag(p.ctx, "file", filepath)
	ctx = cliutil.ContextWithLogTag(p.ctx, "operation", op)

	file, err = os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err == nil {
		p.logger.Debug(ctx, "created dump file")
	} else {
		p.logger.Error(ctx, "error creating dump file", err)
	}

	return
}
