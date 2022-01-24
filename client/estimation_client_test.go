package client


import (
    "testing"
    "fmt"
    "os"
    "context"
	"cloud.google.com/go/firestore"
    "github.com/stretchr/testify/assert"
    "toggl_time_entry_manipulator/domain"
)

var estimationClient *EstimationClient

func TestMain(m *testing.M) {
    if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
        fmt.Println("FIRESTORE_EMULATOR_HOST must be set")
        os.Exit(1)
    }
    estimationClient = initTestClient()
    estimationClient.Insert("1", domain.Estimation{
        Duration: 30,
        Memo: "memo",
    })

    m.Run()

    estimationClient.Close()
}

func TestFetchWhenEntryIdsExist(t *testing.T) {
    // given
    entryIds := [...] int64{1}

    // when
    estimations, _ := estimationClient.Fetch(entryIds[:])

    // then
    assert.Equal(t, 1, len(estimations))

    estimation := estimations[0]
    assert.Equal(t, 30, estimation.Duration)
    assert.Equal(t, "memo", estimation.Memo)
}

func TestFetchWhenEntryIdsEmpty(t * testing.T) {
    // given
    entryIds := [...] int64{}

    // when
    estimations, _ := estimationClient.Fetch(entryIds[:])

    // then
    assert.Equal(t, 0, len(estimations))
}

func TestFetchWhenEntryIdsIncorrect(t * testing.T) {
    // given
    entryIds := [...] int64{3}      // does not exist

    // when
    estimations, _ := estimationClient.Fetch(entryIds[:])

    // then
    assert.Equal(t, 1, len(estimations))
    assert.Nil(t, estimations[0])
}

func TestFetchWhenEntryIdsIncorrect2(t * testing.T) {
    // given
    entryIds := [...] int64{1,3}

    // when
    estimations, _ := estimationClient.Fetch(entryIds[:])

    // then
    assert.Equal(t, 2, len(estimations))
    assert.NotNil(t, estimations[0])
    assert.Nil(t, estimations[1])
}

func initTestClient() (client *EstimationClient) {
    ctx := context.Background()
    fc, _ := firestore.NewClient(ctx, "test")
    client = &EstimationClient{
        firestoreClient: fc,
        firestoreCtx: ctx,
    }
    return
}
