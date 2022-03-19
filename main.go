package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/gosuri/uilive"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pipeline(ctx)
}

func pipeline(ctx context.Context) error {
	errList := [](<-chan error){}
	lines, errc := NewStdinChannel(ctx)
	errList = append(errList, errc)
	aggrLines := NewAggregator(ctx, lines)
	filtered, errc := NewProcessSpawner(ctx, aggrLines, os.Args[1], os.Args[2:]...)
	errList = append(errList, errc)
	errc = NewRedraw(ctx, filtered)
	errList = append(errList, errc)
	errc = MergeErrors(errList...)
	for err := range errc {
		return err
	}
	return nil
}

// NewStdinChannel returns a channel that contains lines read
// from stdin.
func NewStdinChannel(ctx context.Context) (<-chan string, <-chan error) {
	out := make(chan string)
	errc := make(chan error)
	go func() {
		defer close(out)
		defer close(errc)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			select {
			case out <- scanner.Text():
			case <-ctx.Done():
				return
			}
		}
	}()
	return out, errc
}

// NewAggregator returns a channel that supplies all lines read so far.
// Whenever a new line flows in through the input channel, the output
// channel will contain all of the lines so far encountered.
func NewAggregator(ctx context.Context, in <-chan string) <-chan []string {
	out := make(chan []string)
	go func() {
		defer close(out)
		collect := []string{}
		for line := range in {
			collect = append(collect, line)
			select {
			case out <- collect:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

// NewProcessSpawner spawns a process with cmd and args and pipes in the lines
// it gets in the input channel. The stdout bytes of the process are put on
// the output channel.
// Each new item in the input channel results in a new process spawned.
func NewProcessSpawner(ctx context.Context, in <-chan []string, cmd string, args ...string) (<-chan []byte, <-chan error) {
	out := make(chan []byte)
	errc := make(chan error)
	go func() {
		defer close(out)
		defer close(errc)
		for lines := range in {
			bytes, err := spawnProcess(lines, cmd, args...)
			if err != nil {
				errc <- err
				return
			}
			select {
			case out <- bytes:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out, errc
}

// NewRedraw consumes a channel of bytes and overwrites the same part of the
// screen with the new bytes that arrive in the channel.
func NewRedraw(ctx context.Context, in <-chan []byte) <-chan error {
	errc := make(chan error)
	go func() {
		defer close(errc)
		writer := uilive.New()
		writer.RefreshInterval = time.Hour
		writer.Start()
		defer writer.Stop()
		for bytes := range in {
			_, err := writer.Write(bytes)
			if err != nil {
				errc <- err
				return
			}
			writer.Flush()
		}
	}()
	return errc
}

func MergeErrors(channels ...<-chan error) <-chan error {
	wg := sync.WaitGroup{}
	out := make(chan error, len(channels))
	output := func(channel <-chan error) {
		for err := range channel {
			out <- err
		}
		wg.Done()
	}
	wg.Add(len(channels))
	for _, channel := range channels {
		go output(channel)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func spawnProcess(lines []string, cmd string, arg ...string) ([]byte, error) {
	command := exec.Command(cmd, arg...)
	stdin, err := command.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := command.StdoutPipe()
	if err != nil {
		return nil, err
	}
	command.Start()
	for _, line := range lines {
		fmt.Fprintln(stdin, line)
	}
	stdin.Close()
	out, err := io.ReadAll(stdout)
	if err != nil {
		return nil, err
	}
	err = command.Wait()
	if err != nil {
		return nil, err
	}
	return out, nil
}
