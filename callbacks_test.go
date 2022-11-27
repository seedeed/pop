package pop

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UnAfterFindableModel struct {
	ID    int
	After string
}

type AfterFindableModel struct {
	UnAfterFindableModel
}

var (
	modelWith10FriendsForBench     *Model
	modelWith100FriendsForBench    *Model
	modelWith1000FriendsForBench   *Model
	modelWith10000FriendsForBench  *Model
	modelWith100000FriendsForBench *Model
	fakeConn                       = &Connection{}
)

func init() {
	modelWith10FriendsForBench = &Model{Value: newFriends(10)}
	modelWith100FriendsForBench = &Model{Value: newFriends(100)}
	modelWith1000FriendsForBench = &Model{Value: newFriends(1000)}
	modelWith10000FriendsForBench = &Model{Value: newFriends(10000)}
	modelWith100000FriendsForBench = &Model{Value: newFriends(100000)}
}

func newFriends(size int) []Friend {
	return make([]Friend, size)
}

func (m *AfterFindableModel) AfterFind(*Connection) error {
	m.After = makeAfterString(m.ID)
	return nil
}

func makeAfterString(id int) string {
	return fmt.Sprintf("after %d", id)
}

func Test_Callbacks(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		user := &CallbacksUser{
			BeforeS: "BS",
			BeforeC: "BC",
			BeforeU: "BU",
			BeforeD: "BD",
			BeforeV: "BV",
			AfterS:  "AS",
			AfterC:  "AC",
			AfterU:  "AU",
			AfterD:  "AD",
			AfterF:  "AF",
		}

		r.NoError(tx.Save(user))

		r.Equal("BeforeSave", user.BeforeS)
		r.Equal("BeforeCreate", user.BeforeC)
		r.Equal("AfterSave", user.AfterS)
		r.Equal("AfterCreate", user.AfterC)
		r.Equal("BU", user.BeforeU)
		r.Equal("AU", user.AfterU)

		r.NoError(tx.Update(user))

		r.Equal("BeforeUpdate", user.BeforeU)
		r.Equal("AfterUpdate", user.AfterU)
		r.Equal("BD", user.BeforeD)
		r.Equal("AD", user.AfterD)

		r.Equal("AF", user.AfterF)
		r.NoError(tx.Find(user, user.ID))
		r.Equal("AfterFind", user.AfterF)

		r.NoError(tx.Destroy(user))

		r.Equal("BeforeDestroy", user.BeforeD)
		r.Equal("AfterDestroy", user.AfterD)

		verrs, err := tx.ValidateAndSave(user)
		r.False(verrs.HasAny())
		r.NoError(err)
		r.Equal("BeforeValidate", user.BeforeV)
	})
}

func Test_Callbacks_on_Slice(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)
		for i := 0; i < 2; i++ {
			r.NoError(tx.Create(&CallbacksUser{}))
		}

		users := CallbacksUsers{}
		r.NoError(tx.All(&users))

		r.Len(users, 2)

		for _, u := range users {
			r.Equal("AfterFind", u.AfterF)
		}
	})
}

func BenchmarkModelWith10Friends_afterFind(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if e := modelWith10FriendsForBench.afterFind(fakeConn); e != nil {
			b.Fatalf("benchmark fail. %v\n", e)
		}
	}
}

func BenchmarkModelWith100Friends_afterFind(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if e := modelWith100FriendsForBench.afterFind(fakeConn); e != nil {
			b.Fatalf("benchmark fail. %v\n", e)
		}
	}
}

func BenchmarkModelWith1000Friends_afterFind(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if e := modelWith1000FriendsForBench.afterFind(fakeConn); e != nil {
			b.Fatalf("benchmark fail. %v\n", e)
		}
	}
}

func BenchmarkModelWith10000Friends_afterFind(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if e := modelWith10000FriendsForBench.afterFind(fakeConn); e != nil {
			b.Fatalf("benchmark fail. %v\n", e)
		}
	}
}

func BenchmarkModelWith100000Friends_afterFind(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if e := modelWith100000FriendsForBench.afterFind(fakeConn); e != nil {
			b.Fatalf("benchmark fail. %v\n", e)
		}
	}
}

func TestModel_afterFind(t *testing.T) {
	r := require.New(t)
	{
		list := []AfterFindableModel{
			{
				UnAfterFindableModel{ID: 1113},
			},
			{
				UnAfterFindableModel{ID: 1114},
			},
			{
				UnAfterFindableModel{ID: 1115},
			},
			{
				UnAfterFindableModel{ID: 1116},
			},
		}
		model := &Model{Value: list}

		r.NoError(model.afterFind(fakeConn))

		for _, item := range list {
			r.Equal(makeAfterString(item.ID), item.After)
		}
	}

	{
		list := []UnAfterFindableModel{
			{
				ID: 1113,
			},
			{
				ID: 1114,
			},
			{
				ID: 1115,
			},
			{
				ID: 1116,
			},
		}
		model := &Model{Value: list}

		r.NoError(model.afterFind(fakeConn))

		for _, item := range list {
			r.Equal("", item.After)
		}
	}
}
