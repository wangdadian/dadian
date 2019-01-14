package owlist

import (
	"fmt"
	"testing"
	"time"
)

type T1 struct {
	n int
	s string
}

type T2 struct {
	n int
	s string
}

func (self *T1) String() string {
	return fmt.Sprintf("n: %8d, s: %s", self.n, self.s)
}
func TestPushBack(t *testing.T) {
	pList := NewList()
	bys := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	iLen := len(bys)
	iEnd := 0
	MAX_LOOP := iLen * 2
	for i := 0; i < MAX_LOOP; i++ {
		if iEnd >= iLen {
			iEnd = iEnd % iLen
		}
		iEnd += 1
		t1 := &T1{n: i + 1, s: string(bys[:iEnd])}
		//t.Logf("new T1: %s\n", t1.String())
		pList.PushBack(t1)
	}
	iListLen := pList.Size()
	t.Logf("list size: %d\n", iListLen)
	if yes := pList.IsEmpty(); yes {
		t.Errorf("list is empty")
	}
	// 不一致的类型插入
	var iElm int = 100
	err := pList.PushBack(iElm)
	if err != nil {
		t.Logf("value[%v]: push back failed: %s\n", iElm, err.Error())
	}
	var iElmT2 T2 = T2{n: 9999, s: "hello,world"}
	err = pList.PushBack(iElmT2)
	if err != nil {
		t.Logf("value[%v]: push back failed: %s\n", iElmT2, err.Error())
	}

	for it := pList.Begin(); it != pList.End(); it = it.Next() {
		v := it.Value()
		if t1, ok := v.(*T1); ok {
			t.Logf("T1: %s\n", t1.String())
		}
	}
}
func TestPushFront(t *testing.T) {
	pList := NewList()
	bys := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	iLen := len(bys)
	iEnd := 0
	MAX_LOOP := iLen * 2
	for i := 0; i < MAX_LOOP; i++ {
		if iEnd >= iLen {
			iEnd = iEnd % iLen
		}
		iEnd += 1
		t1 := &T1{n: i + 1, s: string(bys[:iEnd])}
		//t.Logf("new T1: %s\n", t1.String())
		pList.PushFront(t1)
	}
	iListLen := pList.Size()
	t.Logf("list size: %d\n", iListLen)
	if yes := pList.IsEmpty(); yes {
		t.Errorf("list is empty")
	}
	for it := pList.Begin(); it != pList.End(); it = it.Next() {
		v := it.Value()
		if t1, ok := v.(*T1); ok {
			t.Logf("T1: %s\n", t1.String())
		}
	}
}

func TestPopFront(t *testing.T) {
	pList := NewList()
	bys := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	iLen := len(bys)
	iEnd := 0
	MAX_LOOP := iLen * 2
	for i := 0; i < MAX_LOOP; i++ {
		if iEnd >= iLen {
			iEnd = iEnd % iLen
		}
		iEnd += 1
		t1 := &T1{n: i + 1, s: string(bys[:iEnd])}
		//t.Logf("new T1: %s\n", t1.String())
		pList.PushBack(t1)
	}
	iListLen := pList.Size()
	t.Logf("list size: %d\n", iListLen)
	if yes := pList.IsEmpty(); yes {
		t.Errorf("list is empty")
	}

	for {
		t.Logf("list size: %d\n", pList.Size())
		if pList.IsEmpty() {
			break
		}

		v := pList.PopFront()
		if v != nil {
			if t1, ok := v.(*T1); ok {
				t.Logf("T1: %s\n", t1.String())
			}

		} else {
			continue
		}
	}
}

func TestPopBack(t *testing.T) {
	pList := NewList()
	bys := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	iLen := len(bys)
	iEnd := 0
	MAX_LOOP := iLen * 2
	for i := 0; i < MAX_LOOP; i++ {
		if iEnd >= iLen {
			iEnd = iEnd % iLen
		}
		iEnd += 1
		t1 := &T1{n: i + 1, s: string(bys[:iEnd])}
		//t.Logf("new T1: %s\n", t1.String())
		pList.PushBack(t1)
	}
	iListLen := pList.Size()
	t.Logf("list size: %d\n", iListLen)
	if yes := pList.IsEmpty(); yes {
		t.Errorf("list is empty")
	}

	for {
		t.Logf("list size: %d\n", pList.Size())
		if pList.IsEmpty() {
			break
		}

		v := pList.PopBack()
		if v != nil {
			if t1, ok := v.(*T1); ok {
				t.Logf("T1: %s\n", t1.String())
			}

		} else {
			continue
		}
	}
}

func TestErase(t *testing.T) {
	pList := NewList()
	bys := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	iLen := len(bys)
	iEnd := 0
	MAX_LOOP := 12
	for i := 0; i < MAX_LOOP; i++ {
		if iEnd >= iLen {
			iEnd = iEnd % iLen
		}
		iEnd += 1
		t1 := &T1{n: i + 1, s: string(bys[:iEnd])}
		//t.Logf("new T1: %s\n", t1.String())
		pList.PushBack(t1)
	}
	t.Logf("before list erase:\n")
	for it := pList.Begin(); it != pList.End(); it = it.Next() {
		v := it.Value()
		if t1, ok := v.(*T1); ok {
			t.Logf("T1: %s\n", t1.String())
		}
	}
	// erase
	var iIndex int = 0
	for it := pList.Begin(); it != pList.End(); it = it.Next() {
		v := it.Value()
		if t1, ok := v.(*T1); ok {
			if t1.n%2 != 0 {
				pList.Erase(it)
			}
		}
		iIndex++
	}

	t.Logf("after list erase:\n")
	for it := pList.Begin(); it != pList.End(); it = it.Next() {
		v := it.Value()
		if t1, ok := v.(*T1); ok {
			t.Logf("T1: %s\n", t1.String())
		}
	}
	tnew := &T1{n: 10000, s: "rrrrrrrrr"}
	pList.PushFront(tnew)
	tnew1 := &T1{n: 87878, s: "QQQQQQQQQ"}
	pList.PushBack(tnew1)
	t.Logf("after list pushfront:\n")
	t.Logf("first: %v\n", &pList.first)
	t.Logf("tail: %v\n", &pList.tail)
	for it := pList.Begin(); it != pList.End(); it = it.Next() {
		t.Logf("%v\n", it)
		// v := it.Value()
		// if t1, ok := v.(*T1); ok {
		// 	t.Logf("T1: %s\n", t1.String())
		// }
	}
}

func TestPushMore(t *testing.T) {
	pList := NewList()

	tStart := time.Now().UnixNano()
	for i := 0; i < 1000000; i++ {
		t1 := &T1{n: i + 1, s: string(make([]byte, 2048))}
		pList.PushBack(t1)
	}
	tNow := time.Now().UnixNano()
	t.Logf("push spend time: %.03f ms", float64(tNow-tStart)/1000.0/1000.0)
}
