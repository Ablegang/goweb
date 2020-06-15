package logob

import "testing"

func TestRegisterExitHandler(t *testing.T) {
	current := len(handlers)

	var results []string

	RegisterExitHandler(func() {
		results = append(results, "test1")
	})

	RegisterExitHandler(func() {
		results = append(results, "test2")
	})

	if len(handlers) != current+2 {
		t.Fatalf("出错，得到长度为 %d ，正确长度应该为 %d \n", len(handlers), current+2)
	}

	runHandlers()

	if len(results) != 2 {
		t.Fatalf("出错，结果数组长度应为 2，实际长度为 %d", len(results))
	}

	if results[0] != "test1" {
		t.Fatalf("出错，结果应该为 test1，实际结果为:%s", results[0])
	}

	if results[1] != "test2" {
		t.Fatalf("出错，结果应该为 test2，实际结果为:%s", results[0])
	}

	t.Log("测试通过！")
}

func TestDeferExitHandler(t *testing.T) {
	current := len(handlers)

	var results []string

	DeferExitHandler(func() {
		results = append(results, "test 111")
	})

	DeferExitHandler(func() {
		results = append(results, "test 222")
	})

	if len(handlers) != current+2 {
		t.Fatalf("出错，期待值 %d，当前值 %d", current+2, len(handlers))
	}

	runHandlers()

	if results[0] != "test 222" {
		t.Fatalf("出错，期待值 %s，当前值 %s", "test 222", results[0])
	}

	if results[1] != "test 111" {
		t.Fatalf("出错，期待值 %s，当前值 %s", "test 111", results[1])
	}

	t.Log("测试通过！")
}

// os.Exit 没法写...
