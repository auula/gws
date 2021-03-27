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
	"net/http"
	"time"
)

var (
	// Memory storage type config
	DefaultCfg = &Configs{
		TimeOut: time.Minute * 30,
		Cookie: &http.Cookie{
			Name:     SessionKey,
			Path:     "/",
			Secure:   false,
			HttpOnly: true,
		},
	}
)

// Configs session option
type Configs struct {
	Cookie *http.Cookie

	// Domain string `json:"domain" validate:"required"`

	// Path string `json:"Path" validate:"required"`

	// Secure string `json:"secure" validate:"required"`

	// HttpOnly bool `json:"http_only" validate:"required"`

	// sessionID value encryption key
	//EncryptedKey string `json:"encrypted_key" validate:"required,len=16"`

	// redis server ip
	RedisAddr string `json:"redis_addr" validate:"required"`
	// redis auth password
	RedisPassword string `json:"redis_password" validate:"required"`
	// redis key prefix
	RedisKeyPrefix string `json:"redis_key_prefix" validate:"required"`
	// redis db
	RedisDB int `json:"redis_db" validate:"gte=0,lte=15"`
	// the life cycle of a session without operations
	TimeOut time.Duration `json:"time_out" validate:"required"`
	// connection pool size
	PoolSize uint8 `json:"pool_size" validate:"gte=5,lte=100"`
}
