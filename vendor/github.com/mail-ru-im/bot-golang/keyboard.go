package botgolang

import (
	"fmt"
)

// Keyboard represents an inline keyboard markup
// Call the NewKeyboard() func to get a keyboard instance
type Keyboard struct {
	Rows [][]Button
}

// NewKeyboard returns a new keyboard instance
func NewKeyboard() Keyboard {
	return Keyboard{
		Rows: make([][]Button, 0),
	}
}

// AddRows adds a row to the keyboard
func (k *Keyboard) AddRow(row ...Button) {
	k.Rows = append(k.Rows, row)
}

// AddButton adds a button to the end of the row
func (k *Keyboard) AddButton(rowIndex int, button Button) error {
	if ok := k.checkRow(rowIndex); !ok {
		return fmt.Errorf("no such row: %d", rowIndex)
	}

	k.Rows[rowIndex] = append(k.Rows[rowIndex], button)
	return nil
}

// DeleteRow removes the row from the keyboard
func (k *Keyboard) DeleteRow(index int) error {
	if ok := k.checkRow(index); !ok {
		return fmt.Errorf("no such row: %d", index)
	}

	k.Rows = append(k.Rows[:index], k.Rows[index+1:]...)
	return nil
}

// DeleteButton removes the button from the row.
// Note - at least one button should remain in a row,
// if you want to delete all buttons, use the DeleteRow function
func (k *Keyboard) DeleteButton(rowIndex, buttonIndex int) error {
	if ok := k.checkButton(rowIndex, buttonIndex); !ok {
		return fmt.Errorf("no button at index %d or row %d", buttonIndex, rowIndex)
	}

	if k.RowSize(rowIndex) < 2 {
		return fmt.Errorf("can't delete button: at least one should remain in a row")
	}

	row := &k.Rows[rowIndex]
	*row = append((*row)[:buttonIndex], (*row)[buttonIndex+1:]...)
	return nil
}

// ChangeButton changes the button to a new one at the specified position
func (k *Keyboard) ChangeButton(rowIndex, buttonIndex int, newButton Button) error {
	if ok := k.checkButton(rowIndex, buttonIndex); !ok {
		return fmt.Errorf("no button at index %d or row %d", buttonIndex, rowIndex)
	}

	k.Rows[rowIndex][buttonIndex] = newButton
	return nil
}

// SwapRows swaps two rows in keyboard
func (k *Keyboard) SwapRows(first, second int) error {
	if ok := k.checkRow(first); !ok {
		return fmt.Errorf("no such index (first): %d", first)
	}
	if ok := k.checkRow(second); !ok {
		return fmt.Errorf("no such index (second): %d", second)
	}

	k.Rows[first], k.Rows[second] = k.Rows[second], k.Rows[first]
	return nil
}

// RowsCount returns the number of rows
func (k *Keyboard) RowsCount() int {
	return len(k.Rows)
}

// RowSize returns the number of buttons in a row.
// If there is no such row, then returns -1
func (k *Keyboard) RowSize(row int) int {
	if ok := k.checkRow(row); !ok {
		return -1
	}
	return len(k.Rows[row])
}

// GetKeyboard returns an array of button rows
func (k *Keyboard) GetKeyboard() [][]Button {
	return k.Rows
}

// checkRow checks that the index of row doesnt go beyond the bounds of the array
func (k *Keyboard) checkRow(i int) bool {
	return i >= 0 && i < len(k.Rows)
}

// checkButton checks that the button and row indexes doesnt go beyond the bounds of the array
func (k *Keyboard) checkButton(row, button int) bool {
	return k.checkRow(row) && button >= 0 && button < len(k.Rows[row])
}
