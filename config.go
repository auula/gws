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
	ram store = iota // Session storage ram type
	rds              // Session storage rds type
	def

	prefix   = "gws_id"                          // Default prefix
	lifeTime = time.Duration(1800) * time.Second // Default session lifetime
)

var (

	// Default option
	defaultOption = option{
		LifeTime:   lifeTime,
		CookieName: prefix,
		Domain:     "",
		Path:       "/",
		HttpOnly:   true,
		Secure:     true,
	}

	// DefaultRAMOptions default RAM config parameter option.
	DefaultRAMOptions = &RAMOption{
		option: defaultOption,
	}

	// NewRDSOptions default RDS config parameter option.
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

	// WithIndex set redis database number
	WithIndex = func(number uint8) func(*RDSOption) {
		return func(r *RDSOption) {
			r.Index = number
		}
	}

	// WithPoolSize set redis connection  pool size
	WithPoolSize = func(poolSize uint8) func(*RDSOption) {
		return func(r *RDSOption) {
			r.PoolSize = poolSize
		}
	}

	// WithPrefix set redis key prefix
	WithPrefix = func(prefix string) func(*RDSOption) {
		return func(r *RDSOption) {
			r.Prefix = prefix
		}
	}

	// WithOpts set base option
	WithOpts = func(opt Options) func(*RDSOption) {
		return func(r *RDSOption) {
			r.option = opt.option
		}
	}
)

// option type is default config parameter option.
type option struct {
	LifeTime   time.Duration `json:"life_time"`
	CookieName string        `json:"cookie_name"`
	HttpOnly   bool          `json:"http_only"`
	Path       string        `json:"path"`
	Secure     bool          `json:"secure"`
	Domain     string        `json:"domain"`
}

// Options type is default config parameter option.
type Options struct {
	option
}

var (
	WithLifeTime = func(d time.Duration) func(*Options) {
		return func(o *Options) {
			o.LifeTime = d
		}
	}
	WithCookieName = func(cn string) func(*Options) {
		return func(o *Options) {
			o.CookieName = cn
		}
	}
	WithPath = func(path string) func(*Options) {
		return func(o *Options) {
			o.Path = path
		}
	}
	WithHttpOnly = func(b bool) func(*Options) {
		return func(o *Options) {
			o.HttpOnly = b
		}
	}
	WithSecure = func(b bool) func(*Options) {
		return func(o *Options) {
			o.Secure = b
		}
	}
	WithDomain = func(domain string) func(*Options) {
		return func(o *Options) {
			o.Domain = domain
		}
	}
)

// NewOptions Initialize default config.
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
	Index    uint8  `json:"db_index" `
	Prefix   string `json:"prefix" `
	Address  string `json:"address" `
	Password string `json:"password" `
	PoolSize uint8  `json:"pool_size" `
}

// Configure is session storage config parameter parser.
type Configure interface {
	Parse() (cfg *Config)
}

// Config is session storage config parameter.
type Config struct {
	store `json:"store,omitempty"`
	*RDSOption
}

func (opt *Options) Parse() (cfg *Config) {
	cfg = new(Config)
	cfg.store = def
	cfg.RDSOption.option = opt.option
	return verifyCfg(cfg)
}

func (opt *RAMOption) Parse() (cfg *Config) {
	cfg = new(Config)
	cfg.store = ram
	cfg.RDSOption = new(RDSOption)
	cfg.RDSOption.option = opt.option
	return verifyCfg(cfg)
}

func (opt *RDSOption) Parse() (cfg *Config) {
	cfg = new(Config)
	cfg.store = rds
	cfg.RDSOption = opt
	return verifyCfg(cfg)
}

// Check the data
func verifyCfg(cfg *Config) *Config {
	// General check
	if cfg.CookieName == "" {
		panic("cookie name is empty.")
	}
	if cfg.Path == "" {
		panic("domain path is empty.")
	}
	if cfg.LifeTime <= 0 {
		cfg.LifeTime = lifeTime
	}

	if cfg.store == ram || cfg.store == def {
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

	// Verification for specific storage
	if net.ParseIP(strings.Split(cfg.Address, ":")[0]) == nil {
		panic("remote ip address illegal.")
	}
	if matched, err := regexp.MatchString("^[0-9]*$", strings.Split(cfg.Address, ":")[1]); err == nil {
		if !matched {
			panic("remote server port illegal.")
		}
	}
	debug.trace(cfg)
	return cfg
}
