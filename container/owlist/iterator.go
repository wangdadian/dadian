package owlist

type Iterator struct {
	k    uint64      // 键值
	v    interface{} // 值
	prev *Iterator   // 前一个元素
	next *Iterator   // 下一个元素
}

// 返回迭代器指示的元素的值
func (self *Iterator) Value() interface{} {
	return self.v
}

// 迭代器向后移动一个
func (self *Iterator) Next() *Iterator {
	return self.nextElement()
}

// 迭代器向后移动一个
func (self *Iterator) nextElement() *Iterator {
	self = self.next
	return self
}

// 迭代器向前移动一个
func (self *Iterator) prevElement() *Iterator {
	self = self.prev
	return self
}
