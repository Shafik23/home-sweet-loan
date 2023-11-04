#!/bin/bash

set -e

remote=root@pollamin01

echo "------ DEPLOYING to $remote ------"

go build

rsync --delete -azP -e ssh ./home-sweet-loan $remote:~/homesweetloan &
rsync --delete -azP -e ssh ./*.js $remote:~/homesweetloan &
rsync --delete -azP -e ssh ./*.html $remote:~/homesweetloan &
rsync --delete -azP -e ssh ./*.ico $remote:~/homesweetloan &

wait

ssh $remote "killall home-sweet-loan && cd ~/homesweetloan && ./home-sweet-loan &"

echo "---------------------"
echo "------ SUCCESS ------"
echo "---------------------"
