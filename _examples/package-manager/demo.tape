# go run github.com/charmbracelet/vhs@latest demo.tape
Output demo.gif
Require go

Set Theme "Dracula+"
Set FontFamily "JetBrainsMono Nerd Font Mono"
Set FontSize 18
Set Width 1100
Set Height 500
Set Margin 20
Set Padding 40
Set MarginFill "#674EFF"
Set BorderRadius 10
Env PS1 "$ "

Type "go run *.go"
Sleep 500ms
Enter
Sleep 10s
