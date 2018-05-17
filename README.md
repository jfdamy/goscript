# GOSCRIPT

Goscript can be use to create small scrit in golang.
It's just for fun, do not use this for anything ;).


Like this : 
```
pipeline := goscript.Command("echo", "this is a main test").
		Pipe(goscript.Func(func(input, output chan (string)) {
			for i := range input {
				fmt.Print("myfunc say : ", i)
				output <- i
			}
			close(output)
		})).
		Pipe(goscript.Command("sed", "s/main/lame/g")).
		Pipe(goscript.Command("grep", "test"))

	output, _ := pipeline.Run()

	for i := range output {
		fmt.Print(i)
	}
```

output : 

```
./main.go 
myfunc say : this is a main test
this is a lame test
```

It works better with [gorun](https://github.com/erning/gorun) and a golang config for [binfmt_misc](https://www.kernel.org/doc/html/v4.14/admin-guide/binfmt-misc.html).

You can gather more information about binfmt_misc and how to make golang your scripting language [here](https://blog.cloudflare.com/using-go-as-a-scripting-language-in-linux/).