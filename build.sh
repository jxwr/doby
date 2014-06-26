#!/bin/sh

cd parser
for y in *.y
do
    go tool yacc -o ${y%.y}.go -p Calc $y 
done
cd ..

go build
