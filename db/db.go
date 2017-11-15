package db

import (
	"context"
	"strings"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type Model interface {
	EntityType() string
	GetKey() *datastore.Key
	SetKey(*datastore.Key) error
	PreSave(context.Context) error
	PostSave(context.Context) error
	PostLoad(context.Context) error
	PreDelete(context.Context) error
	Transform(context.Context, datastore.PropertyList) error
}

func Load(ctx context.Context, k *datastore.Key, m Model) (found bool, myerr error) {
	found = false

	if myerr = datastore.Get(ctx, k, m); myerr != nil {
		if ErrFieldMismatch(ctx, myerr, k, m) != nil {
			myerr = &UnfoundObjectError{
				EntityType: m.EntityType(),
				Key:        "key",
				Value:      k.Encode(),
				Err:        myerr,
			}
			return
		}
	}

	found = true

	if myerr = m.SetKey(k); myerr != nil {
		return
	}

	if myerr = m.PostLoad(ctx); myerr != nil {
		return
	}

	return
}

func LoadS(ctx context.Context, sk string, m Model) (found bool, myerr error) {
	found = false

	newKey, myerr := datastore.DecodeKey(sk)
	if myerr != nil {
		log.Infof(ctx, "LoadS Key: %s", sk)
		log.Infof(ctx, "LoadS Err: %v", myerr)
		return
	}

	found, myerr = Load(ctx, newKey, m)
	return
}

func LoadInt(ctx context.Context, id int64, m Model) (found bool, myerr error) {
	newKey := datastore.NewKey(ctx, m.EntityType(), "", id, nil)

	found, myerr = Load(ctx, newKey, m)
	return
}

func LoadMulti(ctx context.Context, keys []*datastore.Key, models []Model) (found int, myerr error) {
	found = 0

	if myerr = datastore.GetMulti(ctx, keys, models); myerr != nil {
		if ErrFieldMismatchMulti(ctx, myerr, keys, models) != nil {
			return
		}
	}

	for i, k := range keys {
		if myerr = models[i].SetKey(k); myerr != nil {
			return
		}

		if myerr = models[i].PostLoad(ctx); myerr != nil {
			return
		}
	}

	found = len(models)
	return
}

func LoadMultiS(ctx context.Context, skeys []string, models []Model) (found int, myerr error) {
	keys := make([]*datastore.Key, len(skeys))

	for i, sk := range skeys {
		keys[i], myerr = datastore.DecodeKey(sk)
		if myerr != nil {
			found = 0
			return
		}
	}

	found, myerr = LoadMulti(ctx, keys, models)
	return
}

func LoadMultiInt(ctx context.Context, intKeys []int64, entityType string, models []Model) (found int, myerr error) {
	keys := make([]*datastore.Key, len(intKeys))

	for i, ik := range intKeys {
		keys[i] = datastore.NewKey(ctx, entityType, "", ik, nil)
	}

	found, myerr = LoadMulti(ctx, keys, models)
	return
}

func Save(ctx context.Context, m Model) (myerr error) {
	if myerr = m.PreSave(ctx); myerr != nil {
		return
	}

	newKey, myerr := datastore.Put(ctx, m.GetKey(), m)
	if myerr != nil {
		log.Infof(ctx, "Save Err: %v", myerr)
		return
	}

	if myerr = m.SetKey(newKey); myerr != nil {
		return
	}

	if myerr = m.PostSave(ctx); myerr != nil {
		return
	}

	return
}

func SaveMulti(ctx context.Context, models []Model) (myerr error) {
	var keys []*datastore.Key
	for i := range models {
		if myerr = models[i].PreSave(ctx); myerr != nil {
			return
		}
		keys = append(keys, models[i].GetKey())
	}

	newKeys, myerr := datastore.PutMulti(ctx, keys, models)
	if myerr != nil {
		log.Infof(ctx, "Save Err: %v", myerr)
		return
	}

	for i := range models {
		if myerr = models[i].SetKey(newKeys[i]); myerr != nil {
			return
		}

		if myerr = models[i].PostSave(ctx); myerr != nil {
			return
		}
	}

	return
}

func Delete(ctx context.Context, m Model) (myerr error) {
	if myerr = m.PreDelete(ctx); myerr != nil {
		log.Infof(ctx, "PreDelete Err: %v", myerr)
		return
	}

	if myerr = datastore.Delete(ctx, m.GetKey()); myerr != nil {
		log.Infof(ctx, "Delete Err: %v", myerr)
		return
	}

	// if and when PostDelete gets built and/or is needed...
	//if myerr = m.PostDelete(ctx); myerr != nil {
	//	log.Infof(ctx, "PostDelete Err: %v", myerr)
	//	return
	//}

	return
}

func DeleteMultiK(ctx context.Context, keys []*datastore.Key) (myerr error) {
	// the datastore has a query limit, so chunk the query into small enough
	// chunks as to not trip the query limit error
	// chunk it into 0.5MB sizes to prevent query limit
	for 0 < len(keys) {
		chunk := make([]*datastore.Key, 0)
		size := 0
		n := 0
		for k, v := range keys {
			size += len(v.String())
			if size > 512000 {
				break
			}

			if len(keys) <= k-n {
				break
			}

			chunk = append(chunk, v)
			keys = append(keys[:k-n], keys[k-n+1:]...)
		}

		if myerr = datastore.DeleteMulti(ctx, chunk); myerr != nil {
			return
		}
	}

	return
}

func ErrFieldMismatchMulti(ctx context.Context, err error, keys []*datastore.Key, models []Model) (myerr error) {
	myerr, ok := err.(*datastore.ErrFieldMismatch)
	if ok || strings.Contains(err.Error(), "datastore: cannot load field") {
		tList := make([]datastore.PropertyList, len(keys))
		myerr = datastore.GetMulti(ctx, keys, tList)
		if myerr != nil {
			return
		}
		for i, propList := range tList {
			myerr = models[i].SetKey(keys[i])
			if myerr != nil {
				return
			}
			myerr = models[i].Transform(ctx, propList)
			if myerr != nil {
				return
			}
		}
	}
	return
}

func ErrFieldMismatchOnQuery(ctx context.Context, err error, keys []*datastore.Key, models []Model) (myerr error) {
	myerr, ok := err.(*datastore.ErrFieldMismatch)
	if ok || strings.ContainsAny(err.Error(), "datastore: cannot load field") {
		tList := make([]datastore.PropertyList, len(keys))
		myerr = datastore.GetMulti(ctx, keys, tList)
		if myerr != nil {
			return
		}
		for i, propList := range tList {
			myerr = models[i].SetKey(keys[i])
			if myerr != nil {
				return
			}
			myerr = models[i].Transform(ctx, propList)
			if myerr != nil {
				return
			}
		}
	}
	return
}

func ErrFieldMismatch(ctx context.Context, err error, k *datastore.Key, m Model) (myerr error) {
	myerr, ok := err.(*datastore.ErrFieldMismatch)
	if ok || strings.ContainsAny(err.Error(), "datastore: cannot load field") {
		var propList datastore.PropertyList
		myerr = datastore.Get(ctx, k, &propList)
		if myerr != nil {
			return
		}
		myerr = m.SetKey(k)
		if myerr != nil {
			return
		}
		myerr = m.Transform(ctx, propList)
	}
	return
}
