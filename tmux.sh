#!/bin/sh
session="docker"
tmux start-server
tmux new-session -d -s $session -n task-session #"vim -S ~/.vim/sessions/kittybusiness"

tmux set -g pane-border-status bottom
#tmux set -g pane-border-format "#P: #{pane_current_command}"

tmux selectp -t 0
tmux send-keys "nano README.md" C-m
tmux splitw -h -p45 "printf '\033]2;api-servce\033\\'; tail -f /var/log/task/api-service.log"
tmux selectp -t 1
tmux splitw -v -p 66 "printf '\033]2;worker-1\033\\'; tail -f /var/log/task/worker1.log"
tmux selectp -t 2
tmux splitw -v  "printf '\033]2;worker-2\033\\'; tail -f /var/log/task/worker2.log"

tmux selectp -t 0
tmux attach-session -t $session
