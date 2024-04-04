package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ysicing/clash2singbox/convert"
	"github.com/ysicing/clash2singbox/model/clash"
	"gopkg.in/yaml.v3"
)

var (
	path     string
	outPath  string
	include  string
	exclude  string
	insecure bool
)

//go:embed config.json.template
var configByte []byte

func init() {
	flag.StringVar(&path, "i", "", "本地 clash 文件")
	flag.StringVar(&outPath, "o", "config.json", "输出文件")
	flag.StringVar(&include, "include", "", "urltest 选择的节点")
	flag.StringVar(&exclude, "exclude", "", "urltest 排除的节点")
	flag.BoolVar(&insecure, "insecure", false, "所有节点不验证证书")
	flag.Parse()
}

func main() {
	c := clash.Clash{}
	var singList []any
	var tags []string
	if path != "" {
		b, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(b, &c)
		if err != nil {
			panic(err)
		}
	} else {
		panic("i 参数不能都为空")
	}

	if insecure {
		convert.ToInsecure(&c)
	}

	s, err := convert.Clash2sing(c)
	if err != nil {
		fmt.Println(err)
	}
	outb, err := os.ReadFile(outPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			outb = configByte
		} else {
			panic(err)
		}
	}
	name := strings.Split(path, ".")[0]
	name = strings.Split(name, "/")[len(strings.Split(name, "/"))-1]
	outb, err = convert.Patch(outb, s, name, include, exclude, singList, tags...)
	if err != nil {
		panic(err)
	}

	os.WriteFile(outPath, outb, 0644)
}
