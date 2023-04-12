package discovery

import (
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/qrinef/arbitrage-bot/services/storage"
	"go.uber.org/zap"
	"net"
	"sync"
)

type Service struct {
	storageService *storage.Service
	logger         *zap.SugaredLogger
	cfg            discover.Config
	listen         *discover.UDPv5
}

func NewService(storageService *storage.Service, logger *zap.SugaredLogger) *Service {
	bootNodes := []*enode.Node{
		enode.MustParse("enode://d860a01f9722d78051619d1e2351aba3f43f943f6f00718d1b9baa4101932a1f5011f16bb2b1bb35db20d6fe28fa0bf09636d26a87d31de9ec6203eeedb1f666@18.138.108.67:30303"),
		enode.MustParse("enode://22a8232c3abc76a16ae9d6c3b164f98775fe226f0917b0ca871128a74a8e9630b458460865bab457221f1d448dd9791d24c4e5d88786180ac185df813a68d4de@3.209.45.79:30303"),
		enode.MustParse("enode://8499da03c47d637b20eee24eec3c356c9a2e6148d6fe25ca195c7949ab8ec2c03e3556126b0d7ed644675e78c4318b08691b7b57de10e5f0d40d05b09238fa0a@52.187.207.27:30303"),
		enode.MustParse("enode://103858bdb88756c71f15e9b5e09b56dc1be52f0a5021d46301dbbfb7e130029cc9d0d6f73f693bc29b665770fff7da4d34f3c6379fe12721b5d7a0bcb5ca1fc1@191.234.162.198:30303"),
		enode.MustParse("enode://715171f50508aba88aecd1250af392a45a330af91d7b90701c436b618c86aaa1589c9184561907bebbb56439b8f8787bc01f49a7c77276c58c1b09822d75e8e8@52.231.165.108:30303"),
		enode.MustParse("enode://5d6d7cd20d6da4bb83a1d28cadb5d409b64edf314c0335df658c1a54e32c7c4a7ab7823d57c39b6a757556e68ff1df17c748b698544a55cb488b52479a92b60f@104.42.217.25:30303"),
		enode.MustParse("enode://2b252ab6a1d0f971d9722cb839a42cb81db019ba44c08754628ab4a823487071b5695317c8ccd085219c3a03af063495b2f1da8d18218da2d6a82981b45e6ffc@65.108.70.101:30303"),
		enode.MustParse("enode://4aeb4ab6c14b23e2c4cfdce879c04b0748a20d8e9b59e25ded2a08143e265c6c25936e74cbc8e641e3312ca288673d91f2f93f8e277de3cfa444ecdaaf982052@157.90.35.166:30303"),
	}

	privateKey, _ := crypto.GenerateKey()

	return &Service{
		storageService: storageService,
		logger:         logger,
		cfg: discover.Config{
			PrivateKey: privateKey,
			Bootnodes:  bootNodes,
		},
	}
}

func (s *Service) Start() {
	listen, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IP{0, 0, 0, 0}, Port: 3351})
	if err != nil {
		s.logger.Fatal(err)
	}

	db, err := enode.OpenDB("")
	if err != nil {
		s.logger.Fatal(err)
	}

	s.listen, err = discover.ListenV5(listen, enode.NewLocalNode(db, s.cfg.PrivateKey), s.cfg)
	if err != nil {
		s.logger.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		it := s.listen.RandomNodes()

		for it.Next() {
			s.logger.Info(it.Node().URLv4())
		}
	}()

	wg.Wait()
}
