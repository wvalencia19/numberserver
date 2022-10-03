package internal

type numberSet map[int]struct{}

func (s numberSet) add(number int) {
	s[number] = struct{}{}
}

func (s numberSet) has(number int) bool {
	_, ok := s[number]
	return ok
}
