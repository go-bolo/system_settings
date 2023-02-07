package system_settings

import (
	"time"

	"github.com/jellydator/ttlcache/v3"
)

const systemSettingsCacheKey = "system_settings"

var Cache = ttlcache.Cache[string, map[string]string]{}

func GetAllFromCache() map[string]string {
	item := Cache.Get(systemSettingsCacheKey)
	if item == nil || item.IsExpired() {
		return nil
	}

	v := item.Value()
	return v
}

func SetAllInCache(value map[string]string) {
	Cache.Set(systemSettingsCacheKey, value, time.Minute*2)
}

func RefreshSomeItems(data map[string]string) {
	item := Cache.Get(systemSettingsCacheKey)
	if item == nil || item.IsExpired() {
		SetAllInCache(data)
		return
	}

	cachedData := item.Value()

	for k, v := range data {
		cachedData[k] = v
	}
	SetAllInCache(cachedData)
}
