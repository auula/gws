// MIT License

// Copyright (c) 2022 Leon Ding

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package gws

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"
)

type (
	store uint8
)

const (
	ram store = iota // session storage ram type
	rds              // session storage rds type
	customize
	prefix   = "gws_id"
	lifeTime = time.Duration(1800) * time.Second
)

var (
	// default option
	defaultOption = option{
		LifeTime:   lifeTime,
		CookieName: prefix,
		Domain:     "",
		Path:       "/",
		HttpOnly:   true,
		Secure:     true,
	}

	// DefaultRAMOption default RAM config parameter option.
	DefaultRAMOptions = &RAMOption{
		option: defaultOption,
	}

	// NewRDSOption default RDS config parameter option.
	NewRDSOptions = func(ip string, port uint16, passwd string, opts ...func(*RDSOption)) *RDSOption {
		var rdsopt RDSOption
		rdsopt.option = defaultOption

		rdsopt.Index = 6
		rdsopt.Prefix = prefix
		rdsopt.PoolSize = 10
		rdsopt.Password = passwd
		rdsopt.Address = fmt.Sprintf("%s:%v", ip, port)

		for _, opt := range opts {
			opt(&rdsopt)
		}

		return &rdsopt
	}

	// Index set redis database number
	Index = func(number uint8) func(*RDSOption) {
		return func(r *RDSOption) {
			r.Index = number
		}
	}

	// PoolSize set redis connection  pool size
	PoolSize = func(poolsize uint8) func(*RDSOption) {
		return func(r *RDSOption) {
			r.PoolSize = poolsize
		}
	}

	// Prefix set redis key prefix
	Prefix = func(prefix string) func(*RDSOption) {
		return func(r *RDSOption) {
			r.Prefix = prefix
		}
	}

	// Opts set base option
	Opts = func(opt Options) func(*RDSOption) {
		return func(r *RDSOption) {
			r.option = opt.option
		}
	}
)

// option type is default config parameter option.
type option struct {
	LifeTime   time.Duration `json:"life_time,omitempty"`
	CookieName string        `json:"cookie_name,omitempty"`
	Path       string        `json:"path,omitempty"`
	HttpOnly   bool          `json:"http_only,omitempty"`
	Secure     bool          `json:"secure,omitempty"`
	Domain     string        `json:"domain,omitempty"`
}

type Options struct {
	option
}

var (
	LifeTime = func(d time.Duration) func(*Options) {
		return func(o *Options) {
			o.LifeTime = d
		}
	}
	CookieName = func(cn string) func(*Options) {
		return func(o *Options) {
			o.CookieName = cn
		}
	}
	Path = func(path string) func(*Options) {
		return func(o *Options) {
			o.Path = path
		}
	}
	HttpOnly = func(b bool) func(*Options) {
		return func(o *Options) {
			o.HttpOnly = b
		}
	}
	Secure = func(b bool) func(*Options) {
		return func(o *Options) {
			o.Secure = b
		}
	}
	Domain = func(domain string) func(*Options) {
		return func(o *Options) {
			o.Domain = domain
		}
	}
)

func NewOptions(opts ...func(*Options)) Options {
	var opt Options
	opt.option = defaultOption

	for _, v := range opts {
		v(&opt)
	}

	return opt
}

// RAMOption is RAM storage config parameter option.
type RAMOption struct {
	option
}

// RDSOption is Redis storage config parameter option.
type RDSOption struct {
	option
	Index    uint8  `json:"db_index,omitempty"`
	Prefix   string `json:"prefix,omitempty"`
	Address  string `json:"address,omitempty"`
	Password string `json:"password,omitempty"`
	PoolSize uint8  `json:"pool_size,omitempty"`
}

// Configure is session storage config parameter parser.
type Configure interface {
	Parse() (cfg *config)
}

// config is session storage config parameter.
type config struct {
	store `json:"store,omitempty"`
	RDSOption
}

func (opt Options) Parse() (cfg *config) {
	cfg = new(config)
	cfg.store = customize
	// 默认本机内存存储，只需要设置基本设置即可
	cfg.RDSOption.option = opt.option
	return verifyCfg(cfg)
}

func (opt RAMOption) Parse() (cfg *config) {
	cfg = new(config)
	cfg.store = ram
	// 默认本机内存存储，只需要设置基本设置即可
	cfg.RDSOption.option = opt.option
	return verifyCfg(cfg)
}

func (opt RDSOption) Parse() (cfg *config) {
	cfg = new(config)
	cfg.store = rds
	// redis存储相应的设置就会多一点，校验策略根据redis策略
	cfg.RDSOption = opt
	return verifyCfg(cfg)
}

func verifyCfg(cfg *config) *config {

	// 通用校验
	if cfg.CookieName == "" {
		panic("cookie name is empty.")
	}
	if cfg.Path == "" {
		panic("domain path is empty.")
	}
	if cfg.LifeTime <= 0 {
		cfg.LifeTime = lifeTime
	}

	// ram校验通过直接返回
	if cfg.store == ram || cfg.store == customize {
		return cfg
	}

	if cfg.Index > 16 {
		cfg.Index = 6
	}

	if cfg.PoolSize <= 0 {
		cfg.PoolSize = 10
	}

	if cfg.Prefix == "" {
		cfg.Prefix = prefix
	}

	if cfg.Password == "" {
		panic("remote server login passwd is empty.")
	}

	// 针对特定存储校验
	if net.ParseIP(strings.Split(cfg.Address, ":")[0]) == nil {
		panic("remote ip address illegal.")
	}
	if matched, err := regexp.MatchString("^[0-9]*$", strings.Split(cfg.Address, ":")[1]); err == nil {
		if !matched {
			panic("remote server port illegal.")
		}
	}
	return cfg
}
