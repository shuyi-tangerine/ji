package ji

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
)

type ThriftServer struct {
	isServer   bool
	protocol   string
	isBuffered bool
	isFramed   bool
	addr       string
	useSecure  bool
	c          chan error
	processor  thrift.TProcessor
}

func NewThriftServer(protocol string, isBuffered bool, isFramed bool, addr string, useSecure bool, processor thrift.TProcessor) *ThriftServer {
	server := &ThriftServer{
		protocol:   protocol,
		isBuffered: isBuffered,
		isFramed:   isFramed,
		addr:       addr,
		useSecure:  useSecure,
		c:          make(chan error),
		processor:  processor,
	}
	var _ Server = server
	return server
}

func (m *ThriftServer) Start(ctx context.Context) (err error) {
	var protocolFactory thrift.TProtocolFactory
	switch m.protocol {
	case "compact":
		protocolFactory = thrift.NewTCompactProtocolFactoryConf(nil)
	case "simplejson":
		protocolFactory = thrift.NewTSimpleJSONProtocolFactoryConf(nil)
	case "json":
		protocolFactory = thrift.NewTJSONProtocolFactory()
	case "binary", "":
		protocolFactory = thrift.NewTBinaryProtocolFactoryConf(nil)
	default:
		return fmt.Errorf("invalid protocol specified %s", m.protocol)
	}

	var transportFactory thrift.TTransportFactory
	cfg := &thrift.TConfiguration{
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	if m.isBuffered {
		transportFactory = thrift.NewTBufferedTransportFactory(8192)
	} else {
		transportFactory = thrift.NewTTransportFactory()
	}

	if m.isFramed {
		transportFactory = thrift.NewTFramedTransportFactoryConf(transportFactory, cfg)
	}

	var transport thrift.TServerTransport
	if m.useSecure {
		cfg := new(tls.Config)
		if cert, err := tls.LoadX509KeyPair("server.crt", "server.key"); err == nil {
			cfg.Certificates = append(cfg.Certificates, cert)
		} else {
			return err
		}
		transport, err = thrift.NewTSSLServerSocket(m.addr, cfg)
		if err != nil {
			return err
		}
	} else {
		transport, err = thrift.NewTServerSocket(m.addr)
		if err != nil {
			return err
		}
	}
	server := thrift.NewTSimpleServer4(m.processor, transport, transportFactory, protocolFactory)
	fmt.Println("Starting the rpc server... on ", m.addr)
	return server.Serve()
}

func (m *ThriftServer) AsyncStart(ctx context.Context) {
	go func() {
		err := m.Start(ctx)
		if err != nil {
			fmt.Println("[AsyncStart] Start panic", err)
			m.c <- err
		}
	}()
}

func (m *ThriftServer) ErrorC() (c chan error) {
	return m.c
}
