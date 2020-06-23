master/%:
# 需要指定当前系统的 sh
# make master/message
	m=$*
	echo $(m)
	@sh shell/master.sh m=$*

master:
	@make master/

# make feature b=logpkg_0623 f=master
# b 指新建的分支名 f 指基本分支，默认为 master，则从 master 切出新分支
feature/%:
	@sh shell/feature.sh b=$*