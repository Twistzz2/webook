package domain

// User 领域对象，是 DDD 中的 entity
// 也叫 BO(Business Object)
type User struct {
	Id       int64 // 主键
	Email    string
	Password string
}

// type Address struct {
// }
