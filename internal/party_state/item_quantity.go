package party_state

type ItemQuantityLarge struct {
	quantity uint16
}

func (l *ItemQuantityLarge) HasSome() bool {
	return l.Get() > 0
}

type ItemQuantitySmall struct {
	quantity uint16
}

func (s *ItemQuantitySmall) HasSome() bool {
	return s.Get() > 0
}

func (l *ItemQuantityLarge) Get() uint16 {
	return l.quantity
}

func (l *ItemQuantityLarge) Set(quantity uint16) {
	l.quantity = min(quantity, itemQuantityLargeMax)
}

func (s *ItemQuantitySmall) Get() uint16 {
	return s.quantity
}

func (s *ItemQuantitySmall) Set(quantity uint16) {
	s.quantity = min(quantity, itemQuantitySmallMax)
}

const itemQuantityLargeMax uint16 = 9999
const itemQuantitySmallMax uint16 = 99

type ItemQuantity interface {
	IncrementByOne()
	DecrementByOne() bool
	IncrementBy(incBy uint16)
	DecrementBy(incBy uint16) bool
	Set(quantity uint16)
	Get() uint16
	HasSome() bool
}

func incBy(currentQuantity uint16, incBy uint16, maxQuantity uint16) uint16 {
	return min(currentQuantity+incBy, maxQuantity)
}

func decBy(currentQuantity uint16, decBy uint16) (uint16, bool) {
	if decBy > currentQuantity {
		return 0, false
	}
	return currentQuantity - decBy, true
}

func (l *ItemQuantityLarge) IncrementByOne() {
	l.quantity = incBy(l.quantity, 1, itemQuantityLargeMax)
}

func (l *ItemQuantityLarge) DecrementByOne() bool {
	newQuantity, bWasEnough := decBy(l.quantity, 1)
	if !bWasEnough {
		return false
	}
	l.quantity = newQuantity
	return true
}

func (l *ItemQuantityLarge) IncrementBy(incByQuantity uint16) {
	l.quantity = incBy(l.quantity, incByQuantity, itemQuantityLargeMax)
}

func (l *ItemQuantityLarge) DecrementBy(decByQuantity uint16) bool {
	newQuantity, bWasEnough := decBy(l.quantity, decByQuantity)
	if !bWasEnough {
		return false
	}
	l.quantity = newQuantity
	return true
}

func (s *ItemQuantitySmall) IncrementByOne() {
	s.quantity = incBy(s.quantity, 1, itemQuantitySmallMax)
}

func (s *ItemQuantitySmall) DecrementByOne() bool {
	newQuantity, bWasEnough := decBy(s.quantity, 1)
	if !bWasEnough {
		return false
	}
	s.quantity = newQuantity
	return true
}

func (s *ItemQuantitySmall) IncrementBy(incByQuantity uint16) {
	s.quantity = incBy(s.quantity, incByQuantity, itemQuantitySmallMax)
}

func (s *ItemQuantitySmall) DecrementBy(decByQuantity uint16) bool {
	newQuantity, bWasEnough := decBy(s.quantity, decByQuantity)
	if !bWasEnough {
		return false
	}
	s.quantity = newQuantity
	return true
}
