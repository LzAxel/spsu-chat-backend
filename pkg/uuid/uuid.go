package uuid

import "github.com/google/uuid"

type UUIDI interface {
	New() uuid.UUID
	NewString() string
}

type RealUUID struct{}

func (u RealUUID) New() uuid.UUID    { return uuid.New() }
func (u RealUUID) NewString() string { return uuid.NewString() }

type DumbUUID struct{}

func (u DumbUUID) New() uuid.UUID    { return uuid.Nil }
func (u DumbUUID) NewString() string { return uuid.Nil.String() }

var (
	instance UUIDI = RealUUID{}
)

func InitUUID(isDumb bool) {
	if isDumb {
		instance = DumbUUID{}
	}
}

func New() uuid.UUID    { return instance.New() }
func NewString() string { return instance.NewString() }
