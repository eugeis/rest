package core

import (
	"errors"
	"eye/digest"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
	"gee/as"
)

type Http struct {
	Name      string
	AccessKey string
	Url       string

	PingRequest *QueryRequest

	PingTimeoutMillis  int
	QueryTimeoutMillis int
}

type HttpService struct {
	http         *Http
	accessFinder as.AccessFinder

	client    *http.Client
	pingCheck *httpCheck

	pingTimeout  time.Duration
	queryTimeout time.Duration
}

func (o *HttpService) Name() string {
	return o.http.Name
}

func (o *HttpService) Init() (err error) {
	if o.client == nil {
		o.client = digest.NewClient(true, o.queryTimeout)
		o.pingCheck, err = o.newСheck(o.http.PingRequest)
		if err != nil {
			o.Close()
		}
	}
	return
}

func (o *HttpService) Close() {
	o.client = nil
	o.pingCheck = nil
}

func (o *HttpService) Ping() error {
	err := o.Init()
	if err == nil {
		err = o.pingCheck.Validate()
		if err != nil {
			Log.Debug("'%v' can't be reached because of %v", o.Name(), err)
		}
	}
	return err
}

func body(resp *http.Response) string {
	body, _ := ioutil.ReadAll(resp.Body)
	ret := fmt.Sprintf("%v", body)
	return ret
}

func (o *HttpService) NewСheck(req *QueryRequest) (ret Check, err error) {
	return o.newСheck(req)
}

func (o *HttpService) newСheck(req *QueryRequest) (ret *httpCheck, err error) {
	var access as.Access
	access, err = o.accessFinder.FindAccess(o.http.AccessKey)
	if err == nil {
		var pattern *regexp.Regexp
		if pattern, err = compilePattern(req.Expr); err == nil {
			dReq := digest.NewRequest(access.User, access.Password, "GET", o.http.Url+req.Query, "")
			ret = &httpCheck{
				info:    req.CheckKey(o.Name()),
				req:     &dReq,
				pattern: pattern, service: o, not: req.Not}
		}

	}
	return
}

func (o *HttpService) NewExporter(req *ExportRequest) (ret Exporter, err error) {
	return
}

//buildCheck

type httpCheck struct {
	info    string
	req     *digest.Request
	not     bool
	pattern *regexp.Regexp
	service *HttpService
}

func (o httpCheck) Info() string {
	return o.info
}

func (o httpCheck) Validate() error {
	return validate(o, o.pattern, o.not)
}

func (o httpCheck) Query() (data QueryResult, err error) {
	err = o.service.Init()
	if err != nil {
		return
	}
	resp, err := o.req.Execute(o.service.client)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	Log.Debug("http data: %s", data)

	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("Status %d", resp.StatusCode))
	}
	return
}
