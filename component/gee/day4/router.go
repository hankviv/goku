package gee

type router struct {
	handlers    map[string]HandlerFunc
	middlewares []HandlerFunc
}

func NewRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
	}
}

func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	key := r.GetMapKey(method, pattern)
	r.handlers[key] = handler
}

func (r *router) GetMapKey(method, pattern string) string {
	return method + "-" + pattern
}

func (r *router) handle(c *Context) {
	key := r.GetMapKey(c.Method, c.Path)

	if handler, ok := r.handlers[key]; ok {
		c.handlers = append(r.middlewares, handler)
	} else {
		c.String(404, "404 NOT FOUND: %s\n", c.Path)
	}
	c.Next()
}
