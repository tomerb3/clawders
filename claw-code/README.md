The best Agent Harness

https://www.youtube.com/live/RpFh0Nc7RvA

cd ~/src

git clone https://github.com/ultraworkers/claw-code.git

cd claw-code/rust

cargo build --workspace
cd target/debug

./claw prompt "summarize this repository"

./claw prompt "create /tmp/1.txt with 30 random numbers. 3 digits each. with amazing icons each line"


in .zshrc:

alias claw="~/src/claw-code/rust/target/debug/claw $@"
export ANTHROPIC_AUTH_TOKEN="sk-cp-youtoken_minimx_for_example"
export ANTHROPIC_BASE_URL="https://api.minimax.io/anthropic"