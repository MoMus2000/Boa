#!/bin/bash

ALIAS_NAME="boa"
ALIAS_COMMAND="python3 main.py"

alias boa="python3 main.py"
# Determine the user's shell and rc file
if [[ $SHELL == *"bash"* ]]; then
    RC_FILE="$HOME/.bashrc"
elif [[ $SHELL == *"zsh"* ]]; then
    RC_FILE="$HOME/.zshrc"
else
    echo "Unsupported shell: $SHELL"
    exit 1
fi

# Check if the alias already exists in the rc file
if grep -Fxq "alias $ALIAS_NAME='$ALIAS_COMMAND'" "$RC_FILE"; then
    echo "Setting up Boa config."
else
    # Append the alias to the rc file
    echo "alias $ALIAS_NAME='$ALIAS_COMMAND'" >> "$RC_FILE"
    echo "Alias '$ALIAS_NAME' added to $RC_FILE."
fi

# Optionally reload the rc file (useful for interactive shells)
if [[ $SHELL == *"bash"* || $SHELL == *"zsh"* ]]; then
    source "$RC_FILE" >/dev/null 2>&1
    echo "Boa is ready to rip."
fi
