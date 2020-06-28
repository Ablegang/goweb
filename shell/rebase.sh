prefix="===========>"

if [ "$b" = "" ]; then
  echo "$prefix 请指定要作为基础的分支名"
  exit
fi

if [[ $(git status | grep 'nothing to commit, working tree clean') = "" ]]; then
  echo "$prefix 当前分支有修改暂未提交，请先处理"
  exit
fi

echo "当前分支：$(git branch | grep \* | grep -Eo ' .+')"
NOWB=$(git branch | grep \* | grep -Eo ' .+')

git checkout $b
git pull --rebase
git checkout NOWB
git rebase $b
