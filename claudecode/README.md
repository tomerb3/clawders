
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

======================================================

claude code 21 hidden settings

https://www.youtube.com/watch?v=pDoBe4qbFPE

one of them is how to help claude code read files bigger then 2000 lines

start write prompt ... then ctrl + S  to stash it.. then give him another prompt. after ENTER the stashed prompt return

============================================

/insights  will create for you report.html  ~\.claude\usage-data\reports.html
