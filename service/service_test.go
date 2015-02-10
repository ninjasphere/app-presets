package service

import (
	"github.com/ninjasphere/app-presets/model"
	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/logger"
	"testing"
)

var saved = make([]*model.Presets, 0)

func makeService() (error, *PresetsService) {
	service := &PresetsService{
		Model: &model.Presets{
			Version: "1",
			Scenes: []*model.Scene{
				&model.Scene{
					ID:    "existing-uuid",
					Scope: "site/site-id",
				},
			},
		},
		Save: func(m *model.Presets) {
			saved = append(saved, m)
		},
		Conn: &ninja.Connection{},
		Log:  logger.GetLogger("mock"),
	}
	err := service.Init()
	return err, service
}

func TestEmptyInit(t *testing.T) {
	service := &PresetsService{}
	err := service.Init()
	if err == nil {
		t.Fatalf("err was nil but expected not nil")
	}
}

func TestGoodInit(t *testing.T) {
	err, _ := makeService()
	if err != nil {
		t.Fatalf("err was %v but expected nil", err)
	}
}

func TestFetchSceneNotFound(t *testing.T) {
	err, s := makeService()
	if err != nil {
		t.Fatalf("err was %v but expected nil", err)
	}
	if _, err := s.FetchScene("not-existing-uuid"); err == nil {
		t.Fatalf("expected not found error")
	} else if err.Error() != "No such scene: not-existing-uuid" {
		t.Fatalf("err was %s but expected %v", err.Error(), "No such scene: not-existing-uuid")
	}
}

func TestFetchScene(t *testing.T) {
	err, s := makeService()
	if err != nil {
		t.Fatalf("err was %v but expected nil", err)
	}
	if scene, err := s.FetchScene("existing-uuid"); err != nil {
		t.Fatalf("err was %v but expecting nil", err)
	} else if scene == nil {
		t.Fatalf("scene was: nil, expected not nil")
	}
}

func TestFetchScenes(t *testing.T) {
	err, s := makeService()
	if err != nil {
		t.Fatalf("err was %v but expected nil", err)
	}
	if scenes, err := s.FetchScenes("site/site-id"); err != nil {
		t.Fatalf("err was %v but expecting nil", err)
	} else if len(scenes) != 1 {
		t.Fatalf("number of results was: %d but expected: 1", len(scenes))
	}
}

func TestStoreSceneWithNilScope(t *testing.T) {
	err, s := makeService()
	if err != nil {
		t.Fatalf("err was %v but expected nil", err)
	}
	if scene, err := s.FetchScene("existing-uuid"); err != nil {
		t.Fatalf("err was %v but expecting nil", err)
	} else if scene == nil {
		t.Fatalf("scene was: nil, expected not nil")
	} else {
		saved = make([]*model.Presets, 0)
		scene.Scope = ""
		_, err := s.StoreScene(scene)
		if len(saved) != 0 || err == nil {
			t.Fatalf("save should not have been called")
		}
	}
}

func TestStoreScene(t *testing.T) {
	err, s := makeService()
	if err != nil {
		t.Fatalf("err was %v but expected nil", err)
	}
	if scene, err := s.FetchScene("existing-uuid"); err != nil {
		t.Fatalf("err was %v but expecting nil", err)
	} else if scene == nil {
		t.Fatalf("scene was: nil, expected not nil")
	} else {
		saved = make([]*model.Presets, 0)
		_, err := s.StoreScene(scene)
		if len(saved) == 0 || err != nil {
			t.Fatalf("save was not called")
		}
	}
}
