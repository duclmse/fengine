package uuid

import (
	"github.com/duclmse/fengine/pkg/errors"
	"github.com/duclmse/fengine/viot"
	"github.com/google/uuid"
)

// ErrGeneratingID indicates error in generating UUID
var ErrGeneratingID = errors.New("generating id failed")

var _ viot.UUIDProvider = (*uuidProvider)(nil)

type uuidProvider struct{}

// New instantiates a UUID provider.
func New() viot.UUIDProvider {
	return &uuidProvider{}
}

func (up *uuidProvider) ID() (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", errors.Wrap(ErrGeneratingID, err)
	}

	return id.String(), nil
}
