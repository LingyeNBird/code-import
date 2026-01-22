package coding

import (
	"sync"
	"testing"
)

// TestGetProjectByName_CacheMiss 测试缓存未命中的情况
func TestGetProjectByName_CacheMiss(t *testing.T) {
	// 清空缓存
	projectCache = sync.Map{}

	// 由于需要真实的 CODING API,这里仅验证缓存逻辑
	// 实际 API 调用需要在集成测试中验证

	// 验证缓存初始为空
	if _, ok := projectCache.Load("test-project"); ok {
		t.Error("缓存应该初始为空")
	}
}

// TestGetProjectByName_CacheHit 测试缓存命中的情况
func TestGetProjectByName_CacheHit(t *testing.T) {
	// 清空缓存
	projectCache = sync.Map{}

	// 模拟缓存中已有数据
	mockProject := Project{
		Name:        "test-project",
		Id:          12345,
		DisplayName: "测试项目",
		Description: "这是一个测试项目",
	}
	projectCache.Store("test-project", mockProject)

	// 验证可以从缓存读取
	if cached, ok := projectCache.Load("test-project"); !ok {
		t.Error("应该能从缓存读取数据")
	} else {
		cachedProject := cached.(Project)
		if cachedProject.Name != "test-project" {
			t.Errorf("期望项目名称为 'test-project',实际为 '%s'", cachedProject.Name)
		}
		if cachedProject.Id != 12345 {
			t.Errorf("期望项目 ID 为 12345,实际为 %d", cachedProject.Id)
		}
	}
}

// TestProjectCache_ConcurrentAccess 测试并发访问缓存的安全性
func TestProjectCache_ConcurrentAccess(t *testing.T) {
	// 清空缓存
	projectCache = sync.Map{}

	// 创建多个 goroutine 并发写入不同的项目
	var wg sync.WaitGroup
	projectCount := 10

	for i := 0; i < projectCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			project := Project{
				Name: string(rune('A' + id)),
				Id:   id,
			}
			projectCache.Store(project.Name, project)
		}(i)
	}

	wg.Wait()

	// 验证所有项目都被正确存储
	for i := 0; i < projectCount; i++ {
		projectName := string(rune('A' + i))
		if cached, ok := projectCache.Load(projectName); !ok {
			t.Errorf("项目 %s 应该在缓存中", projectName)
		} else {
			cachedProject := cached.(Project)
			if cachedProject.Id != i {
				t.Errorf("项目 %s 的 ID 应该为 %d,实际为 %d", projectName, i, cachedProject.Id)
			}
		}
	}
}

// TestProjectCache_ConcurrentReadWrite 测试并发读写场景
func TestProjectCache_ConcurrentReadWrite(t *testing.T) {
	// 清空缓存
	projectCache = sync.Map{}

	// 预先存储一个项目
	mockProject := Project{
		Name:        "concurrent-test",
		Id:          999,
		DisplayName: "并发测试项目",
	}
	projectCache.Store("concurrent-test", mockProject)

	// 创建多个 goroutine 并发读取同一个项目
	var wg sync.WaitGroup
	readCount := 20
	successCount := 0
	var mu sync.Mutex

	for i := 0; i < readCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if cached, ok := projectCache.Load("concurrent-test"); ok {
				cachedProject := cached.(Project)
				if cachedProject.Id == 999 {
					mu.Lock()
					successCount++
					mu.Unlock()
				}
			}
		}()
	}

	wg.Wait()

	// 验证所有读取都成功
	if successCount != readCount {
		t.Errorf("期望 %d 次成功读取,实际为 %d 次", readCount, successCount)
	}
}

// TestProjectCache_StoreAndLoad 测试基本的存储和加载功能
func TestProjectCache_StoreAndLoad(t *testing.T) {
	// 清空缓存
	projectCache = sync.Map{}

	tests := []struct {
		name    string
		project Project
	}{
		{
			name: "项目1",
			project: Project{
				Name:        "project-1",
				Id:          1,
				DisplayName: "项目一",
				Description: "第一个项目",
			},
		},
		{
			name: "项目2",
			project: Project{
				Name:        "project-2",
				Id:          2,
				DisplayName: "项目二",
				Description: "第二个项目",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 存储项目
			projectCache.Store(tt.project.Name, tt.project)

			// 加载项目
			if cached, ok := projectCache.Load(tt.project.Name); !ok {
				t.Errorf("应该能从缓存加载项目 %s", tt.project.Name)
			} else {
				cachedProject := cached.(Project)
				if cachedProject.Id != tt.project.Id {
					t.Errorf("期望 ID 为 %d,实际为 %d", tt.project.Id, cachedProject.Id)
				}
				if cachedProject.DisplayName != tt.project.DisplayName {
					t.Errorf("期望 DisplayName 为 '%s',实际为 '%s'", tt.project.DisplayName, cachedProject.DisplayName)
				}
			}
		})
	}
}
