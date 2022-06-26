package main

import "fmt"

type Boss struct {
	name     string
	delegate map[string]func()
}
type Staff struct {
	name string
}

//消息发送者
type Notificator interface {
	Attach(name string, handleFunc func())
	Detach(name string) error
	Notify() error
}

//观察者
type Observer interface {
	Handle()
}

func (b *Boss) Attach(name string, handleFunc func()) {
	if b.delegate == nil {
		b.delegate = make(map[string]func())
	}
	b.delegate[name] = handleFunc

}

func (b *Boss) Detach(name string) error {
	if _, ok := b.delegate[name]; ok {
		delete(b.delegate, name)
		return nil
	} else {
		return fmt.Errorf("handle function <%s> does not exist", name)
	}
}

func (b *Boss) Notify() error {
	if b.delegate != nil {
		for _, handleFunc := range b.delegate {
			handleFunc()
		}
		return nil
	} else {
		return fmt.Errorf("no handle function was attached")
	}
}

func (s *Staff) Handle() {
	fmt.Printf("%s is handling boss coming\n", s.name)
}

func main() {
	s1 := Staff{name: "jack"}
	s2 := Staff{name: "tom"}
	boss := Boss{name: "boss"}
	boss.Attach(s1.name, s1.Handle)
	boss.Attach(s2.name, s2.Handle)
	err := boss.Notify()
	if err != nil {
		fmt.Println(err)
	}

	boss.Detach(s1.name)
	boss.Notify()
	err = boss.Detach(s1.name)
	fmt.Println(err)

}
