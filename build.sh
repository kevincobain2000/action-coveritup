#! /bin/bash

cd app/frontend
npm install
npm run build
cd ..

go build main.go
