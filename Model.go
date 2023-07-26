package system_settings

import (
	"encoding/json"

	"github.com/go-bolo/bolo"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Settings - Settings Model
type Settings struct {
	Key   string `gorm:"primary_key;type:varchar(255);not null" json:"key"`
	Value string `gorm:"type:varchar(255);not null" json:"value"`
}

type SettingsQueryOpts struct {
	Records *[]Settings
	Count   *int64
	Limit   int
	Offset  int
	C       echo.Context
	IsHTML  bool
}

func (r *Settings) Delete() error {
	db := bolo.GetDefaultDatabaseConnection()
	return db.Unscoped().Delete(&r).Error
}

func (m *Settings) Save() error {
	db := bolo.GetDefaultDatabaseConnection()

	saved := Settings{}
	err := SettingsFindOne(m.Key, &saved)
	if err != nil {
		return err
	}

	if saved.Key == "" {
		// create ....
		err = db.Create(&m).Error
		if err != nil {
			return err
		}
	} else {
		// update ...
		err := db.Save(&m).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// TableName - Method for set settings tablename as settings
func (r *Settings) TableName() string {
	return "system_settings"
}

func (r *Settings) ToJSON() string {
	jsonString, _ := json.MarshalIndent(r, "", "  ")
	return string(jsonString)
}

func (r *Settings) LoadData() error {
	return nil
}

func (r *Settings) GetKey() string {
	return r.Key
}

func SettingsFindOne(key string, target *Settings) error {
	logrus.WithFields(logrus.Fields{
		"key": key,
	}).Debug("FindSettingsByKey Will find by")

	db := bolo.GetDefaultDatabaseConnection()

	if err := db.
		Where("`key` = ? ", key).
		First(target).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	}

	return nil
}

func SettingsQueryAndCountReq(opts *SettingsQueryOpts) error {
	db := bolo.GetDefaultDatabaseConnection()

	c := opts.C

	q := c.QueryParam("q")

	query := db

	if q != "" {
		query = db.Where("key LIKE ?", "%"+q+"%")
	}

	err := query.Limit(opts.Limit).
		Offset(opts.Offset).
		Find(opts.Records).Error
	if err != nil {
		return err
	}

	return SettingsCountReq(opts)

}

func SettingsCountReq(opts *SettingsQueryOpts) error {
	db := bolo.GetDefaultDatabaseConnection()

	c := opts.C

	q := c.QueryParam("q")

	// Count ...
	queryCount := db

	if q != "" {
		queryCount = queryCount.Or(
			db.Where("key LIKE ?", "%"+q+"%"),
		)
	}

	return queryCount.
		Table("system_settings").
		Count(opts.Count).Error

}

func SettingsFindAll(records *[]*Settings, limit int) error {
	db := bolo.GetDefaultDatabaseConnection()

	if err := db.
		Limit(limit).
		Find(records).Error; err != nil {
		return err
	}

	return nil
}

func FindAllAsMap() (map[string]string, error) {
	cachedData := GetAllFromCache()
	if cachedData != nil {
		return cachedData, nil
	}

	s := map[string]string{}

	records := []*Settings{}
	err := SettingsFindAll(&records, 10000)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"records": records,
			"error":   err,
		}).Debug("system_settings form cache")
		return s, err
	}

	for _, v := range records {
		s[v.Key] = v.Value
	}

	SetAllInCache(s)

	return s, nil
}

// Get system setting by key, returns nil if not set
func Get(key string) string {
	d, _ := FindAllAsMap()

	if v, ok := d[key]; ok {
		return v
	}

	return ""
}

// Get with default
func GetD(key string, d string) string {
	data, _ := FindAllAsMap()

	if v, ok := data[key]; ok {
		return v
	}

	return d
}

func Set(key string, value string) {

}
