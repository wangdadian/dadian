// 单向列表： One-Way List
// 说明：插入第一个元素后，后续的插入必须与第一个元素类型相同，否则插入失败
// 插入：支持从头部或者尾部插入元素
// 弹出：支持从头部或者尾部弹出元素
// 遍历：仅支持单向从头部开始遍历，使用迭代器Iterator遍历

package owlist

import (
	"errors"
	"reflect"
)

var (
	EOL = errors.New("end of list")
	TNC = errors.New("type is not consistent with the first element")
)

type List struct {
	size  uint64    // 列表长度
	first *Iterator // 列表头
	tail  *Iterator // 列表尾
}

func NewList() *List {
	pl := &List{
		size:  0,
		first: nil,
		tail:  nil,
	}
	return pl
}

// 判断插入的数据是否和之前的数据类型一致
func (self *List) isSameType(v interface{}) bool {
	if self.size == 0 {
		return true
	}
	szTypeMe := reflect.TypeOf(v).String()
	szTypeFirst := reflect.TypeOf(self.first.Value()).String()
	if szTypeMe == szTypeFirst {
		return true
	}
	return false
}

// 在列表首部添加元素
func (self *List) PushFront(v interface{}) error {
	if ok := self.isSameType(v); !ok {
		return TNC
	}
	pNewIt := &Iterator{
		v:    v,
		prev: nil,
		next: self.first,
	}
	// 如果这是插入的第一个元素，尾部也指向这个元素
	if self.size == 0 {
		self.tail = pNewIt
	} else {
		self.first.prev = pNewIt
	}
	self.first = pNewIt

	// 列表长度增1
	self.size++
	return nil
}

// 在列表尾部添加元素
func (self *List) PushBack(v interface{}) error {
	if ok := self.isSameType(v); !ok {
		return TNC
	}
	pNewIt := &Iterator{
		v:    v,
		prev: self.tail,
		next: nil,
	}
	if self.size == 0 {
		self.first = pNewIt
	} else {
		self.tail.next = pNewIt
	}
	// 列表尾部元素指向最新元素
	self.tail = pNewIt
	// 如果之前列表为空，则列表首部元素也指向此新增元素

	// 列表长度增1
	self.size++
	return nil
}

// 获取列表长度
func (self *List) Size() uint64 {
	return self.size
}

// 是否为空
func (self *List) IsEmpty() bool {
	if self.size == 0 {
		return true
	}
	return false
}

// 清空元素
func (self *List) Clear() {
	self.size = 0
	self.first = nil
	self.tail = nil
}

// 取出首部位置元素或者nil，并从列表删除
func (self *List) PopFront() interface{} {
	if self.size == 0 {
		return nil
	}
	v := self.first.Value()

	self.first = self.first.next

	self.size--
	if self.size == 0 {
		self.tail = nil
	} else {
		self.first.prev = nil
	}

	return v
}

// 弹出尾部位置元素或者nil，并从列表删除
func (self *List) PopBack() interface{} {
	if self.size == 0 {
		return nil
	}
	v := self.tail.Value()
	self.tail = self.tail.prev
	self.size--

	if self.size == 0 {
		self.first = nil
	} else {
		self.tail.next = nil
	}
	return v
}

// 清除迭代器指向的元素，返回下一个元素或者nil
func (self *List) Erase(it *Iterator) *Iterator {
	prevIt := it.prevElement()
	nextIt := it.nextElement()
	if prevIt != nil {
		// 删除的非首元素
		prevIt.next = nextIt
	} else {
		// 删除的首元素
		self.first = nextIt
	}
	if nextIt != nil {
		// 删除的非尾部元素
		nextIt.prev = prevIt
	} else {
		// 删除的尾部元素
		self.tail = prevIt
	}
	return nextIt
}

// 头部元素
func (self *List) Begin() *Iterator {
	return self.first
}

// 尾部元素的下一个nil元素
func (self *List) End() *Iterator {
	if self.tail == nil {
		return nil
	}
	return self.tail.next
}
