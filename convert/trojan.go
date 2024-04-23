package convert

import (
	"fmt"

	"github.com/ysicing/clash2singbox/model/clash"
	"github.com/ysicing/clash2singbox/model/singbox"
)

func trojan(p *clash.Proxies, s *singbox.SingBoxOut) error {
	p.Tls = true
	tls(p, s)
	if p.WsOpts.Path != "" || p.Network == "ws" {
		err := vmessWsOpts(p, s)
		if err != nil {
			return fmt.Errorf("trojan: %w", err)
		}
	}
	if p.GrpcOpts.GrpcServiceName != "" {
		err := vmessGrpcOpts(p, s)
		if err != nil {
			return fmt.Errorf("trojan: %w", err)
		}
	}
	return nil
}

func vmessGrpcOpts(p *clash.Proxies, s *singbox.SingBoxOut) error {
	if s.Transport == nil {
		s.Transport = &singbox.SingTransport{}
	}
	s.Transport.Type = "grpc"
	s.Transport.ServiceName = p.GrpcOpts.GrpcServiceName
	return nil
}

func vmessWsOpts(p *clash.Proxies, s *singbox.SingBoxOut) error {
	t := "ws"
	if p.WsOpts.V2rayHttpUpgrade {
		t = "httpupgrade"
	}
	if s.Transport == nil {
		s.Transport = &singbox.SingTransport{}
	}
	s.Transport.Type = t
	m := map[string][]string{}

	if len(p.WsHeaders) != 0 {
		for k, v := range p.WsHeaders {
			m[k] = []string{v}
		}
	} else {
		for k, v := range p.WsOpts.Headers {
			m[k] = []string{v}
		}
	}
	if p.WsOpts.V2rayHttpUpgrade {
		host := p.Servername
		if host == "" {
			host = p.WsOpts.Headers["Host"]
		}
		s.Transport.Host = host
	}
	s.Transport.Headers = m
	s.Transport.Path = p.WsOpts.Path
	s.Transport.EarlyDataHeaderName = p.WsOpts.EarlyDataHeaderName
	s.Transport.MaxEarlyData = int(p.WsOpts.MaxEarlyData)
	return nil
}
