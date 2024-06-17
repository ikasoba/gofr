package gofr

type LifeTime struct {
	CleanupHandlers []func()
	MountHandlers   []func()
}

func (lifeTime *LifeTime) Destroy() {
	for _, fn := range lifeTime.CleanupHandlers {
		fn()
	}
}

func (lifeTime *LifeTime) DispatchMount() {
	for _, fn := range lifeTime.MountHandlers {
		fn()
	}
}

var ctx = &LifeTime{}

func GetCurrentLifeTime() *LifeTime {
	return ctx
}

func BeginLifeTime() func() *LifeTime {
	prev := ctx
	ctx = &LifeTime{}

	return func() *LifeTime {
		res := ctx
		ctx = prev

		return res
	}
}

// コンポーネントのクリーンアップ時に処理を実行する
func OnCleanup(fn func()) {
	ctx.CleanupHandlers = append(ctx.CleanupHandlers, fn)
}

// コンポーネントがページにマウントされたときに処理を実行する
func OnMount(fn func()) {
	ctx.MountHandlers = append(ctx.MountHandlers, fn)
}
