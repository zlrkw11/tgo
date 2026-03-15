Output demo.gif

Set FontSize 14
Set Width 900
Set Height 600
Set Theme "Catppuccin Mocha"
Set Shell "zsh"

Type "cd ~/projects/tgo && tgo ./tests/..."
Enter
Sleep 5s

# expand package
Enter
Sleep 1s

# scroll through tests
Down
Sleep 0.5s
Down
Sleep 0.5s
Down
Sleep 0.5s
Down
Sleep 0.5s

# expand error detail
Enter
Sleep 2s

# close error
Enter
Sleep 1s

# rerun test
Type "r"
Sleep 4s

# jump to top
Type "g"
Sleep 1s

# jump to bottom
Type "G"
Sleep 1.5s

# quit
Type "q"
Sleep 0.5s
