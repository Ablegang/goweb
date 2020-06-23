prefix="===========>"

if [ "$b" = "" ]; then
  echo "$prefix 请指定新分支名"
  exit
fi

if [[ $(git status | grep 'nothing to commit, working tree clean') = "" ]]; then
  echo "$prefix 当前分支有修改暂未提交，请先处理"
  exit
fi

from=$f

if [ "$from" = "" ]; then
  echo "$prefix 当前未指定基础分支，默认从 master 切出新分支"
  from="master"
else
  if [[ $(git branch | grep "$from") = "" ]]; then
      echo "$prefix 没有 $from 分支"
      exit
  fi
fi

echo "$prefix 即将从 $from 检出 feature/$b 分支"

git checkout $from > /dev/null
git checkout -b feature/$b > /dev/null
git status

echo "$prefix finish"