package estimation_client

import (
    "testing"
	"google.golang.org/api/option"
    "github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
    // given
    serviceAccount := option.WithCredentialsFile("../credential/secret.json")
    // when
    estimationClient, _ := NewEstimationClient(serviceAccount)
    // then
    assert.NotNil(t, estimationClient.firestoreClient)
    assert.NotNil(t, estimationClient.firestoreCtx)
}
