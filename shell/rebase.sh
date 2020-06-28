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
git checkout $NOWB
git rebase $b

echo "$prefix rebase finish"

echo "$prefix 接下来的合并步骤"
echo "      1、将冲突文件一一处理"
echo "      2、处理完后，运行 git add ."
echo "      3、运行 git rebase --continue"
echo "      4、如果想放弃 rebase，可以运行 git rebase --abort"