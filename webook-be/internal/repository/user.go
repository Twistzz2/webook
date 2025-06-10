package repository

type UserRepository struct {
}

func (r *UserRepository) FindByID(int64) {
	// 先从 cache 找
	// 再从 dao 找
	// 找到了回写 cache
}
