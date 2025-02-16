package model

type Inventory struct {
	Items map[string]int // key: item; val: amount
}

func NewInventory() Inventory {
	return Inventory{
		Items: make(map[string]int, 16),
	}
}

func (i *Inventory) AddItem(item string) error {
	if _, found := i.Items[item]; found {
		i.Items[item]++
		return nil
	}

	i.Items[item] = 1

	return nil
}

func (i *Inventory) GetItems() map[string]int {
	return i.Items
}
