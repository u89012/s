package s

import (
	"fmt"
	"net/http"
	"regexp"
)

type (
	C struct {
		r *http.Request
		w http.ResponseWriter
		m M
	}

	M map[string]interface{}
)

var (
	routes  = map[string]map[string]func(*C){}
	befores = map[*regexp.Regexp]func(*C){}
	afters  = map[*regexp.Regexp]func(*C){}
)

func init() {
	for _, j := range []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS", "TRACE"} {
		routes[j] = map[string]func(*C){}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c := &C{r, w, M{}}
		if f, k := routes[r.Method][r.URL.Path]; k {
			for p, f := range befores {
				if p.MatchString(r.URL.Path) {
					f(c)
				}
			}
			f(c)
			for p, f := range afters {
				if p.MatchString(r.URL.Path) {
					f(c)
				}
			}
		} else {
			http.NotFound(w, r)
		}
	})
}

func Get(path string, f func(*C))     { routes["GET"][path] = f }
func Post(path string, f func(*C))    { routes["POST"][path] = f }
func Put(path string, f func(*C))     { routes["PUT"][path] = f }
func Delete(path string, f func(*C))  { routes["DELETE"][path] = f }
func Patch(path string, f func(*C))   { routes["PATCH"][path] = f }
func Head(path string, f func(*C))    { routes["HEAD"][path] = f }
func Options(path string, f func(*C)) { routes["OPTIONS"][path] = f }
func Trace(path string, f func(*C))   { routes["TRACE"][path] = f }
func Before(pattern string, f func(*C)) {
	befores[regexp.MustCompile(pattern)] = f
}
func After(pattern string, f func(*C)) {
	afters[regexp.MustCompile(pattern)] = f
}

func Serve(port int) {
	panic(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
