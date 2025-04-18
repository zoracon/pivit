package pivit

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/go-piv/piv-go/v2/piv"
	"github.com/pkg/errors"
)

type ResetOpts struct {
	// Pin new PIN to set for performing certain operations on the Yubikey
	Pin string
}

// ResetYubikey resets the yubikey, sets a new pin, and sets a random PIN unblock key
func ResetYubikey(yk Pivit, opts *ResetOpts) error {
	if err := yk.Reset(); err != nil {
		return errors.Wrap(err, "reset PIV applet")
	}

	if err := yk.SetPIN(piv.DefaultPIN, opts.Pin); err != nil {
		return errors.Wrap(err, "failed to change pin")
	}

	newManagementKey, err := RandomManagementKey()
	if err != nil {
		return errors.Wrap(err, "failed to generate random management key")
	}

	if err = yk.SetManagementKey(piv.DefaultManagementKey, *newManagementKey); err != nil {
		return errors.Wrap(err, "set new management key")
	}
	if err = yk.SetMetadata(*newManagementKey, &piv.Metadata{
		ManagementKey: newManagementKey,
	}); err != nil {
		return errors.Wrap(err, "failed to store new management key")
	}

	randomPuk, err := rand.Int(rand.Reader, big.NewInt(100_000_000))
	if err != nil {
		return errors.Wrap(err, "create new random puk")
	}

	if err = yk.SetPUK(piv.DefaultPUK, fmt.Sprintf("%08d", randomPuk)); err != nil {
		return errors.Wrap(err, "set new puk")
	}

	return nil
}
