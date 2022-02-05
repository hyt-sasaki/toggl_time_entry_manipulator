package client
import (
	"toggl_time_entry_manipulator/domain"
)

type EmptyEstimationClient struct {
}


func (client *EmptyEstimationClient) Fetch(entryIds []int64) (estimations []*domain.Estimation, err error) {
    for range entryIds {
        estimations = append(estimations, nil)
    }
    return
}

func (client *EmptyEstimationClient) Insert(id string, estimation domain.Estimation) (err error) {
    return
}

func (client *EmptyEstimationClient) Update(id string, estimation domain.Estimation) (err error) {
    return
}

func (client *EmptyEstimationClient) Delete(id string) (err error) {
    return
}

func (client *EmptyEstimationClient) Close() {
}
