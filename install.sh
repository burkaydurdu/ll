#!/bin/bash

mkdir tmp-ll

cd ./tmp-ll || exit

curl -LO https://github.com/burkaydurdu/ll/releases/download/v0.0.1/tureng.zip

unzip tureng.zip

mv ./tureng /usr/local/bin

cd ..

rm -rf ./tmp-ll

echo installed successfully