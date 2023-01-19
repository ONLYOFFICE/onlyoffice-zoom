package worker

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	zclient "github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/client"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/response"
	"github.com/gocraft/work"
	"go-micro.dev/v4/client"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

type workerContext struct{}

type callbackWorker struct {
	client        client.Client
	zoomFilestore zclient.ZoomFilestore
	uploadTimeout int
	logger        log.Logger
}

func NewWorkerContext() workerContext {
	return workerContext{}
}

func NewCallbackWorker(client client.Client, uploadTimeout int, logger log.Logger) callbackWorker {
	return callbackWorker{
		client:        client,
		zoomFilestore: zclient.NewZoomFilestoreClient(),
		uploadTimeout: uploadTimeout,
		logger:        logger,
	}
}

func (c callbackWorker) UploadFile(job *work.Job) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.uploadTimeout)*time.Second)
	defer cancel()

	tracer := otel.GetTracerProvider().Tracer("zoom-onlyoffice/pool")
	tctx, span := tracer.Start(ctx, "upload")
	defer span.End()

	uid, filename, url := job.ArgString("uid"), job.ArgString("filename"), job.ArgString("url")
	c.logger.Debugf("got a new file %s upload job (%s)", filename, uid)

	var wg sync.WaitGroup
	userChan := make(chan response.UserResponse)
	sizeChan := make(chan int64)
	errChan := make(chan error, 2)

	go func() {
		wg.Add(1)
		defer wg.Done()

		req := c.client.NewRequest("onlyoffice:auth", "UserSelectHandler.GetUser", uid)
		var ures response.UserResponse
		if err := c.client.Call(tctx, req, &ures); err != nil {
			errChan <- err
			return
		}

		userChan <- ures
	}()

	go func() {
		wg.Add(1)
		defer wg.Wait()

		headResp, _ := otelhttp.Head(tctx, url)
		size, err := strconv.ParseInt(headResp.Header.Get("Content-Length"), 10, 64)
		if err != nil {
			errChan <- err
			return
		}

		sizeChan <- size
	}()

	wg.Wait()

	select {
	case err := <-errChan:
		return err
	default:
	}

	ures := <-userChan
	if err := c.zoomFilestore.UploadFile(tctx, url, ures.AccessToken, ures.ID, filename, <-sizeChan); err != nil {
		c.logger.Debugf("could not upload an onlyoffice file to zoom: %s", err.Error())
		return err
	}

	return nil
}
