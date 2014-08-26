#!/bin/sh

cd parser
for y in *.y
do
    go tool yacc -o ${y%.y}.go -p Doby $y 
done
cd ..

go build
