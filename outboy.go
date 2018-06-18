package outboy

import (
	"fmt"
	"github.com/evil-go/fall"
	"net/http"
	"reflect"
	"strings"
)

var (
	rwt = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()
	rt  = reflect.TypeOf((*http.Request)(nil))
)

type regInfo struct {
	controller interface{}
	paths      map[string]string
}

type outboy struct {
	registrations []regInfo
}

var ob = &outboy{}

func (o *outboy) register() {
	for _, v := range o.registrations {
		i := v.controller
		paths := v.paths
		t := reflect.TypeOf(i)
		vi := reflect.ValueOf(i)
		for k, v := range paths {
			if m, ok := t.MethodByName(k); ok {
				f := m.Func
				ft := f.Type()
				if ft.NumIn() != 3 || ft.NumOut() != 0 || ft.In(0) != t || ft.In(1) != rwt || ft.In(2) != rt {
					panic("Invalid method " + f.String())
				}
				fmt.Println("registering endpoint", m.Name, "at", v)
				http.HandleFunc(v, func(rw http.ResponseWriter, req *http.Request) {
					if strings.EqualFold(m.Name[:len(req.Method)], req.Method) {
						f.Call([]reflect.Value{vi, reflect.ValueOf(rw), reflect.ValueOf(req)})
					}
				})
			}
		}
	}
}

var port = ":8080"

func Port(p int) {
	port = fmt.Sprintf(":%d", p)
}

func Register(i interface{}, paths map[string]string) {
	ob.registrations = append(ob.registrations, regInfo{i, paths})
}

func (o *outboy) InitLast() {
	o.register()
	http.ListenAndServe(port, nil)
}

func init() {
	fall.Register(ob)
}
