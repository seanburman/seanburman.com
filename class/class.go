package class

import "strings"

type Class string

func (c Class) String() string {
	return string(c)
}

func (c Class) Add(class ...Class) Class {
	for _, v := range class {
		c2 := Class(c.String() + " " + v.String())
		c = c2
	}
	return c
}

func (c Class) Remove(class Class) Class {
	return Class(strings.Replace(c.String(), class.String(), "", -1))
}
