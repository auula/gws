// MIT License

// Copyright (c) 2021 Jarvib Ding

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

package sessionx

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"
)

var (
	// Memory storage type config
	DefaultCfg = &Configs{
		TimeOut:  time.Minute * 30,
		Domain:   "",
		Path:     "/",
		Name:     SessionKey,
		Secure:   false,
		HttpOnly: true,
	}
)

// Configs session option
type Configs struct {

	// sessionID value encryption key
	//EncryptedKey string `json:"encrypted_key" validate:"required,len=16"`

	// redis server ip
	RedisAddr string `json:"redis_addr" validate:"required"`
	// redis auth password
	RedisPassword string `json:"redis_password" validate:"required"`
	// redis key prefix
	RedisKeyPrefix string `json:"redis_key_prefix" validate:"required"`
	// redis db
	RedisDB int `json:"redis_db" validate:"gte=0,lte=15,redis"`
	// the life cycle of a session without operations
	TimeOut time.Duration `json:"time_out" validate:"required"`
	// connection pool size
	PoolSize uint8 `json:"pool_size" validate:"gte=5,lte=100"`

	// cookie domain
	Domain string `json:"domain" `

	// cookie domain url path
	Path string `json:"Path" validate:"required"`

	// the browser can only use the cookie through a secure encrypted connection
	Secure bool `json:"secure" `

	// not allow javascript operating
	HttpOnly bool `json:"http_only" validate:"required"`

	// cookie key name
	Name string `json:"name" validate:"required"`

	// cookie config template
	_cookie *http.Cookie
}

// VerifyRedis config parameter
func (c *Configs) VerifyRedis() error {
	IPAndPortRegex := "^(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])\\.(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])\\.(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])\\.(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5]):([0-9]|[1-9]\\d|[1-9]\\d{2}|[1-9]\\d{3}|[1-5]\\d{4}|6[0-4]\\d{3}|65[0-4]\\d{2}|655[0-2]\\d|6553[0-5])$"
	if err := c.VerifyMemory(); err != nil {
		return fmt.Errorf("verify config param fail: %s", err.Error())
	}
	// check ip address
	match, _ := regexp.MatchString(IPAndPortRegex, c.RedisAddr)
	if !match {
		return errors.New("redis ip address format error")
	}
	// redis redis key prefix
	if len(c.RedisKeyPrefix) <= 0 {
		return errors.New("redis key prefix required")
	}
	// check redis database number index
	if c.RedisDB < 0 || c.RedisDB > 15 {
		return errors.New("redis database number index not exist")
	}
	// check redis pool size
	if c.PoolSize < 1 || c.PoolSize > 100 {
		return errors.New("redis pool size range is 1 - 100")
	}
	return nil
}

// VerifyMemory config parameter
func (c *Configs) VerifyMemory() error {
	// check time out
	if c.TimeOut <= 0 {
		c.TimeOut = time.Minute * 30
	}
	// check cookie path
	if c.Path == "" {
		return errors.New("cookie path not exist")
	}
	// check cookie key name
	if c.Name == "" {
		return errors.New("cookie key name not exist")
	}
	// no allow javascript operating
	c.HttpOnly = true
	return nil
}

// Parse config parameter
func (c *Configs) Parse() *Configs {
	c._cookie = &http.Cookie{
		Domain:   c.Domain,
		Path:     c.Path,
		Name:     c.Name,
		Secure:   c.Secure,
		HttpOnly: c.HttpOnly,
		Expires:  time.Now().Add(c.TimeOut),
	}
	return c
}

// Reload loading config
func (c *Configs) Reload(cookie *http.Cookie) {
	c._cookie = cookie
}
