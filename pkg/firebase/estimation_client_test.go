package estimation_client


import (
    "testing"
	"google.golang.org/api/option"
)

var estimationClient *EstimationClient

func TestMain(m *testing.M) {
    serviceAccount := option.WithCredentialsFile("../../credential/secret.json")
    estimationClient, _ = Init(serviceAccount)

    m.Run()
}

func TestFetch(t *testing.T) {
    entryIds := [...] int64{2279929209}
    estimations := estimationClient.Fetch(entryIds[:])
    t.Log(len(estimations))
    t.Log(estimations[0])
}

