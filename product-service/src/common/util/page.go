package util

func CreateLimitAndOffset(page int) (limit, offset int) {
	const limitDefault = 50
	
	if page == 0 {
		return limitDefault, 0
	}

	return limitDefault, (page - 1) * limitDefault
}