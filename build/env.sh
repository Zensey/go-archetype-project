#!/bin/sh

set -e

if [ ! -f "build/env.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

# Create fake Go workspace if it doesn't exist yet.
workspace="$PWD/build/_workspace"
root="$PWD"
root_name="go-archetype-project"
sdir="$workspace/src/dev.rubetek.com"

if [ ! -L "$sdir/$root_name" ]; then
    mkdir -p "$sdir"
    cd "$sdir"
    ln -s ../../../../. "$root_name"
    cd "$root"
fi

# Set up the environment to use the workspace.
export GOPATH="$workspace"
export PATH="$PATH:$workspace/bin"

# Run the command inside the workspace.
cd "$sdir/$root_name"
PWD="$sdir/$root_name"

# Launch the arguments with the configured environment.
exec "$@"
