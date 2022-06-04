# goboxes
A BYOS (bring your own scripts) plain text status generator written in Go. No built-in modules, no fancy features, no X11 or Wayland dependencies, just a tiny Go binary that manages executing your scripts at varying intervals and prints their output back to stdout, either delimited or according to your own formatting string!

goboxes is designed to be used with status bars that take in plain text like Lemonbar or dwm, but can easily be used with anything else that interacts with plain text!