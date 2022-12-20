package worker

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	zclient "github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/client"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/response"
	"github.com/gocraft/work"
	"go-micro.dev/v4/client"
	"go.opentelemetry.io/otel"
)

type workerContext struct{}

type callbackWorker struct {
	client        client.Client
	zoomFilestore zclient.ZoomFilestore
	logger        log.Logger
}

func NewWorkerContext() workerContext {
	return workerContext{}
}

func NewCallbackWorker(client client.Client, logger log.Logger) callbackWorker {
	return callbackWorker{
		client:        client,
		logger:        logger,
		zoomFilestore: zclient.NewZoomFilestoreClient(),
	}
}

func (c callbackWorker) UploadFile(job *work.Job) error {
	uid, filename, url := job.ArgString("uid"), job.ArgString("filename"), job.ArgString("url")

	c.logger.Debugf("got a new file %s upload job (%s)", filename, uid)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tracer := otel.GetTracerProvider().Tracer("zoom-onlyoffice/pool")
	tctx, span := tracer.Start(ctx, "upload")
	defer span.End()

	errorsChan := make(chan error, 2)
	fileChan := make(chan io.Reader)
	userChan := make(chan response.UserResponse)
	defer close(errorsChan)
	defer close(fileChan)
	defer close(userChan)

	var wg sync.WaitGroup

	go func() {
		wg.Add(1)
		defer wg.Done()
		file, err := c.zoomFilestore.GetFile(tctx, url)
		if err != nil {
			errorsChan <- err
			return
		}
		fileChan <- file
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		req := c.client.NewRequest("onlyoffice:auth", "UserHandler.GetUser", uid)

		var ures response.UserResponse
		if err := c.client.Call(tctx, req, &ures); err != nil {
			errorsChan <- err
			return
		}

		userChan <- ures
	}()

	wg.Wait()

	select {
	case err := <-errorsChan:
		c.logger.Debugf("could not execute a file %s upload operation (%s). Reason: %s", filename, uid, err.Error())
		return err
	default:
	}

	ures := <-userChan
	if err := c.zoomFilestore.UploadFile(tctx, ures.AccessToken, ures.ID, filename, <-fileChan); err != nil {
		c.logger.Debugf("could not upload an onlyoffice file to zoom: %s", err.Error())
		return err
	}

	return nil
}
