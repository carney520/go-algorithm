package set

import (
	"fmt"
	"testing"
)

func TestSetInsert(t *testing.T) {
	s := New(DefaultMatch)
	s.Insert(1)
	if s.Len() != 1 {
		t.Error("集合插入异常, 长度错误")
	}
	err := s.Insert(1)
	if err != ErrConflict || s.Len() != 1 {
		t.Error("集合插入异常, 成员已存在, 不应该插入")
	}

	err = s.Insert(2)
	if err != nil || s.Len() != 2 || !s.Has(2) {
		t.Error("集合插入失败")
	}
}

func TestSetEqual(t *testing.T) {
	s := New(DefaultMatch, 1, 2, 3, 4)
	cls := s.Clone()
	if !s.Equal(cls) {
		t.Error("集合不等于克隆集合")
	}

	ns := New(DefaultMatch, 3, 1, 4, 2)
	if !s.Equal(ns) {
		t.Error("相同成员集合, 但无需的两个集合应该相等")
	}

	ans := New(DefaultMatch, 1, 2)
	if s.Equal(ans) {
		t.Error("集合应该不相等")
	}
}

func TestSetRemove(t *testing.T) {
	s := New(DefaultMatch, 1, 2, 3, 4)
	s.Remove(2)
	if s.Len() != 3 || s.Has(2) {
		t.Error("集合移除失败, 2没有被移除")
	}

	err := s.Remove(5)
	if err != ErrNotFound {
		t.Error("集合移除异常, 5是不存在的")
	}

	s.Remove(4)
	if s.Len() != 2 || s.Has(4) {
		t.Error("集合移除失败, 4没有被移除")
	}
}

func TestSetUnion(t *testing.T) {
	s := New(DefaultMatch, 1, 2, 3, 4)
	v := New(DefaultMatch, 2, 3, 9)
	u := s.Union(v)
	exp := New(DefaultMatch, 1, 2, 3, 4, 9)
	if !u.Equal(exp) {
		t.Error("集合并集异常")
	}

	empty := New(DefaultMatch)
	if !s.Union(empty).Equal(s) {
		t.Error("集合并集异常, 和空集并集等于自身")
	}
}

func TestSetIntersection(t *testing.T) {
	s := New(DefaultMatch, 1, 2, 3, 4)
	v := New(DefaultMatch, 2, 1, 9, 8)
	i := s.Intersection(v)
	if !i.Equal(New(DefaultMatch, 2, 1)) {
		t.Error("集合交集异常")
	}

	empty := New(DefaultMatch)
	if !s.Intersection(empty).Equal(empty) {
		t.Error("集合并集异常, 和空集交集等于空集")
	}
}

func TestSetDiff(t *testing.T) {
	s := New(DefaultMatch, 1, 3, 4, 5)
	v := New(DefaultMatch, 2, 6, 3)
	d := s.Diff(v)
	if !d.Equal(New(DefaultMatch, 1, 4, 5)) {
		t.Error("集合差集异常")
	}

	empty := New(DefaultMatch)
	if !s.Diff(empty).Equal(s) {
		t.Error("集合差集异常, 集合-空集等于集合本身")
	}

	if !empty.Diff(s).Equal(empty) {
		t.Error("集合差集异常, 空集-集合等于空集")
	}
}

func TestSetSubset(t *testing.T) {
	s := New(DefaultMatch, 1, 2, 3)
	if !s.Subset(New(DefaultMatch, 1, 2)) {
		t.Error("集合子集异常, {1, 2} 是 {1, 2, 3}的子集")
	}

	if s.Subset(New(DefaultMatch, 1, 2, 3, 4)) {
		t.Error("集合子集异常, {1, 2, 3, 4} 不是 {1, 2, 3}的子集")
	}

	if !s.Subset(s) {
		t.Error("集合子集异常, 集合自身也是子集")
	}

	if !s.Subset(New(DefaultMatch)) {
		t.Error("集合子集异常, 空集是任何集合的子集")
	}
}

// cover 实现集合覆盖
// members是所有技能集合
// subsets是members的子集A1， 到An组成的集合。 cover 从subsets中找出覆盖members最高的项
func cover(members *Set, subsets *Set) (*Set, error) {
	members = members.Clone()
	subsets = subsets.Clone()
	rt := New(DefaultMatch)
	var maxLen int
	var maxMember *Set
	for members.Len() > 0 && subsets.Len() > 0 {
		maxLen = 0

		// 遍历subset
		subsets.Each(func(data interface{}, _ int) bool {
			member, ok := data.(*Set)
			if !ok {
				panic("Type error")
			}
			i := member.Intersection(members)
			if i.Len() > maxLen {
				maxMember = member
				maxLen = i.Len()
			}
			return false
		})

		// 没有任何子集覆盖
		if maxLen == 0 {
			return nil, ErrNotFound
		}
		rt.Insert(maxMember)
		members = members.Diff(maxMember)
		subsets.Remove(maxMember)
	}
	if rt.Len() == 0 {
		return nil, ErrNotFound
	}
	return rt, nil
}

func Example() {
	skills := New(DefaultMatch, "c++", "go", "python", "ruby", "java")
	group := New(func(a, b interface{}) bool {
		sa := a.(*Set)
		sb := b.(*Set)
		return sa.Equal(sb)
	},
		New(DefaultMatch, "python", "go", "ruby"),
		New(DefaultMatch, "c++", "python"),
		New(DefaultMatch, "java", "ruby"),
		New(DefaultMatch, "go", "java", "ruby"),
	)
	best, err := cover(skills, group)
	if err != nil {
		fmt.Println("覆盖失败")
	} else {
		fmt.Println(best)
	}
}
