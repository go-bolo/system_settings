package system_settings

import (
	"net/http"
	"strconv"

	"github.com/go-catupiry/catu"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type SettingsJSONResponse struct {
	catu.BaseListReponse
	Settings *[]Settings `json:"system-settings"`
}

type settingsCountJSONResponse struct {
	catu.BaseMetaResponse
}

type SettingsFindOneJSONResponse struct {
	Settings *Settings `json:"system-settings"`
}

type SettingsCreateJSONResponse struct {
	Settings map[string]string `json:"system-settings"`
}

type SettingsBodyRequest struct {
	Settings map[string]interface{} `json:"system-settings"` //
}

func (b *SettingsBodyRequest) GetValue(key string) string {
	d := b.Settings[key]

	switch v := d.(type) {
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.Itoa(int(v))
	case string:
		return d.(string)
	default:
		return ""
	}
}

func (b *SettingsBodyRequest) GetAllValues() map[string]string {
	d := map[string]string{}

	for k, _ := range b.Settings {
		d[k] = b.GetValue(k)
	}

	return d
}

// Http symbol controller | struct with http handlers
type SettingsController struct {
	App catu.App
}

func (ctl *SettingsController) Query(c echo.Context) error {
	var err error
	RequestContext := c.(*catu.RequestContext)

	var count int64
	var records []Settings
	err = SettingsQueryAndCountReq(&SettingsQueryOpts{
		Records: &records,
		Count:   &count,
		Limit:   RequestContext.GetLimit(),
		Offset:  RequestContext.GetOffset(),
		C:       c,
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Debug("SettingsFindAll Error on find settings")
	}

	RequestContext.Pager.Count = count

	logrus.WithFields(logrus.Fields{
		"count":             count,
		"len_records_found": len(records),
	}).Debug("SettingsFindAll count result")

	for i := range records {
		records[i].LoadData()
	}

	resp := SettingsJSONResponse{
		Settings: &records,
	}

	resp.Meta.Count = count

	return c.JSON(200, &resp)
}

func (ctl *SettingsController) Create(c echo.Context) error {
	logrus.Debug("SettingsController.Create running")
	var err error
	ctx := c.(*catu.RequestContext)

	can := ctx.Can("create_settings")
	if !can {
		return echo.NewHTTPError(http.StatusForbidden, "Forbidden")
	}

	var body SettingsBodyRequest

	if err := c.Bind(&body); err != nil {
		if _, ok := err.(*echo.HTTPError); ok {
			return err
		}
		return c.NoContent(http.StatusNotFound)
	}

	bodyData := body.GetAllValues()

	ms := map[string]string{}

	for k, v := range bodyData {
		s := &Settings{Key: k, Value: v}

		err = s.Save()
		if err != nil {
			return err
		}

		ms[k] = v
	}

	RefreshSomeItems(ms)

	logrus.WithFields(logrus.Fields{
		"body": body,
	}).Info("SettingsController.Create params")

	return c.JSON(http.StatusCreated, SettingsCreateJSONResponse{Settings: GetAllFromCache()})
}

func (ctl *SettingsController) Count(c echo.Context) error {
	var err error
	RequestContext := c.Get("ctx").(*catu.RequestContext)

	var count int64
	err = SettingsCountReq(&SettingsQueryOpts{
		Count:  &count,
		Limit:  RequestContext.GetLimit(),
		Offset: RequestContext.GetOffset(),
		C:      c,
	})

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Debug("SettingsFindAll Error on find symbols")
	}

	RequestContext.Pager.Count = count

	resp := settingsCountJSONResponse{}
	resp.Count = count

	return c.JSON(200, &resp)
}

func (ctl *SettingsController) FindOne(c echo.Context) error {
	key := c.Param("key")

	logrus.WithFields(logrus.Fields{
		"key": key,
	}).Debug("FindOne key from params")

	record := Settings{}
	SettingsFindOne(key, &record)

	if record.Key == "" {
		logrus.WithFields(logrus.Fields{
			"key": key,
		}).Debug("FindOneHandler key record not found")

		return echo.NotFoundHandler(c)
	}

	record.LoadData()

	resp := SettingsFindOneJSONResponse{
		Settings: &record,
	}

	return c.JSON(200, &resp)
}

func (ctl *SettingsController) Update(c echo.Context) error {
	return ctl.Create(c) // upsert
}

func (ctl *SettingsController) Delete(c echo.Context) error {
	var err error

	key := c.Param("key")

	logrus.WithFields(logrus.Fields{
		"key": key,
	}).Debug("settings.DeleteOneHandler key from params")

	RequestContext := c.Get("ctx").(*catu.RequestContext)

	can := RequestContext.Can("delete_widget")
	if !can {
		return echo.NewHTTPError(http.StatusForbidden, "Forbidden")
	}

	record := Settings{}
	err = SettingsFindOne(key, &record)
	if err != nil {
		return err
	}

	err = record.Delete()
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

type SettingsControllerConfiguration struct {
}

func NewSymbolController(cfg *SettingsControllerConfiguration) *SettingsController {
	ctx := SettingsController{}

	return &ctx
}
