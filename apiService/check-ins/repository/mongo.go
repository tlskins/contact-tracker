package repository

import (
	"context"
	"time"

	"github.com/globalsign/mgo"

	t "github.com/contact-tracker/apiService/check-ins/types"
	m "github.com/contact-tracker/apiService/pkg/mongo"
)

var ColCheckIns = "checkIns"

// MongoCheckInRepository -
type MongoCheckInRepository struct {
	session *mgo.Session
	dbname  string
}

// NewMongoCheckInRepository -
func NewMongoCheckInRepository(m *mgo.Session, dbname string) *MongoCheckInRepository {
	return &MongoCheckInRepository{m, dbname}
}

func (r *MongoCheckInRepository) DB() (*mgo.Session, *mgo.Database) {
	sess := r.session.Copy()
	return sess, sess.DB(r.dbname)
}

func (r *MongoCheckInRepository) C(colName string) (*mgo.Session, *mgo.Collection) {
	sess := r.session.Copy()
	return sess, sess.DB(r.dbname).C(colName)
}

func (r *MongoCheckInRepository) Get(_ context.Context, id string) (resp *t.CheckIn, err error) {
	sess, c := r.C(ColCheckIns)
	defer sess.Close()

	err = m.FindOne(c, &resp, m.M{"_id": id})
	return
}

func (r *MongoCheckInRepository) GetHistory(_ context.Context, placeID string) (resp []*t.CheckInHistory, err error) {
	sess, c := r.C(ColCheckIns)
	defer sess.Close()

	resp = []*t.CheckInHistory{}
	if err = m.Find(c, &resp, m.M{"$match": m.M{"place.id": placeID}}); err != nil {
		return
	}
	for _, check := range resp {
		contacts := []*t.CheckIn{}
		if err = m.Find(c, &contacts, m.M{"$match": m.M{
			"place.id": placeID,
			"$or": []m.M{
				{"in": m.M{"$lte": check.Out}},
				{"out": m.M{"$gte": check.In}},
			},
		}}); err != nil {
			return
		}
		check.Contacts = contacts
	}
	return
}

func (r *MongoCheckInRepository) GetAll(_ context.Context, userID, placeID *string, start, end *time.Time) (resp []*t.CheckIn, err error) {
	sess, c := r.C(ColCheckIns)
	defer sess.Close()

	query := m.M{}
	if userID != nil {
		query["user.id"] = *userID
	}
	if placeID != nil {
		query["place.id"] = *placeID
	}
	if start != nil {
		query["in"] = m.M{"$gte": *start}
	}
	if end != nil {
		query["in"] = m.M{"$lt": *end}
	}

	err = m.Find(c, &resp, query)
	return
}

func (r *MongoCheckInRepository) LastCheckIn(_ context.Context, userID, placeID string) (*t.CheckIn, error) {
	sess, c := r.C(ColCheckIns)
	defer sess.Close()

	var resp []*t.CheckIn
	if err := c.Find(m.M{"place.id": placeID, "user.id": userID, "out": nil}).Sort("-in").Limit(1).All(&resp); err != nil {
		return nil, err
	}
	if len(resp) > 0 {
		return resp[0], nil
	}
	return nil, nil
}

func (r *MongoCheckInRepository) Create(_ context.Context, checkIn *t.CheckIn) (*t.CheckIn, error) {
	sess, c := r.C(ColCheckIns)
	defer sess.Close()

	var resp t.CheckIn
	err := m.Upsert(c, &resp, m.M{"_id": checkIn.ID}, m.M{"$set": checkIn})
	return &resp, err
}

func (r *MongoCheckInRepository) CheckOut(_ context.Context, id string) (*t.CheckIn, error) {
	sess, c := r.C(ColCheckIns)
	defer sess.Close()

	var resp t.CheckIn
	err := m.Update(c, &resp, m.M{"_id": id}, m.M{"$set": m.M{"out": time.Now()}})
	return &resp, err
}

func (r *MongoCheckInRepository) Delete(_ context.Context, id string) error {
	sess, c := r.C(ColCheckIns)
	defer sess.Close()

	var out t.CheckIn
	if err := m.FindOne(c, &out, m.M{"_id": id}); err != nil {
		return err
	}

	return m.Remove(c, m.M{"_id": id})
}
