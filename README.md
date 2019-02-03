# i3icons2
go version of i3icons - native deamon to rename workspaces with fontawesome-icons based on open applications

## i3 config to switch to a workspace by number, not by name

```

exec_always --no-startup-id killall i3icons && i3icons || i3icons

set $ws1 number 1
set $ws2 number 2
set $ws3 number 3
set $ws4 number 4
set $ws5 number 5
set $ws6 number 6
set $ws7 number 7
set $ws8 number 8
set $ws9 number 9
set $ws0 number 10

# switch to workspace
bindsym $mod+1 workspace $ws1
bindsym $mod+2 workspace $ws2
bindsym $mod+3 workspace $ws3
bindsym $mod+4 workspace $ws4
bindsym $mod+5 workspace $ws5
bindsym $mod+6 workspace $ws6
bindsym $mod+7 workspace $ws7
bindsym $mod+8 workspace $ws8
bindsym $mod+9 workspace $ws9
bindsym $mod+0 workspace $ws0

# move focused container to workspace
bindsym $mod+Shift+1 move container to workspace $ws1
bindsym $mod+Shift+2 move container to workspace $ws2
bindsym $mod+Shift+3 move container to workspace $ws3
bindsym $mod+Shift+4 move container to workspace $ws4
bindsym $mod+Shift+5 move container to workspace $ws5
bindsym $mod+Shift+6 move container to workspace $ws6
bindsym $mod+Shift+7 move container to workspace $ws7
bindsym $mod+Shift+8 move container to workspace $ws8
bindsym $mod+Shift+9 move container to workspace $ws9
bindsym $mod+Shift+0 move container to workspace $ws0
```
