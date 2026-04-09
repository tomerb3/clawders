export ZSH="$HOME/.oh-my-zsh"

ZSH_THEME="powerlevel10k/powerlevel10k"
# Powerlevel10k instant prompt (recommended)
if [[ -r "${XDG_CACHE_HOME:-$HOME/.cache}/p10k-instant-prompt-${(%):-%n}.zsh" ]]; then
  source "${XDG_CACHE_HOME:-$HOME/.cache}/p10k-instant-prompt-${(%):-%n}.zsh"
fi
#ZSH_THEME="robbyrussell"
# Example format: plugins=(rails git textmate ruby lighthouse)
# Add wisely, as too many plugins slow down shell startup.
plugins=(git history terraform copybuffer zsh-autosuggestions zsh-syntax-highlighting dirhistory)
source $ZSH/oh-my-zsh.sh
export PATH=$PATH:$HOME/.bin
export PATH=$PATH:$HOME/docker

alias d='/home/user/docker/docker'
export PATH=$PATH:/home/user/bin

export pw=qwer80

export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"  # This loads nvm
[ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"  # This loads nvm bash_completion
service cron status >/dev/null || sudo service cron start > /dev/null 2>&1

export DISPLAY=$(ip route list default | awk '{print $3}'):0
export LIBGL_ALWAYS_INDIRECT=1

export LIBGL_ALWAYS_INDIRECT=1
export DISPLAY=$(awk '/nameserver / {print $2; exit}' /etc/resolv.conf 2>/dev/null):0
export PATH="/snap/docker/current/bin:$PATH"

export GPG_TTY=$(tty)


export DISPLAY=$(ip route list default | awk '{print $3}'):0
export LIBGL_ALWAYS_INDIRECT=1

source ${HOME}/Dropbox/bash/generalFunctions
export SDKMAN_DIR="$HOME/.sdkman"
[[ -s "$HOME/.sdkman/bin/sdkman-init.sh" ]] && source "$HOME/.sdkman/bin/sdkman-init.sh"
export PATH=$PATH:/usr/local/go/bin

export PYENV_ROOT="$HOME/.pyenv"
[[ -d $PYENV_ROOT/bin ]] && export PATH="$PYENV_ROOT/bin:$PATH"
eval "$(pyenv init - bash)"

_cdmkdir()
{
  if [ ! -z $1 ];then 
     mkdir -p $1
    cd $1
  else
	  echo "you need to enter directory name"
  fi

}
alias cdmkdir="_cdmkdir $1" 
export PATH=~/.npm-global/bin:$PATH
nvm use 21
# Added by Hugging Face CLI installer
export PATH="/home/user/.local/bin:$PATH"


# To customize prompt, run `p10k configure` or edit ~/.p10k.zsh.
[[ ! -f ~/.p10k.zsh ]] || source ~/.p10k.zsh
 typeset -g POWERLEVEL9K_INSTANT_PROMPT=quiet

# bun completions
[ -s "/home/user/.bun/_bun" ] && source "/home/user/.bun/_bun"

OLAMA_CONTEXT_LENGTH=64000

# bun
export BUN_INSTALL="$HOME/.bun"
export PATH="$BUN_INSTALL/bin:$PATH"

# Tmux aliases
alias t='tmux'
alias ta='tmux attach -t'
alias td='tmux detach'
alias tl='tmux ls'
alias tkill='tmux kill-server'

# Sesh + Zoxide integration
eval "$(zoxide init zsh)"
alias cd="z"
export SESH_TMUX_BINARY="tmux"
export SESH_SESSION_PATTERN='sed -E "s|'"'"'||g; s|~||g; s|/|-|g; s| |-|g"'
export SESH_TMUX_LAYOUT="even-horizontal"
alias sz='sesh list -z'

# Open 3 panes: left narrow, right stacked
work() {
  session="work"
  tmux kill-session -t $session 2>/dev/null
  tmux new-session -d -s $session -n "work" -c ~/src/name2
  tmux send-keys -t $session "clear;pwd" Enter

  tmux split-window -h -l 40 -t $session -c ~/src/nameclawders
  tmux send-keys -t $session "clear;pwd" Enter

  tmux split-window -v -l 30 -t $session -c ~/src/name1/a_folder

  tmux send-keys -t $session "clear;l general" Enter
  tmux attach -t $session
}

# Eza aliases
#alias ls='eza --icons --group-directories-first --color=always'
alias ll='eza -lah --icons --group-directories-first --color=always'
alias la='eza -a --icons --group-directories-first --color=always'
export ANTHROPIC_AUTH_TOKEN="sk-cp-example"
export ANTHROPIC_BASE_URL="https://api.minimax.io/anthropic"
export ANTHROPIC_MODEL="MiniMax-M2.7"
export ANTHROPIC_SMALL_FAST_MODEL="MiniMax-M2.7"
export ANTHROPIC_DEFAULT_SONNET_MODEL="MiniMax-M2.7"
export ANTHROPIC_DEFAULT_OPUS_MODEL="MiniMax-M2.7"
export ANTHROPIC_DEFAULT_HAIKU_MODEL="MiniMax-M2.7"

eval "$(fzf --zsh)"
export FZF_DEFAULT_COMMAND="fd --hidden --strip-cwd-prefix --exclude .git"
export FDF_CTRL_T_COMMAND="$FZF_DEFAULT_COMMAND"
export FDF_ALT_C_COMMAND="fd --type=d --hidden --strip-cwd-prefix --exclude .git"
_fzf_compgen_path(){
  fd --hidden --exclude .git . "$1"
}

_fzf_compgen_dir(){
  fd --type=d --hidden --exclude .git . "$1"
}
source ~/fzf-git.sh/fzf-git.sh
eval $(thefuck --alias)
eval $(thefuck --alias fuck)
