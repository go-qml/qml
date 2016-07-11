package qml

type ItemModelDefaultImpl struct {
}

func (l *ItemModelDefaultImpl) ColumnCount(parent ModelIndex) int {
	return 1
}

func (l *ItemModelDefaultImpl) RowCount(parent ModelIndex) int {
	return 1
}

func (l *ItemModelDefaultImpl) Data(index ModelIndex, role Role) interface{} {
	return nil
}

func (l *ItemModelDefaultImpl) Index(row int, column int, parent ModelIndex) ModelIndex {
	// if !parent.IsValid() && column == 0 && row >= 0 && row < len(l.lines) {
	// 	return l.internal.CreateIndex(row, column, 0)
	// }
	return nil
}

func (l *ItemModelDefaultImpl) Parent(index ModelIndex) ModelIndex {
	return nil
}

// Editing functions

func (l *ItemModelDefaultImpl) Flags(index ModelIndex) ItemFlags {
	return NoItemFlags
}

func (l *ItemModelDefaultImpl) SetData(index ModelIndex, value interface{}, role Role) bool {
	return false
}
