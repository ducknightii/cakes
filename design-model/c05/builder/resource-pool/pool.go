package resource_pool

import "errors"

const (
	defaultTotal   = 8
	defaultMaxIdle = 8
	defaultMinIdle = 0
)

type Pool struct {
	name    string
	total   int
	maxIdle int
	minIdle int

	// some func
}

type Builder struct {
	name    string
	total   int
	maxIdle int
	minIdle int
}

func (b *Builder) SetName(name string) {
	b.name = name
}

func (b *Builder) SetTotal(total int) {
	b.total = total
}

func (b *Builder) SetMaxIdle(num int) {
	b.maxIdle = num
}

func (b *Builder) SetMinIdle(num int) {
	b.minIdle = num
}

func (b *Builder) Build() (pool Pool, err error) {
	// set default
	if b.total == 0 {
		b.total = defaultTotal
	}
	if b.maxIdle == 0 {
		b.maxIdle = defaultMaxIdle
	}
	if b.minIdle == 0 {
		b.minIdle = defaultMinIdle
	}

	// conf check
	if b.name == "" {
		err = errors.New("name empty")

		return
	}

	if b.minIdle > b.maxIdle || b.maxIdle > b.total {
		err = errors.New("invalid conf")

		return
	}

	// new instance
	pool = Pool{
		name:    b.name,
		total:   b.total,
		maxIdle: b.maxIdle,
		minIdle: b.minIdle,
	}

	return
}
