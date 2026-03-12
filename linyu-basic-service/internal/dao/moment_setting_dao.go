package dao

var MomentSetDao = newMomentSetDao()

func newMomentSetDao() *momentSetDao {
	return &momentSetDao{}
}

type momentSetDao struct{}
