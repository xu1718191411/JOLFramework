package framework

type Group struct {
	prefix   string
	handlers map[string]*Tree
}

func NewGroup(router *Router, prefix string) *Group {
	return &Group{
		prefix:   prefix,
		handlers: router.handlers,
	}
}

func (g *Group) Get(url string, h func(ctx *JolContext)) {
	g.handlers["GET"].Add(g.prefix+url, h)
}

func (g *Group) Post(url string, h func(ctx *JolContext)) {
	g.handlers["POST"].Add(g.prefix+url, h)
}

func (g *Group) Put(url string, h func(ctx *JolContext)) {
	g.handlers["PUT"].Add(g.prefix+url, h)
}

func (g *Group) Head(url string, h func(ctx *JolContext)) {
	g.handlers["HEAD"].Add(g.prefix+url, h)
}

func (g *Group) Delete(url string, h func(ctx *JolContext)) {
	g.handlers["DELETE"].Add(g.prefix+url, h)
}
