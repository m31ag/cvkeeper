package model

func GetChoicesByLevel(level int) ([]Item, int) {
	menu0 := []string{"token", "passwords"}
	var items = make([]Item, 0)
	var actual []string
	if level == 0 {
		actual = menu0
	}

	for i := 0; i < len(actual); i++ {
		items = append(items, Item{
			Name:      actual[i],
			Operation: nil,
		})
	}

	return items, level
}
