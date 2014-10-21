#!/bin/bash

# build
echo "Building..."
go build

# move the created 'main' binary into 'bin/step_agent_osx'
echo "Moving to bin..."
mkdir -p _bin
mv step-agent-go _bin/step_agent_osx

echo "Done"