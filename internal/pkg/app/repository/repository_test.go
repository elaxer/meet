package repository

import (
	"context"
	"testing"
)

type modelTest struct {
	added, updated bool
}

func (m *modelTest) BeforeAdd() {
	m.added = true
}
func (m *modelTest) BeforeUpdate() {
	m.updated = true
}

func (m *modelTest) Validate() error {
	return nil
}

func TestRepository_Add(t *testing.T) {
	r := &collectionRepository[*modelTest]{
		models: []*modelTest{{}, {}, {}},
	}

	m := &modelTest{}

	if err := r.Add(context.Background(), m); err != nil {
		t.Errorf("collectionRepository.Add(): %s", err)
	}

	if has := r.has(m); !has {
		t.Errorf("collectionRepository.has() = false, want true")
	}

	if !m.added {
		t.Errorf("modelTest.added = false, want true")
	}
}

func TestRepository_Update(t *testing.T) {
	m := &modelTest{}

	r := &collectionRepository[*modelTest]{
		models: []*modelTest{m},
	}

	if err := r.Update(m); err != nil {
		t.Errorf("collectionRepository.Add(): %s", err)
	}

	if !m.updated {
		t.Errorf("modelTest.updated = false, want true")
	}
}

func TestRepository_Remove(t *testing.T) {
	m := &modelTest{}

	r := &collectionRepository[*modelTest]{
		models: []*modelTest{
			{},
			{},
			m,
			{},
			{},
		},
	}

	if err := r.Remove(m); err != nil {
		t.Errorf("collectionRepository.Remove(): %s", err)
	}

	if r.has(m) {
		t.Errorf("collectionRepository.has() = true, want false")
	}

	if len(r.models) != 4 {
		t.Errorf("repository models length is %d, want %d", len(r.models), 4)
	}
}
