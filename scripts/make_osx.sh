#!/bin/bash

# build
echo "Building..."
go build main.go

# move the created 'main' binary into 'bin/step_agent_osx'
echo "Moving to bin..."
mkdir -p bin
mv main bin/step_agent_osx

echo "Done"