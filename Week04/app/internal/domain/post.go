package domain

//DO  可以用orm解决,这里使用gorm
type Post struct {
	Id    int64
	Title string
}
