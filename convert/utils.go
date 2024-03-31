package convert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/ysicing/clash2singbox/model/clash"
	"github.com/ysicing/clash2singbox/model/singbox"
)

func filter(isinclude bool, reg string, sl []string) ([]string, error) {
	r, err := regexp.Compile(reg)
	if err != nil {
		return sl, fmt.Errorf("filter: %w", err)
	}
	return getForList(sl, func(v string) (string, bool) {
		has := r.MatchString(v)
		if has && isinclude {
			return v, true
		}
		if !isinclude && !has {
			return v, true
		}
		return "", false
	}), nil
}

func getForList[K, V any](l []K, check func(K) (V, bool)) []V {
	sl := make([]V, 0, len(l))
	for _, v := range l {
		s, ok := check(v)
		if !ok {
			continue
		}
		sl = append(sl, s)
	}
	return sl
}

func getTags(s []singbox.SingBoxOut) []string {
	return getForList(s, func(v singbox.SingBoxOut) (string, bool) {
		tag := v.Tag
		if tag == "" || v.Ignored {
			return "", false
		}
		return tag, true
	})
}

func Patch(b []byte, s []singbox.SingBoxOut, include, exclude string, extOut []interface{}, extags ...string) ([]byte, error) {
	d, err := PatchMap(b, s, include, exclude, extOut, extags, true)
	if err != nil {
		return nil, fmt.Errorf("Patch: %w", err)
	}
	bw := &bytes.Buffer{}
	jw := json.NewEncoder(bw)
	jw.SetIndent("", "    ")
	err = jw.Encode(d)
	if err != nil {
		return nil, fmt.Errorf("Patch: %w", err)
	}
	return bw.Bytes(), nil
}

func ToInsecure(c *clash.Clash) {
	for i := range c.Proxies {
		p := c.Proxies[i]
		p.SkipCertVerify = true
		c.Proxies[i] = p
	}
}

func PatchMap(
	tpl []byte,
	s []singbox.SingBoxOut,
	include, exclude string,
	extOut []interface{},
	extags []string,
	urltestOut bool,
) (map[string]any, error) {
	d := map[string]interface{}{}
	err := json.Unmarshal(tpl, &d)
	if err != nil {
		return nil, fmt.Errorf("PatchMap: %w", err)
	}
	tags := getTags(s)

	tags = append(tags, extags...)

	ftags := tags
	if include != "" {
		ftags, err = filter(true, include, ftags)
		if err != nil {
			return nil, fmt.Errorf("PatchMap: %w", err)
		}
	}
	if exclude != "" {
		ftags, err = filter(false, exclude, ftags)
		if err != nil {
			return nil, fmt.Errorf("PatchMap: %w", err)
		}
	}

	if urltestOut {
		s = append([]singbox.SingBoxOut{{
			Type:      "selector",
			Tag:       "select",
			Outbounds: append([]string{"urltest"}, tags...),
			Default:   "urltest",
		}}, s...)
		s = append(s, singbox.SingBoxOut{
			Type:      "urltest",
			Tag:       "urltest",
			Outbounds: ftags,
		})
	}
	anyList := make([]any, 0, len(s)+len(extOut))
	for _, v := range s {
		anyList = append(anyList, v)
	}
	anyList = append(anyList, extOut...)

	d["outbounds"] = anyList

	return d, nil
}

func tls(p *clash.Proxies, s *singbox.SingBoxOut) {
	if p.Tls {
		s.TLS = &singbox.SingTLS{}
		s.TLS.Enabled = bool(p.Tls)
		if p.Servername != "" {
			s.TLS.ServerName = p.Servername
		} else if p.Sni != "" {
			s.TLS.ServerName = p.Sni
		} else {
			s.TLS.ServerName = p.Server
		}
		if p.Fingerprint != "" || p.ClientFingerprint != "" {
			s.TLS.Utls = &singbox.SingUtls{}
			s.TLS.Utls.Enabled = true
			if p.ClientFingerprint != "" {
				s.TLS.Utls.Fingerprint = p.ClientFingerprint
			} else {
				s.TLS.Utls.Fingerprint = p.Fingerprint
			}
		}
		s.TLS.Insecure = bool(p.SkipCertVerify)
		s.TLS.Alpn = p.Alpn
	}
}
