
# redraw

An interactive terminal tool to filter and redraw output continuously.

I was frustrated with terminal commands like `column`, `sort`, `wc` not
showing any output until they encounter an EOF. So I wrote `redraw` to run
those commands on every new line and redraw the terminal output with 
terminal escape sequences.

With `redraw` you can provide a updates to terminal commands before the input
reaches EOF, resulting in a more interactive experience.

Example usage:

You'd want to format output with `column -t` but the lines are arriving 
slowly. Instead of staring at a blank screen for minutes, you can appreciate
the formatted output as the results are coming in.

Instead of this:

```bash
your_command | column -t
```

Write this:

```bash
your_command | redraw column -t
```

If your downstream pipeline consists of multiple piped commands:

Instead of this:

```bash
your_command | command_a | command_b | command_c
```

You could write this:

```bash
your_command | redraw bash -c 'command_a | command_b | command_c'
```

TODO: we could have a line editor script performing the back and forth conversion.

# Shout outs 

* to Claudio Fahey for writing an inspiring blog post: https://medium.com/statuscode/pipeline-patterns-in-go-a37bb3a7e61d
