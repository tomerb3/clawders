1. create free AI token in NVIDIA website

2. clone this repository 
   https://github.com/zhangrr/claude-nvidia-proxy

3. I did some changes to fix some problems 
   check 'main.go' file 

   also I wanted to run it as docker and not : ' go run .'
   so you will also see Dockerfile, dockerbuild and dockerrun

4. put these files inside the clone repository

5. make sure you ran the docker run of proxy in the background...
    then run these commands 

export ANTHROPIC_BASE_URL=http://localhost:3001

export ANTHROPIC_AUTH_TOKEN=<your nvidia token> >

export ANTHROPIC_DEFAULT_HAIKU_MODEL=qwen/qwen2.5-coder-32b-instruct

export ANTHROPIC_DEFAULT_SONNET_MODEL=qwen/qwen2.5-coder-32b-instruct

export ANTHROPIC_DEFAULT_OPUS_MODEL=qwen/qwen2.5-coder-32b-instruct

export DISABLE_TOOLS=true

claude