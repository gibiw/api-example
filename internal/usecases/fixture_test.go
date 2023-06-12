package usecases

import (
	"testing"

	"gihub.com/gibiw/api-example/internal/usecases/mocks"
	"github.com/golang/mock/gomock"
)

type Fixture struct {
	repository *mocks.Mockrepository
}

func NewFixture(t *testing.T) *Fixture {
	mockCtrl := gomock.NewController(t)
	repoMock := mocks.NewMockrepository(mockCtrl)

	return &Fixture{repository: repoMock}
}
