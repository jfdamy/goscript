package goscript

import (
	"bufio"
	"io"
	"log"
	"os/exec"
)

// Pipeline represent a pipeline of Pipeable elements
type Pipeline struct {
	commands []Pipeable
}

//Pipe add a new Pipeable element to the pipeline
func (p *Pipeline) Pipe(pipeable Pipeable) *Pipeline {
	p.commands = append(p.commands, pipeable)
	return p
}

//Run run the all pipeline and return the output chan
func (p *Pipeline) Run() (chan (string), error) {
	var output chan (string)
	var input chan (string)
	for _, command := range p.commands {
		output = make(chan string)
		go command.Exec(input, output)
		input = output
	}
	return output, nil
}

// Pipeable represent a executable that can be "pipe" into a pipeline
type Pipeable interface {

	//Exec execute the executable
	Exec(input, output chan (string))

	//Pipe add the element to a new Pipeline
	Pipe(pipeable Pipeable) *Pipeline
}

type cmd struct {
	command *exec.Cmd
}

//Exec execute the executable
func (cmd *cmd) Exec(input, output chan (string)) {
	stdout, _ := cmd.command.StdoutPipe()
	if input != nil {
		stdin, _ := cmd.command.StdinPipe()
		go func() {
			defer stdin.Close()
			for i := range input {
				io.WriteString(stdin, i)
			}
		}()
	}
	go func() {
		defer stdout.Close()
		rd := bufio.NewReader(stdout)
		for {
			str, err := rd.ReadString('\n')
			if err != nil {
				if output != nil {
					close(output)
				}
				return
			}
			output <- str
		}
	}()
	err := cmd.command.Run()
	if err != nil {
		log.Println("Command : ", cmd.command.Path, " Error :", err)
	}
}

//Pipe create a new pipeline with this Pipeable element
func (cmd *cmd) Pipe(pipeable Pipeable) *Pipeline {
	p := new(Pipeline)
	p.commands = []Pipeable{cmd, pipeable}
	return p
}

// Command instantiate a Pipeable shell command
func Command(name string, arg ...string) Pipeable {
	cmd := new(cmd)
	cmd.command = exec.Command(name, arg...)
	return cmd
}

type funcPipeable struct {
	exec func(input, output chan (string))
}

//Exec execute the function
func (f *funcPipeable) Exec(input, output chan (string)) {
	f.exec(input, output)
}

//Pipe create a new pipeline with this Pipeable element
func (f *funcPipeable) Pipe(pipeable Pipeable) *Pipeline {
	p := new(Pipeline)
	p.commands = []Pipeable{f}
	return p
}

//Func instanciate a Pipeable function
func Func(myExec func(input, output chan (string))) Pipeable {
	exec := new(funcPipeable)
	exec.exec = myExec
	return exec
}
