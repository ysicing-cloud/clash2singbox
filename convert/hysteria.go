package convert

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/ysicing/clash2singbox/model/clash"
	"github.com/ysicing/clash2singbox/model/singbox"
)

func hysteia2(p *clash.Proxies, s *singbox.SingBoxOut) ([]singbox.SingBoxOut, error) {
	p.Tls = true
	tls(p, s)
	var err error
	s.UpMbps, err = anyToMbps(p.Up)
	if err != nil {
		return nil, fmt.Errorf("hysteia2: %w", err)
	}
	s.DownMbps, err = anyToMbps(p.Down)
	if err != nil {
		return nil, fmt.Errorf("hysteia2: %w", err)
	}
	s.Password = p.Password
	if p.ObfsPassword != "" {
		s.Obfs = &singbox.SingObfs{
			Type:  p.Obfs,
			Value: p.ObfsPassword,
		}
	}
	return []singbox.SingBoxOut{*s}, nil
}

var rateStringRegexp = regexp.MustCompile(`^(\d+)\s*([KMGT]?)([Bb])ps$`)

func anyToMbps(s string) (int, error) {
	if s == "" {
		return 0, nil
	}

	if mb, err := strconv.Atoi(s); err == nil {
		return mb, nil
	}

	m := rateStringRegexp.FindStringSubmatch(s)
	if m == nil {
		return 0, fmt.Errorf("anyToMbps: %w", ErrNotSupportType)
	}

	n := 1.0
	switch m[2] {
	case "K":
		n = 1.0 / 1000.0
	case "M":
		n = 1
	case "G":
		n = 1000
	case "T":
		n = 1000 * 1000
	}
	if m[3] == "B" {
		n = n * 8.0
	}
	v, err := strconv.Atoi(m[1])
	if err != nil {
		return 0, fmt.Errorf("anyToMbps: %w", ErrNotSupportType)
	}
	mb := int(float64(v) * n)
	if mb == 0 {
		mb = 1
	}
	return mb, nil
}
