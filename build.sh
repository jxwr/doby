#!/bin/sh

cd parser
for y in *.y
do
    goyacc -o ${y%.y}.go -p Doby $y 
done
cd ..

go build
