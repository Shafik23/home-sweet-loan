#!/bin/bash

set -e

remote=root@pollamin01

echo "------ DEPLOYING to $remote ------"

go build

rsync --delete -azP -e ssh ./home-sweet-loan $remote:~/homesweetloan &
rsync --delete -azP -e ssh ./db_migrations $remote:~/homesweetloan &
rsync --delete -azP -e ssh ./*.js $remote:~/homesweetloan &
rsync --delete -azP -e ssh ./*.html $remote:~/homesweetloan &
rsync --delete -azP -e ssh ./*.ico $remote:~/homesweetloan &

wait

ssh $remote "killall home-sweet-loan || true"
ssh $remote "cd ~/homesweetloan && screen -d -m bash -c './home-sweet-loan > output.log 2>&1'"

echo "---------------------"
echo "------ SUCCESS ------"
echo "---------------------"
