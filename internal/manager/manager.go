package manager

import (
	"go-broker/internal/storage"
	"go-broker/internal/tcp"
)

type Manager struct {
	socketServer *tcp.Server
	storage      storage.Storage
	router       *Router
}

func InitManager(basePath string) (*Manager, error) {

	// create socket server config
	socketServerConf := tcp.ServerConfig{
		ConnectionPort: 8085,
	}

	// create socket server
	socketServer := tcp.New(socketServerConf)

	// start socket server
	socketServer.Start()

	// init router
	router := NewRouter()

	storageConfig := storage.StorageConfig{
		Path:           "C:\\Users\\m.shakiba.PSZ021-PC\\Desktop\\data",
		FileMaxSize:    1024,
		FileNamePrefix: "go",
	}

	// init storage
	s := storage.New(storageConfig)

	err := s.Init()

	if err != nil {
		return nil, err
	}

	mgr := &Manager{
		socketServer: socketServer,
		storage:      s,
		router:       router,
	}

	// register handlers
	socketServer.RegisterHandler("SUB", mgr.handleSubscribeMessage)
	socketServer.RegisterHandler("PUB", mgr.handlePublishMessage)
	socketServer.RegisterHandler("ACK", mgr.handlePublishMessage)

	return mgr, nil
}

func (m *Manager) processMessage(p *PayloadMessage) {
	subscribers := m.router.Match(p.Routes)

	for _, s := range subscribers {
		s.OnMessage(p)
	}
}
