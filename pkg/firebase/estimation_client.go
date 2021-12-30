package estimation_client

import (
    "fmt"
    "strconv"
    "context"
    "time"
	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
	firebase "firebase.google.com/go"
)

type Estimation struct {
    Duration int        `firestore:"duration"`
    Memo string         `firestore:"memo"`
    CreatedTm time.Time `firestore:"createdTm"`
    UpdatedTm time.Time `firestore:"updatedTm"`
}

const collectionName = "time_entry_estimations"

type EstimationClient struct {
    firestoreClient *firestore.Client
    firestoreCtx context.Context
}

func Init(serviceAccount option.ClientOption) (client *EstimationClient, err error) {
    var firestoreClient *firestore.Client
    var firestoreCtx = context.Background()

    var app *firebase.App
    app, err = firebase.NewApp(firestoreCtx, nil, serviceAccount)
    fmt.Printf("app: %v\n", app)
    fmt.Printf("firestoreCtx: %v\n", firestoreCtx)
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
    }
    return 
}

func (client *EstimationClient) Fetch(entryIds []int64) (estimations []Estimation) {
    // https://qiita.com/miyukiaizawa/items/88c174c00e9e99d3871b
    collectionRef := client.firestoreClient.Collection(collectionName)

    tmpDocs := make([]*firestore.DocumentRef, len(entryIds))
    for idx, id := range entryIds {
        tmpDocs[idx] = collectionRef.Doc(strconv.FormatInt(id, 10))
    }

    docSnaps, err := client.firestoreClient.GetAll(client.firestoreCtx, tmpDocs)
    if err != nil {
        // TODO error handling
    }
    for _, ds := range docSnaps {
        var estimation = Estimation{}
        ds.DataTo(&estimation)
        estimations = append(estimations, estimation)
    }
    return
}
