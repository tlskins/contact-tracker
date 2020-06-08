package mongo

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

var session *mgo.Session

func NewClient(host, user, pwd string) (*mgo.Session, error) {
	if user == "" {
		return mgo.Dial(host)
	}

	hostParts := strings.Split(host, "-")
	hostPre := hostParts[0]
	hostSuff := hostParts[1]
	hosts := []string{
		fmt.Sprintf("%s-shard-00-00-%s:27017", hostPre, hostSuff),
		fmt.Sprintf("%s-shard-00-01-%s:27017", hostPre, hostSuff),
		fmt.Sprintf("%s-shard-00-02-%s:27017", hostPre, hostSuff),
	}

	dialInfo := &mgo.DialInfo{
		Addrs:     hosts,
		Username:  user,
		Password:  pwd,
		PoolLimit: 10,
	}
	tlsConfig := &tls.Config{}
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}
	return mgo.DialWithInfo(dialInfo)
}

type IdGetter interface {
	GetId() interface{}
}

func Insert(c *mgo.Collection, result, data interface{}) error {
	var err error
	if result == nil {
		err = c.Insert(data)
	} else {
		change := mgo.Change{
			Update:    data,
			ReturnNew: true,
			Upsert:    true,
		}
		if ider, ok := data.(IdGetter); ok {
			_, err = c.Find(M{"_id": ider.GetId()}).Apply(change, result)
		} else {
			_, err = c.Find(M{"_id": bson.NewObjectId()}).Apply(change, result)
		}
	}
	return err
}

func Update(c *mgo.Collection, result, query, update interface{}) error {
	var err error
	if result == nil {
		err = c.Update(query, update)
	} else {
		change := mgo.Change{
			Update:    update,
			ReturnNew: true,
		}
		_, err = c.Find(query).Apply(change, result)
	}
	return err
}

func Upsert(c *mgo.Collection, result, query, update interface{}) error {
	var err error
	if result == nil {
		_, err = c.Upsert(query, update)
	} else {
		change := mgo.Change{
			Update:    update,
			Upsert:    true,
			ReturnNew: true,
		}
		_, err = c.Find(query).Apply(change, result)
		if p, ok := result.(PostProcessable); ok {
			p.PostProcess()
		}
	}
	return err
}

func UpdateAll(c *mgo.Collection, query, update interface{}) error {
	_, err := c.UpdateAll(query, update)
	return err
}

// First optional arg is Fields
// Second optional arg is slice of sort strings, ie. []string{"price", "-created_at"}
func Find(c *mgo.Collection, result, query interface{}, args ...interface{}) error {
	q := c.Find(query)
	if args != nil {
		if len(args) > 0 && args[0] != nil {
			q = q.Select(args[0])
		}
		if len(args) > 1 && args[1] != nil {
			q = q.Sort(args[1].([]string)...)
		}
	}
	if err := q.All(result); err != nil {
		return err
	}
	return nil
}

func FindOne(c *mgo.Collection, result, query interface{}, args ...interface{}) error {
	q := c.Find(query)
	if args != nil {
		if len(args) > 0 && args[0] != nil {
			q = q.Select(args[0])
		}
		if len(args) > 1 && args[1] != nil {
			q = q.Sort(args[1].([]string)...)
		}
	}
	if err := q.One(result); err != nil {
		return err
	}
	if p, ok := result.(PostProcessable); ok {
		p.PostProcess()
	}
	return nil
}

func Aggregate(c *mgo.Collection, result, pipe interface{}) error {
	if err := c.Pipe(pipe).All(result); err != nil {
		return err
	}
	return nil
}

func AggregateOne(c *mgo.Collection, result, pipe interface{}) error {
	if err := c.Pipe(pipe).One(result); err != nil {
		return err
	}
	return nil
}

func Remove(c *mgo.Collection, query interface{}) error {
	if _, err := c.RemoveAll(query); err != nil {
		return err
	}
	return nil
}

func CreateIndexKey(c *mgo.Collection, key ...string) error {
	return c.EnsureIndexKey(key...)
}
