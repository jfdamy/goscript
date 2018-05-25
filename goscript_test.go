package goscript

import (
	"testing"
)

//TestPipeline test an exemple of a all pipeline
func TestPipeline(t *testing.T) {
	pipeline := Command("echo", "this is a main test").
		Pipe(Func(func(input, output chan (string)) {
			for i := range input {
				t.Log("myfunc say : ", i, "\n")
				output <- i
			}
			close(output)
		})).
		Pipe(Command("sed", "s/main/lame/g")).
		Pipe(Command("grep", "test"))

	output, _ := pipeline.Run()

	for i := range output {
		t.Log(i)
		if i != "this is a lame test\n" {
			t.Error("Expected this is a lame test, got ", i)
		}
	}

	pipeline = Func(func(input, output chan (string)) {
		output <- "this is a main test\n"
		close(output)
	}).
		Pipe(Command("sed", "s/main/lame/g")).
		Pipe(Command("grep", "test"))

	output, _ = pipeline.Run()

	for i := range output {
		t.Log(i)
		if i != "this is a lame test\n" {
			t.Error("Expected this is a lame test, got ", i)
		}
	}
}
