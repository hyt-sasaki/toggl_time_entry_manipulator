package client

import (
	"testing"
	"toggl_time_entry_manipulator/config"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option"
)

func TestInit(t *testing.T) {
    // given
    serviceAccount := option.WithCredentialsFile("../credential/secret.json")
    config := config.FirestoreConfig{
        CollectionName: "test",
    }
    // when
    estimationClient, _ := NewEstimationClient(serviceAccount, config)
    // then
    assert.NotNil(t, estimationClient.firestoreClient)
    assert.NotNil(t, estimationClient.firestoreCtx)
}
