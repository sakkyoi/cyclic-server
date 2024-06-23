package figleaf

import "github.com/matthewhartstonge/argon2"

type FigLeaf struct{}

func (f *FigLeaf) Cover(secret string) ([]byte, error) {
	argon := argon2.DefaultConfig()
	encoded, err := argon.HashEncoded([]byte(secret))
	if err != nil {
		return nil, err
	}

	return encoded, nil
}

func (f *FigLeaf) Peep(secret string, target []byte) (bool, error) {
	ok, err := argon2.VerifyEncoded([]byte(secret), target)
	if err != nil {
		return false, err
	}

	return ok, nil
}
