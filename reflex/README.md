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

```golang
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
## License
This software licensed under the 2-clause BSD license:

Copyright 2021 Paul C. Hoehne

Redistribution and use in source and binary forms, with or without modification, 
are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this 
    list of conditions and the following disclaimer.
2. Redistributions in binary form must reproduce the above copyright notice, 
    this list of conditions and the following disclaimer in the documentation 
    and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND 
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED 
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. 
IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, 
INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, 
BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, 
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF 
LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR 
OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF 
THE POSSIBILITY OF SUCH DAMAGE.