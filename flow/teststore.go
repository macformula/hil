package flow

import (
	"context"

	"github.com/macformula/hil/config"
)

type Store struct {
	*config.FireDB
}

// NewStore returns a Store.
func NewStore() *Store {
	d := config.FirebaseDB()
	return &Store{
		FireDB: d,
	}
}

// Create a new BIN object
func (s *Store) Create(tag *Tag) error {
	if err := s.NewRef("tags/"+tag.ID).Set(context.Background(), tag); err != nil {
		return err
	}
	return nil
}

func (s *Store) Delete(tag *Tag) error {
	return s.NewRef("tags/" + tag.ID).Delete(context.Background())
}

func (s *Store) GetByID(tadID string) (*Tag, error) {
	tag := &Tag{}
	if err := s.NewRef("tags/"+tadID).Get(context.Background(), tag); err != nil {
		return nil, err
	}
	if tag.ID == "" {
		return nil, nil
	}
	return tag, nil
}

func (s *Store) Update(tagID string, m map[string]interface{}) error {
	return s.NewRef("tagID/"+tagID).Update(context.Background(), m)
}
