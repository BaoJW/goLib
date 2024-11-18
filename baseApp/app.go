package baseApp

type Config struct {
}

type BaseApp struct {
	config Config
	http   HttpConfig
}

func New(cfg Config) *BaseApp {
	// TODO 关于一些默认初始化的配置可以在这里增加 比如权限的初始化校验等
	return &BaseApp{config: cfg}
}

func (ba *BaseApp) Start() {
	if ba.http.ready {
		go ba.startHttp()
	}
}
