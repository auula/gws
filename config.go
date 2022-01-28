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

import "time"

type StoreType uint8

const (
	RAM StoreType = iota
	RDS
)

// option type is default config parameter option.
type option struct {
	LifeTime   time.Duration
	CookieName string
	DomainPath string
	HttpOnly   bool
	Secure     bool
}

// RAMOption is RAM storage config parameter option.
type RAMOption struct {
	option
}

func (ram RAMOption) Parse() *config {
	return nil
}

// RedisOption is Redis storage config parameter option.
type RedisOption struct {
	option
}

func (rds RedisOption) Parse() *config {
	return nil
}

// config is session storage config parameter.
type config struct {
	option
	StoreType
}

// Parser is session storage config parameter parser.
type Parser interface {
	Parse() *config
}
