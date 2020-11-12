package repository

import (
	"context"
	"time"

	"github.com/globalsign/mgo"
	"github.com/google/uuid"

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

func (r *MongoCheckInRepository) GetHistory(_ context.Context, userID string, start, end *time.Time) (resp []*t.CheckInHistory, err error) {
	sess, c := r.C(ColCheckIns)
	defer sess.Close()

	match := m.M{}
	if len(userID) != 0 {
		match["user.id"] = userID
	}

	pipeline := []m.M{
		{"$match": match},
		{"$addFields": m.M{
			"out": m.M{
				"$ifNull": []interface{}{
					"$out",
					m.M{"$add": []interface{}{"$in", 60000 * 5}}, // add 5 mins to in if out is null
				},
			},
			"tentative": m.M{
				"$cond": m.M{
					"if":   m.M{"$eq": []interface{}{"$out", nil}},
					"then": true,
					"else": false,
				},
			},
		}},
	}

	if start != nil && end != nil {
		pipeline = append(pipeline, m.M{
			"$match": m.M{
				"$or": []m.M{
					m.M{"out": m.M{"$gte": start, "$lte": end}},
					m.M{"in": m.M{"$gte": start, "$lte": end}},
					m.M{"in": m.M{"$lte": start}, "out": m.M{"$gte": end}},
				},
			},
		})
	}

	resp = []*t.CheckInHistory{}
	if err = m.Aggregate(c, &resp, pipeline); err != nil {
		return
	}
	for _, check := range resp {
		contacts := []*t.CheckIn{}
		if err = m.Aggregate(c, &contacts, []m.M{
			{"$match": m.M{
				"user.id": m.M{"$ne": check.User.ID},
			}},
			{"$addFields": m.M{
				"out": m.M{
					"$ifNull": []interface{}{
						"$out",
						m.M{"$add": []interface{}{"$in", 60000 * 5}}, // add 5 mins to in if out is null
					},
				},
				"tentative": m.M{
					"$cond": m.M{
						"if":   m.M{"$eq": []interface{}{"$out", nil}},
						"then": true,
						"else": false,
					},
				},
			}},
			{"$match": m.M{
				"$or": []m.M{
					m.M{"out": m.M{"$gte": check.In, "$lte": check.Out}},
					m.M{"in": m.M{"$gte": check.In, "$lte": check.Out}},
					m.M{"in": m.M{"$lte": check.In}, "out": m.M{"$gte": check.Out}},
				},
			}},
		}); err != nil {
			return
		}
		check.Contacts = contacts
	}
	return
}

func (r *MongoCheckInRepository) GetAll(_ context.Context, userID *string, start, end *time.Time) (resp []*t.CheckIn, err error) {
	sess, c := r.C(ColCheckIns)
	defer sess.Close()

	query := m.M{}
	if userID != nil {
		query["user.id"] = *userID
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

func (r *MongoCheckInRepository) LastCheckIn(_ context.Context, userID string) (resp *t.CheckIn, err error) {
	sess, c := r.C(ColCheckIns)
	defer sess.Close()

	resp = &t.CheckIn{}
	err = m.FindOne(c, resp, m.M{"user.id": userID, "out": m.M{"$eq": nil}})
	return
}

func (r *MongoCheckInRepository) Create(_ context.Context, checkIn *t.CheckIn) (*t.CheckIn, error) {
	sess, c := r.C(ColCheckIns)
	defer sess.Close()

	if checkIn.ID == "" {
		checkIn.ID = uuid.New().String()
	}

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
