package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	zclient "github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/client"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/message"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/response"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/util/backoff"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

type callbackWorker struct {
	namespace     string
	client        client.Client
	zoomFilestore zclient.ZoomFilestore
	uploadTimeout int
	logger        log.Logger
}

func NewCallbackWorker(namespace string, client client.Client, uploadTimeout int, logger log.Logger) callbackWorker {
	return callbackWorker{
		namespace:     namespace,
		client:        client,
		zoomFilestore: zclient.NewZoomFilestoreClient(logger),
		uploadTimeout: uploadTimeout,
		logger:        logger,
	}
}

func (c callbackWorker) UploadFile(ctx context.Context, payload []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.uploadTimeout)*time.Second)
	defer cancel()

	tracer := otel.GetTracerProvider().Tracer("zoom-onlyoffice/pool")
	tctx, span := tracer.Start(ctx, "upload")
	defer span.End()

	var msg message.JobMessage
	if err := json.Unmarshal(payload, &msg); err != nil {
		logger.Errorf("could not notify ws clients. Reason: %s", err.Error())
		return err
	}

	c.logger.Debugf("got a new file %s upload job (%s)", msg.Filename, msg.UID)

	var wg sync.WaitGroup
	userChan := make(chan response.UserResponse)
	sizeChan := make(chan int64)
	ferrChan := make(chan error)
	serrChan := make(chan error)

	req := c.client.NewRequest(fmt.Sprintf("%s:auth", c.namespace), "UserSelectHandler.GetUser", msg.UID)

	go func() {
		wg.Add(1)
		defer wg.Done()

		c.logger.Debugf("trying to get an access token")
		var ures response.UserResponse
		if res, ok := c.client.Options().Cache.Get(ctx, &req); ok {
			ures = res.(response.UserResponse)
		} else {
			if err := c.client.Call(tctx, req, &ures, client.WithRetries(3), client.WithBackoff(func(ctx context.Context, req client.Request, attempts int) (time.Duration, error) {
				return backoff.Do(attempts), nil
			})); err != nil {
				ferrChan <- err
				return
			}
			c.client.Options().Cache.Set(ctx, &req, ures, time.Duration((ures.ExpiresAt-time.Now().UnixMilli())*1e6/6))
		}

		c.logger.Debugf("populating user channel")
		userChan <- ures
		c.logger.Debugf("successfully populated user channel")
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()

		headResp, err := otelhttp.Head(tctx, msg.Url)
		if err != nil {
			serrChan <- err
			return
		}

		size, err := strconv.ParseInt(headResp.Header.Get("Content-Length"), 10, 64)
		if err != nil {
			serrChan <- err
			return
		}

		c.logger.Debugf("populating file size channel")
		sizeChan <- size
		c.logger.Debugf("successfully populated file size channel")
	}()

	select {
	case err := <-ferrChan:
		return err
	case err := <-serrChan:
		return err
	case ures := <-userChan:
		if err := c.zoomFilestore.UploadFile(tctx, msg.Url, ures.AccessToken, ures.ID, msg.Filename, <-sizeChan); err != nil {
			c.logger.Errorf("could not upload an onlyoffice file to zoom: %s", err.Error())
			c.client.Options().Cache.Set(ctx, &req, nil, time.Duration(time.Now().Add(1*time.Second).Nanosecond()))
			return err
		}
		return nil
	}
}
