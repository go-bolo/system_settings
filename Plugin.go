package system_settings

import (
	"time"

	"github.com/go-bolo/bolo"
	"github.com/go-bolo/system_settings/migrations"
	"github.com/gookit/event"
	"github.com/jellydator/ttlcache/v3"
	"github.com/sirupsen/logrus"
)

type Plugin struct {
	bolo.Pluginer
	Name       string
	Controller *SettingsController
}

func (p *Plugin) GetName() string {
	return p.Name
}

func (p *Plugin) Init(app bolo.App) error {
	p.Controller = &SettingsController{
		App: app,
	}

	Cache = *ttlcache.New[string, map[string]string](
		ttlcache.WithTTL[string, map[string]string](2 * time.Minute),
	)

	app.GetEvents().On("bindRoutes", event.ListenerFunc(func(e event.Event) error {
		return p.BindRoutes(app)
	}), event.Normal)

	return nil
}

func (r *Plugin) BindRoutes(app bolo.App) error {
	logrus.Debug(r.GetName() + " BindRoutes")

	ctl := r.Controller

	routerApi := app.SetRouterGroup("system-settings-api", "/api/system-settings")
	app.SetResource("system-settings", ctl, routerApi)

	routerApiOld := app.SetRouterGroup("system-settings-api-old", "/system-settings")
	app.SetResource("system-settings-old", ctl, routerApiOld)

	return nil
}

func (r *Plugin) GetMigrations() []*bolo.Migration {
	return []*bolo.Migration{
		migrations.Get00001Migration(),
	}
}

type PluginCfgs struct {
	PublicSystemSettings map[string]bool
}

func NewPlugin(cfg *PluginCfgs) *Plugin {
	p := Plugin{Name: "settings"}
	if cfg.PublicSystemSettings == nil {
		cfg.PublicSystemSettings = map[string]bool{}
	}
	return &p
}
