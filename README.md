<p align="center">
    <img style="width:20em;" src="./assets/mascott.png" alt="jim">
</p>

# todo

A simple `TUI` app built in golang to handle your todo-s


<img style="width:30em;" src="./assets/demo.gif" alt="jim">

## configuration

`todo` doesn't create the config file for you, but it looks in the following location:

```shell
$HOME/.config/todo/config.yml
```

this is the default configuration:

```yml
# keybindings for editing
bubble:
  quit: q
  help: ?
  check: space
  insert: enter
  remove: backspace
  edit: right
  editexit: esc
  swapup: shift+up
  swapdown: shift+down
```