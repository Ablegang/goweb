master:
	
	@echo "============================>提交 master 分支代码"

	git add .
	git commit -m "$(m)"
	git pull origin master
	git push origin master