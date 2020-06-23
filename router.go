package rester

import (
	"reflect"

	"github.com/buaazp/fasthttprouter"
	"github.com/henrylee2cn/ameda"
)

type Router struct {
	router          fasthttprouter.Router
	controllers     map[string]Controller // {relativePath:Controller}
	controllerNames map[string]string     // {controllerName:relativePath}
}

// Control routes controller.
// NOTE:
// The same routing controller can be registered repeatedly, but only for the first time;
// If the controller of the same route registered twice is different, panic
func (r *Router) Control(path string, controller Controller) {
	if r.controllers == nil {
		r.controllers = make(map[string]Controller)
	}
	if r.controllerNames == nil {
		r.controllerNames = make(map[string]string)
	}
	ctl, ok := r.controllers[path]
	if ok && reflect.TypeOf(ctl) == reflect.TypeOf(controller) {
		return
	}
	handlerMap := MustToHandlers(controller)
	controllerName := getControllerName(controller)
	for _, method := range httpMethodList {
		handler := handlerMap[method]
		if handler != nil {
			r.router.Handle(method, path, handler)
			r.controllerNames[controllerName] = path
		}
	}
	r.controllers[path] = controller
}

// Path returns router path of the controller
// NOTE:
//  Must be called after routing
func (r *Router) Path(controller Controller) string {
	return r.controllerNames[getControllerName(controller)]
}

func getControllerName(controller Controller) string {
	t := ameda.DereferenceValue(reflect.ValueOf(controller)).Type()
	return t.PkgPath() + "." + t.Name()
}