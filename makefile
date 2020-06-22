# 默认值
MSG := empty msg
ORIGIN := origin
BRANCH := master

ifdef m
	MSG = $(m)
endif

ifdef o
	ORIGIN = $(o)
endif

ifdef b
	BRANCH = $(b)
endif

master:
	
	@echo "现在开始提交，消息：$(MSG)，仓库：$(ORIGIN)，分支：$(BRANCH)"

	git add .
	git commit -m "$(MSG)"
	git pull $(ORIGIN) $(BRANCH)
	git push $(ORIGIN) $(BRANCH)