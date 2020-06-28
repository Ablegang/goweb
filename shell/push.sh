if [ "$m" = "" ];then
  m="default msg"
fi

git add .
git commit -m "$m"
git push origin $(git branch | grep \* | grep -Eo ' .+')
