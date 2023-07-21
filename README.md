<p align="center">
    <img style="width:20em;" src="./assets/mascott.png" alt="mascott.png">
</p>

# todo

A simple `TUI` app built in golang to handle your todo-s


<img style="width:30em;" src="./assets/demo.gif" alt="jim">

## Configuration

`todo` doesn't create the config file for you, but it looks in the following location:

```shell
$HOME/.config/todo/config.yml
```

this is the default configuration:

```yml
# keybindings for editing
bubble:
  quit: q
  up: [up, j]
  down: [down, k]
  help: ?
  check: space
  insert: enter
  remove: backspace
  edit: right
  editexit: esc
  swapup: shift+up
  swapdown: shift+down
  undo: ctrl+z
```

## Change Log

- `02/07/2023` : multiple keybindings can be associated with the same action
- `04/07/2023` : ctrl+z
- `14/07/2023` : 🗑 bin
- `21/07/2023` : `tui` improvements

 