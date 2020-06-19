package logrus_docs

// hook 接口，限定了所有 hook 实例必须实现的方法
// 所有添加的 hook ，都会在指定的日志级别 logs 发生时执行
// 比如执行当一个 entry 实例执行 logs 时，会根据该 entry 的 level 或调用 logs 时指定的 level 来执行相关 level 的 hook
// logs.Hooks 的类型是 LevelHooks，则一个 logs 可以绑定多个 level，一个 level 又可以绑定多个 Hook
// 当 entry.logs 时，会调用 LevelHooks.Fire，会顺序地将该级别已注册的 Hook 一一执行
// 这是阻塞调用、顺序调用，如果想要让这些 Hook 并发运行，则需要自己处理
type Hook interface {
	// 一个 hook 可以同时注册给多个日志级别
	Levels() []Level
	// 一个 hook 必须包含一个 Fire 方法，Fire 就是执行具体 Hook 的逻辑
	Fire(*Entry) error
}

// Hooks map 的类型
type LevelHooks map[Level][]Hook

// 给 logs 实例添加 Hook
func (hooks LevelHooks) Add(hook Hook) {
	for _, level := range hook.Levels() {
		hooks[level] = append(hooks[level], hook)
	}
}

// 执行指定日志级别所有已注册的 Hook
func (hooks LevelHooks) Fire(level Level, entry *Entry) error {
	for _, hook := range hooks[level] {
		//fmt.Println(hook.Fire(entry))
		if err := hook.Fire(entry); err != nil {
			return err
		}
	}

	return nil
}
