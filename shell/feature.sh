prefix="===========>"

if [ "$b" = "" ]; then
  echo "$prefix 请指定新分支名"
  exit
fi

if [ "$(git status | grep 'nothing to commit, working tree clean')" = "" ]; then
  echo "$prefix 当前分支有修改暂未提交，请先处理"
  exit
fi

echo "$prefix 即将创建 feature/$b 分支"

from=$f

if [ "$from" = "" ]; then
  echo "$prefix 当前未指定基础分支，默认从 master 切出新分支"
  from="master"
fi