master:
# 需要指定当前系统的 sh
# make master m=message
	@sh shell/master.sh

# make feature b=logpkg_0623 f=master
# b 指新建的分支名 f 指基本分支，默认为 master，则从 master 切出新分支
feature:
	@sh shell/feature.sh