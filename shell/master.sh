# 分支判断
if [ $(git branch | grep \* | grep -Eo ' .+') != "master" ]; then
  echo "================>当前分支不是 master，请切换到 master 分支再执行当前命令";
  exit
fi

git add .
git commit -m "$m"
git pull origin master
git push origin master