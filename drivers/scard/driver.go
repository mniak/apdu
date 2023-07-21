package scard

import (
	"errors"
	"io"
	"log"

	"github.com/ebfe/scard"
	"github.com/mniak/apdu/internal/noop"
)

type driver struct {
	context *scard.Context
	card    *scard.Card
	logger  *log.Logger
}

func New() (*driver, error) {
	cont, err := scard.EstablishContext()
	if err != nil {
		return nil, err
	}
	result := driver{
		context: cont,
		logger:  noop.Logger(),
	}
	return &result, nil
}

func (d *driver) LoggingTo(w io.Writer) *driver {
	d.logger = log.New(w, "[scard] ", 0)
	return d
}

func (d *driver) ListReaders() ([]string, error) {
	return d.context.ListReaders()
}

func (d *driver) Connect(reader string) error {
	card, err := d.context.Connect(reader, scard.ShareShared, scard.ProtocolAny)
	if err != nil {
		return err
	}
	d.card = card
	return nil
}

func (d driver) SendBytes(b []byte) ([]byte, error) {
	if d.card == nil {
		return nil, ErrNotConnected
	}

	d.logger.Printf("TPDU sent: %2X\n", b)
	r, err := d.card.Transmit(b)
	d.logger.Printf("TPDU received: %2X\n", r)
	return r, err
}

var ErrNotConnected = errors.New("not connected. Connect() should be called first")
