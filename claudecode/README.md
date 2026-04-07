
in your ~/.zshrc or ~/.bashrc add the following:

```bash
alias cc="claude --dangerously-skip-permissions $@"  
```

then when it start you can switch modes:  plan / accept edits on / bypass permissions on  

if you are new to claude code. do not use the flag   --dangerously-skip-permissions

youtube links that can help you understand more claude code : 

https://www.youtube.com/watch?v=TiNpzxoBPz0&t=360s

https://www.youtube.com/watch?v=-O6MEtleOdA

https://www.youtube.com/watch?v=rVEoyx349Hk


==== to add bottom status lines ====

in your terminal type: npx ccstatusline@latest

2 lines reccomendation: 

first line: Model
            Seperator
            Context %   ( most importent information )
            Seperator
            Session Cost
            Seperator
            Session Clock

second line: Git branch 
             Separator
             Git Worktree

====================================

install warp terminal from WARP.DEV 
  it will help you navigate files also in the left area 
  for example to see the superpower plans..

===================================

to continue work on your claude code home terminal from your phone

install https://happy.engineering

===========================

Get Shit Done - Meta-prompting framework.

https://github.com/gsd-build/get-shit-done

npx get-shit-done-cc@latest

then 

claude --dangerously-skip-permissions

=========================================

how to connect to CLI tools that make claude code unstoppable

https://www.youtube.com/watch?v=uULvhQrKB_c

need to check browser cli: vercel browser automation vs playright cli - and need to check if this working in wsl2 ubuntu 

https://github.com/microsoft/playwright-cli


https://www.youtube.com/watch?v=P7JrP57AxR0    claude code agent-browser (by vercel) 

======================================================

claude code 21 hidden settings

https://www.youtube.com/watch?v=pDoBe4qbFPE

one of them is how to help claude code read files bigger then 2000 lines

start write prompt ... then ctrl + S  to stash it.. then give him another prompt. after ENTER the stashed prompt return

============================================

/insights  will create for you report.html  ~\.claude\usage-data\reports.html

/effort auto or low ...etc  how many tokens will use per task

 create agents and claude team for these agents... they will know how to work together.

make sure in settings.json :  "CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS": "1",

/simplify  simple my code

=============================================

https://www.youtube.com/watch?v=W9igiY2JdHA
   
claude --channels plugin:telegram@claude-plugins-official
https://www.youtube.com/watch?v=GjDRlqmfoT8
https://bun.sh/
curl -fsSL https://bun.sh/install | bash

pip install git+https://github.com/RichardAtCT/claude-code-telegram@latest
Package 'claude-code-telegram' requires a different Python: 3.10.14 not in '>=3.11'
============================================

https://github.com/oraios/serena
  like elastic 

https://skillsmp.com/
  skill market

==============================
Dan Github with his CLAUDE.md : 

 https://github.com/85danf/agent-skills/blob/master/claude/CLAUDE.md
===========================
   claude code history explorer 

   https://github.com/jhlee0409/claude-code-history-viewer?tab=readme-ov-file
  curl -fsSL https://raw.githubusercontent.com/jhlee0409/claude-code-history-viewer/main/install-server.sh | sh
  
  sudo apt-get install -y libwebkit2gtk-4.1-0

  cchv-server --serve

 _cchv(){
   pkill cchv
   cd
   (cchv-server --serve) & 
   cd .claude-history-viewer
   sleep 5
   token=$(cat webui-token.txt)
   echo "http://localhost:3727?token=$token"
 }
 alias cchv="_cchv"


=======================================

after claude code cracked .
  https://www.youtube.com/watch?v=mBHRPeg8zPU
  www.ccunpacked.dev
  https://github.com/ultraworkers/claw-code

===========================================

claude code self improve https://www.youtube.com/watch?v=wQ0duoTeAAU
 
 ===================================

 cli utils that can help alongside Claude-code

 https://www.youtube.com/watch?v=3NzCBIcIqD0
   Lazygit
   glow ( for md files )
   llmfit   ( table what model i can run in my hardware )
   models ( list models  providers and how much thet charge )
   taproom    what brew package i installed
   ranger  file explorer
   z jump to folder
   btop  https://github.com/aristocratos/btop   system monitor
   mactop
   chafa file.png    you will see the image in cli 
   csvlens file.csv
   
=====================================================
   https://www.youtube.com/watch?v=AhXfI1rSUPc
 5 skills / plugins must use 

 remotion  - motion graphics and videos - using simple prompts
    npx skills add remotion-dev/skills
     https://www.youtube.com/watch?v=7OR-L0AySn8
     mkdir v1
     cd v1
     bun create video    .... then claude  ...then prompt that will use remotion skill
 ==========================================
 claw code : https://www.youtube.com/watch?v=YJR8TfpPUrs
    https://github.com/ultraworkers/claw-code
    https://github.com/Yeachan-Heo/oh-my-claudecode

=======================================

https://www.linkedin.com/feed/update/urn:li:activity:7446945416140619777/
  let openclaw use claude code in remote control
======================================

