package convert

import (
	"github.com/ysicing/clash2singbox/model/clash"
	"github.com/ysicing/clash2singbox/model/singbox"
)

func httpOpts(p *clash.Proxies, s *singbox.SingBoxOut) error {
	tls(p, s)
	p.Username = s.Username
	return nil
}

func socks5(p *clash.Proxies, s *singbox.SingBoxOut) error {
	tls(p, s)
	p.Username = s.Username
	if !p.Udp {
		s.Network = "tcp"
	}
	return nil
}
