# 默认值
MSG := emptymsg
ORIGIN := origin
BRANCH := branch

ifdef m
	MSG = $(m)
endif

ifdef o
	ORIGIN = $(o)
endif

ifdef b
	BRANCH = $(b)

endif

git:
	
	@echo "现在开始提交，消息：$(MSG)，仓库：$(ORIGIN)，分支：$(BRANCH)"

	git add .
	git commit -m "$(MSG)"
	git push $(ORIGIN) $(BRANCH)