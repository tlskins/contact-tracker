package repository

import (
	"context"

	"github.com/globalsign/mgo"
	"github.com/google/uuid"

	m "github.com/contact-tracker/apiService/pkg/mongo"
	t "github.com/contact-tracker/apiService/users/types"
)

var ColUsers = "users"

// MongoUserRepository -
type MongoUserRepository struct {
	session *mgo.Session
	dbname  string
}

// NewMongoUserRepository -
func NewMongoUserRepository(m *mgo.Session, dbname string) *MongoUserRepository {
	return &MongoUserRepository{m, dbname}
}

func (r *MongoUserRepository) DB() (*mgo.Session, *mgo.Database) {
	sess := r.session.Copy()
	return sess, sess.DB(r.dbname)
}

func (r *MongoUserRepository) C(colName string) (*mgo.Session, *mgo.Collection) {
	sess := r.session.Copy()
	return sess, sess.DB(r.dbname).C(colName)
}

func (r *MongoUserRepository) Get(_ context.Context, id string) (resp *t.User, err error) {
	sess, c := r.C(ColUsers)
	defer sess.Close()

	err = m.FindOne(c, &resp, m.M{"_id": id})
	return
}

func (r *MongoUserRepository) GetAll(_ context.Context) (resp []*t.User, err error) {
	sess, c := r.C(ColUsers)
	defer sess.Close()

	err = m.Find(c, &resp, m.M{})
	return
}

func (r *MongoUserRepository) Update(_ context.Context, user *t.UpdateUser) (*t.User, error) {
	sess, c := r.C(ColUsers)
	defer sess.Close()

	var resp t.User
	err := m.Update(c, &resp, m.M{"_id": user.ID}, m.M{"$set": user})
	return &resp, err
}

func (r *MongoUserRepository) Create(_ context.Context, user *t.User) (*t.User, error) {
	sess, c := r.C(ColUsers)
	defer sess.Close()

	var resp t.User
	user.ID = r.newID()
	err := m.Upsert(c, &resp, m.M{"_id": user.ID}, m.M{"$set": user})
	return &resp, err
}

func (r *MongoUserRepository) Delete(_ context.Context, id string) error {
	sess, c := r.C(ColUsers)
	defer sess.Close()

	var out t.User
	if err := m.FindOne(c, &out, m.M{"_id": id}); err != nil {
		return err
	}

	return m.Remove(c, m.M{"_id": id})
}

func (r *MongoUserRepository) newID() string {
	uid := uuid.New()
	return uid.String()
}
