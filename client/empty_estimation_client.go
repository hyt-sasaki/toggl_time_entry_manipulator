package client
import (
	"toggl_time_entry_manipulator/domain"
)

type emptyEstimationClient struct {
}


func (client *emptyEstimationClient) Fetch(entryIds []int64) (estimations []*domain.Estimation, err error) {
    for range entryIds {
        estimations = append(estimations, nil)
    }
    return
}

func (client *emptyEstimationClient) Insert(id string, estimation domain.Estimation) (err error) {
    return
}

func (client *emptyEstimationClient) Update(id string, estimation domain.Estimation) (err error) {
    return
}

func (client *emptyEstimationClient) Delete(id string) (err error) {
    return
}

func (client *emptyEstimationClient) Close() {
}
