package client

import (
	"context"
	"strconv"
	"time"

	"toggl_time_entry_manipulator/config"
	"toggl_time_entry_manipulator/domain"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type EstimationClient struct {
    firestoreClient *firestore.Client
    firestoreCtx context.Context
    config config.FirestoreConfig
}

type IEstimationClient interface {
    Fetch([]int64) ([]*domain.Estimation, error)
    Insert(string, domain.Estimation) error
    Update(string, domain.Estimation) error
    Delete(string) error
    Close()
}

func NewEstimationClient(serviceAccount option.ClientOption, config config.FirestoreConfig) (client IEstimationClient, err error) {
    if serviceAccount == nil {
        client = &EmptyEstimationClient{}
        return
    }
    var firestoreClient *firestore.Client
    var firestoreCtx = context.Background()

    var app *firebase.App
    app, err = firebase.NewApp(firestoreCtx, nil, serviceAccount)
    if err != nil {
        return 
    }

    firestoreClient, err = app.Firestore(firestoreCtx)
    if err != nil {
        return
    }

    client = &EstimationClient{
        firestoreClient: firestoreClient,
        firestoreCtx: firestoreCtx,
        config: config,
    }
    return 
}

func (client *EstimationClient) Fetch(entryIds []int64) (estimations []*domain.Estimation, err error) {
    // https://qiita.com/miyukiaizawa/items/88c174c00e9e99d3871b
    collectionRef := client.firestoreClient.Collection(client.config.CollectionName)

    tmpDocs := make([]*firestore.DocumentRef, len(entryIds))
    for idx, id := range entryIds {
        tmpDocs[idx] = collectionRef.Doc(strconv.FormatInt(id, 10))
    }

    docSnaps, err := client.firestoreClient.GetAll(client.firestoreCtx, tmpDocs)
    for _, ds := range docSnaps {
        if ds.Exists() {
            var estimation = domain.Estimation{}
            if err := ds.DataTo(&estimation); err == nil {
                estimations = append(estimations, &estimation)
            }
        } else {
            estimations = append(estimations, nil)
        }
    }
    return
}

func (client *EstimationClient) Insert(id string, estimation domain.Estimation) (err error){
    _, err = client.firestoreClient.Collection(client.config.CollectionName).Doc(id).Set(client.firestoreCtx, estimation)

    return
}

func (client *EstimationClient) Update(id string, estimation domain.Estimation) (err error){
    estimation.UpdatedTm = time.Now()
    _, err = client.firestoreClient.Collection(client.config.CollectionName).Doc(id).Set(client.firestoreCtx, estimation)

    return
}

func (client *EstimationClient) Delete(id string) (err error){
    _, err = client.firestoreClient.Collection(client.config.CollectionName).Doc(id).Delete(client.firestoreCtx)
    return
}

func (client *EstimationClient) Close() {
    client.firestoreClient.Close()
}
