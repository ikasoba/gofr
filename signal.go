package gofr

import "syscall/js"

var deferedBatches = []func(){}
var isBatchProcessing = false

// Batch の実行後に処理を実行する
func DeferBatch(fn func()) {
	deferedBatches = append(deferedBatches, fn)
}

// signal の更新を一括で行うもの
func Batch(fn func()) {
	if isBatchProcessing {
		fn()
		return
	}

	isBatchProcessing = true

	fn()

	isBatchProcessing = false

	ProcessDeferedBatch()
}

func ProcessDeferedBatch() {
	for len(deferedBatches) > 0 {
		fn := deferedBatches[0]
		deferedBatches = deferedBatches[1:]

		fn()
	}
}

// signalの更新イベントを重複させないためのもの

var broadcastSignals = map[Transmitter]func(){}
var isBroadcasting = false

func deliver(s Transmitter, fn func()) {
	if _, ok := broadcastSignals[s]; ok {
		return
	}

	broadcastSignals[s] = fn

	if isBroadcasting {
		return
	}

	isBroadcasting = true

	DeferBatch(func() {
		m := broadcastSignals
		broadcastSignals = make(map[Transmitter]func())

		for _, fn := range m {
			fn()
		}

		isBroadcasting = false
	})

	if !isBatchProcessing {
		Batch(func() {})
	}
}

type Transmitter interface {
	Subscribe(func())
}

var referencedSignals = map[Transmitter]any{}

func captureReferencedSignals() func() []Transmitter {
	prev := referencedSignals
	signals := map[Transmitter]any{}
	referencedSignals = signals

	return func() []Transmitter {
		referencedSignals = prev

		var res = []Transmitter{}

		for s := range signals {
			res = append(res, s)
		}

		return res
	}
}

// 値の変更を受信、送信するもの

type Signal[T comparable] struct {
	val         T
	subscribers []func()
}

func NewSignal[T comparable](val T) *Signal[T] {
	return &Signal[T]{val, nil}
}

func (s *Signal[T]) GetAsAny() any {
	return any(s.val)
}

func (s *Signal[T]) Get() T {
	referencedSignals[s] = nil

	return s.val
}

// 依存関係として追跡させずに signal から値を取り出す
func (s *Signal[T]) GetSilent() T {
	return s.val
}

func (s *Signal[T]) Set(val T) {
	s.val = val

	deliver(s, s.broadcast)
}

func (s *Signal[T]) broadcast() {
	for _, fn := range s.subscribers {
		fn()
	}
}

func (s *Signal[T]) Subscribe(fn func()) {
	s.subscribers = append(s.subscribers, fn)
}

func (s *Signal[T]) ToNode() js.Value {
	node := newNode(s.val)

	s.Subscribe(func() {
		node = replaceNode(s.val, node)
	})

	return node
}

// 複数の signal に依存した値を持つ signal を生成する
func Computed[T comparable](fn func() T) *Signal[T] {
	escapeLifeTime := BeginLifeTime()
	collectDeps := captureReferencedSignals()

	val := fn()

	deps := collectDeps()
	lifeTime := escapeLifeTime()

	s := NewSignal(val)

	isProcessing := false

	{
		l := GetCurrentLifeTime()
		l.MountHandlers = append(l.MountHandlers, lifeTime.MountHandlers...)
	}

	var handler = func() {
		if isProcessing {
			return
		}

		isProcessing = true

		DeferBatch(func() {
			lifeTime.Destroy()

			escapeLifeTime = BeginLifeTime()
			val = fn()
			lifeTime = escapeLifeTime()

			s.Set(val)

			lifeTime.DispatchMount()

			isProcessing = false
		})
	}

	for _, d := range deps {
		d.Subscribe(handler)
	}

	return s
}
