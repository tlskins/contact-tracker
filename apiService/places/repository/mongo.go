package repository

import (
	"context"
	"fmt"

	"github.com/globalsign/mgo"
	"github.com/google/uuid"

	m "github.com/contact-tracker/apiService/pkg/mongo"
	t "github.com/contact-tracker/apiService/places/types"
)

var ColPlaces = "places"

// MongoPlaceRepository -
type MongoPlaceRepository struct {
	session *mgo.Session
	dbname  string
}

// NewMongoPlaceRepository -
func NewMongoPlaceRepository(m *mgo.Session, dbname string) *MongoPlaceRepository {
	return &MongoPlaceRepository{m, dbname}
}

func (r *MongoPlaceRepository) DB() (*mgo.Session, *mgo.Database) {
	sess := r.session.Copy()
	return sess, sess.DB(r.dbname)
}

func (r *MongoPlaceRepository) C(colName string) (*mgo.Session, *mgo.Collection) {
	sess := r.session.Copy()
	return sess, sess.DB(r.dbname).C(colName)
}

func (r *MongoPlaceRepository) Get(_ context.Context, id string) (resp *t.Place, err error) {
	sess, c := r.C(ColPlaces)
	defer sess.Close()

	err = m.FindOne(c, &resp, m.M{"_id": id})
	return
}

func (r *MongoPlaceRepository) FindByEmail(_ context.Context, email string) (resp *t.Place, err error) {
	sess, c := r.C(ColPlaces)
	defer sess.Close()

	err = m.FindOne(c, &resp, m.M{"em": m.M{"$regex": fmt.Sprintf("(?i)^%s$", email)}})
	return
}

func (r *MongoPlaceRepository) GetAll(_ context.Context) (resp []*t.Place, err error) {
	sess, c := r.C(ColPlaces)
	defer sess.Close()

	err = m.Find(c, &resp, m.M{})
	return
}

func (r *MongoPlaceRepository) Update(_ context.Context, place *t.UpdatePlace) (*t.Place, error) {
	sess, c := r.C(ColPlaces)
	defer sess.Close()

	var resp t.Place
	err := m.Update(c, &resp, m.M{"_id": place.ID}, m.M{"$set": place})
	return &resp, err
}

func (r *MongoPlaceRepository) Create(_ context.Context, place *t.Place) (*t.Place, error) {
	sess, c := r.C(ColPlaces)
	defer sess.Close()

	var resp t.Place
	place.ID = r.newID()
	err := m.Upsert(c, &resp, m.M{"_id": place.ID}, m.M{"$set": place})
	return &resp, err
}

func (r *MongoPlaceRepository) Delete(_ context.Context, id string) error {
	sess, c := r.C(ColPlaces)
	defer sess.Close()

	var out t.Place
	if err := m.FindOne(c, &out, m.M{"_id": id}); err != nil {
		return err
	}

	return m.Remove(c, m.M{"_id": id})
}

func (r *MongoPlaceRepository) newID() string {
	uid := uuid.New()
	return uid.String()
}
