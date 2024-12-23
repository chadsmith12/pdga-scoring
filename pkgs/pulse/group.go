package pulse

import "fmt"

type Group struct {
    prefix string
    router *PulseRouter
    middleware []MiddlewareFunc
}

func NewGroup(prefix string, router *PulseRouter) *Group {
    return &Group{
        prefix: prefix,
        router: router,
    }
}

func (group *Group) Get(pattern string, endpoint EndpointHandler) {
    group.router.Get(routeUrl(group.prefix, pattern), endpoint, group.middleware...)
}

func (group *Group) Post(pattern string, endpoint EndpointHandler) {
    group.router.Post(routeUrl(group.prefix, pattern), endpoint, group.middleware...)
}

func (group *Group) Use(middleware ...MiddlewareFunc) {
    group.middleware = append(group.middleware, middleware...)
}

func routeUrl(prefix, pattern string) string {
    return fmt.Sprintf("%s%s", prefix, pattern)
}
