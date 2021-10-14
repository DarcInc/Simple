# The Reflex
The reflex is not a full dependency injection.
It's more of a dependency assistance framework.
It can handle basic types like `int`, `string`, etc.
It can handle compound types like structs.
It can handle factory methods that construct values.
It can also handle types that will be used to construct values.

It does not build clever graphs or intelligently cache constructed values.
Instead, if you need that fanciness, wrap it with your own logic.
You may still want to wrap its functions so that you don't have to worry about casting.

## Example Usage

```Go
package myPackage

import (
	"github.com/darcinc/Simple/reflex"
	"net/http"
)

func main() {
	reflex :=  reflex.NewReflex()
    reflex.Register("Message", "Say Hello")
	reflex.Register("Times", 5)
	
	http.Handle("/", MyHandler{R: reflex})
}

type MyHandler struct {
	R reflex.Reflex
}

func (mh MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	message, _ := mh.R.MustGet("Message").(string)
	repeat, _ := mh.R.MustGet("Times").(int)
	
	for i := 0; i < repeat; i++ {
		w.Write([]byte(message))
    }
}

```