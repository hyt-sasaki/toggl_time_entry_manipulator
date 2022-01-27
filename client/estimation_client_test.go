package client


import (
    "testing"
    "time"
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
        CreatedTm: time.Now(),
        UpdatedTm: time.Now(),
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

func TestUpdate(t *testing.T) {
    // given
    id := "1"
    beforeUpdate, _ := estimationClient.Fetch([]int64{1})
    estimation := beforeUpdate[0]
    estimation.Duration = 10
    estimation.Memo = "updated memo"
    // when
    estimationClient.Update(id, *estimation)

    // then
    afterUpdate, _ := estimationClient.Fetch([]int64{1})
    assert.Equal(t, 1, len(afterUpdate))
    assert.Equal(t, estimation.Duration, afterUpdate[0].Duration)
    assert.Equal(t, estimation.Memo, afterUpdate[0].Memo)
    assert.Equal(t, estimation.CreatedTm, afterUpdate[0].CreatedTm)
    assert.NotEqual(t, estimation.UpdatedTm, afterUpdate[0].UpdatedTm)
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
