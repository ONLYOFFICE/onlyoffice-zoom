package handler

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	plog "github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/builder/web/server/core/domain"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/builder/web/server/core/port"
	zclient "github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/client"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/constants"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/crypto"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/request"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/response"
	"github.com/google/uuid"
	"github.com/mileusna/useragent"
	"go-micro.dev/v4/client"
	"golang.org/x/sync/singleflight"
)

type ConfigHandler struct {
	namespace   string
	logger      plog.Logger
	client      client.Client
	zoomAPI     zclient.ZoomAuth
	service     port.SessionService
	jwtManager  crypto.JwtManager
	callbackURL string
	group       singleflight.Group
}

func NewConfigHandler(
	namespace string,
	logger plog.Logger,
	client client.Client,
	zoomAPI zclient.ZoomAuth,
	service port.SessionService,
	jwtManager crypto.JwtManager,
	callbackURL string,
) ConfigHandler {
	return ConfigHandler{
		namespace:   namespace,
		logger:      logger,
		client:      client,
		zoomAPI:     zoomAPI,
		service:     service,
		jwtManager:  jwtManager,
		callbackURL: callbackURL,
	}
}

func (c ConfigHandler) processConfig(user response.UserResponse, request request.BuildConfigRequest, ctx context.Context) (response.BuildConfigResponse, error) {
	var config response.BuildConfigResponse

	u, err := c.zoomAPI.GetZoomUser(ctx, user.AccessToken)
	if err != nil {
		return config, err
	}

	t := "desktop"
	ua := useragent.Parse(request.UserAgent)

	if ua.Mobile || ua.Tablet {
		t = "mobile"
	}

	config = response.BuildConfigResponse{
		Document: response.Document{
			Key:   uuid.NewString(),
			Title: request.Filename,
			URL:   request.FileURL,
			Permissions: response.Permissions{
				Edit: false,
			},
		},
		EditorConfig: response.EditorConfig{
			User: response.User{
				ID:   u.ID,
				Name: strings.Join([]string{u.Firstname, u.Lastname}, " "),
			},
			CallbackURL: c.callbackURL,
			Customization: response.Customization{
				Goback: response.Goback{
					RequestClose: true,
				},
			},
			Lang: u.Language,
		},
		Type:  t,
		Owner: true,
	}

	if strings.TrimSpace(request.Filename) != "" {
		ext := strings.ReplaceAll(filepath.Ext(request.Filename), ".", "")
		fileType, err := constants.GetFileType(ext)
		if err != nil {
			return config, err
		}
		config.Document.FileType = ext
		config.Document.Permissions = response.Permissions{
			Edit: constants.IsExtensionEditable(ext),
		}
		config.DocumentType = fileType
	}

	return config, nil
}

func (c ConfigHandler) BuildConfig(ctx context.Context, payload request.BuildConfigRequest, res *response.BuildConfigResponse) error {
	c.logger.Debugf("processing a docs config: %s", payload.Filename)

	config, err, _ := c.group.Do(payload.Uid, func() (interface{}, error) {
		req := c.client.NewRequest(fmt.Sprintf("%s:auth", c.namespace), "UserSelectHandler.GetUser", payload.Uid)

		var ures response.UserResponse
		if err := c.client.Call(ctx, req, &ures); err != nil {
			return nil, err
		}

		config, err := c.processConfig(ures, payload, ctx)
		if err != nil {
			return nil, err
		}

		cbURL := config.EditorConfig.CallbackURL
		if payload.Mid == "" {
			c.logger.Debugf("request has no mid")
			config.IssuedAt, config.ExpiresAt = time.Now().UnixMilli(), time.Now().Add(3*time.Minute).UnixMilli()
			config.EditorConfig.CallbackURL = fmt.Sprintf("%s?filename=%s", cbURL, url.QueryEscape(payload.Filename))
			if config.Token, err = c.jwtManager.Sign(config); err != nil {
				c.logger.Errorf("could not sign a docs config. Error: %s", err.Error())
				return nil, err
			}
			return config, nil
		}

		md := md5.Sum([]byte(payload.Mid))
		payload.Mid = hex.EncodeToString(md[:])
		c.logger.Debugf("appending mid to callback url: %s", payload.Mid)
		config.EditorConfig.CallbackURL = fmt.Sprintf("%s?mid=%s&filename=%s", cbURL, payload.Mid, url.QueryEscape(payload.Filename))

		c.logger.Debugf("trying to find a docs session for mid: %s", payload.Mid)
		if session, err := c.service.GetSession(ctx, payload.Mid); err == nil {
			c.logger.Debugf("mid %s session has been found", payload.Mid)
			ext := strings.ReplaceAll(filepath.Ext(session.Filename), ".", "")
			fileType, err := constants.GetFileType(ext)
			if err != nil {
				c.logger.Errorf("could not get fileType for mid %s. Error: %s", payload.Mid, err.Error())
				return nil, err
			}

			config.Session = true
			config.Owner = session.Owner == payload.Uid
			config.Document = response.Document{
				FileType: ext,
				Key:      session.DocKey,
				Title:    session.Filename,
				URL:      session.FileURL,
				Permissions: response.Permissions{
					Edit: constants.IsExtensionEditable(ext),
				},
			}
			config.DocumentType = fileType
			config.EditorConfig.CallbackURL = fmt.Sprintf("%s?mid=%s&filename=%s", cbURL, payload.Mid, url.QueryEscape(session.Filename))
			config.IssuedAt, config.ExpiresAt = time.Now().UnixMilli(), time.Now().Add(3*time.Minute).UnixMilli()
			if config.Token, err = c.jwtManager.Sign(config); err != nil {
				c.logger.Errorf("could not sign a docs config for mid: %s. Error: %s", payload.Mid, err.Error())
				return nil, err
			}

			return config, nil
		} else {
			c.logger.Debugf("mid %s session hasn't been found. Generating a new one", payload.Mid)

			session, err := c.service.CreateSession(ctx, payload.Mid, domain.Session{
				Owner:    payload.Uid,
				Filename: payload.Filename,
				FileURL:  payload.FileURL,
				DocKey:   uuid.NewString(),
				Initial:  true,
			})

			if err != nil {
				c.logger.Error(err.Error())
				return nil, err
			}

			ext := strings.ReplaceAll(filepath.Ext(session.Filename), ".", "")
			fileType, err := constants.GetFileType(ext)
			if err != nil {
				c.logger.Errorf("could not get fileType for mid %s. Error: %s", payload.Mid, err.Error())
				return nil, err
			}

			config.Document = response.Document{
				FileType: ext,
				Key:      session.DocKey,
				Title:    session.Filename,
				URL:      session.FileURL,
				Permissions: response.Permissions{
					Edit: constants.IsExtensionEditable(ext),
				},
			}
			config.DocumentType = fileType
			config.IssuedAt, config.ExpiresAt = time.Now().UnixMilli(), time.Now().Add(3*time.Minute).UnixMilli()
			if config.Token, err = c.jwtManager.Sign(config); err != nil {
				c.logger.Errorf("could not sign a docs config for mid: %s. Error: %s", payload.Mid, err.Error())
				return nil, err
			}

			return config, nil
		}
	})

	if cfg, ok := config.(response.BuildConfigResponse); ok {
		*res = cfg
		return nil
	}

	return err
}
