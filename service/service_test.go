package service

import (
	"github.com/ninjasphere/app-presets/model"
	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/logger"
	nmodel "github.com/ninjasphere/go-ninja/model"
	"github.com/ninjasphere/go-ninja/rpc"
	"testing"
)

var saved = make([]*model.Presets, 0)

type mockConnection struct {
}

func (*mockConnection) ExportService(service interface{}, topic string, ann *nmodel.ServiceAnnouncement) (*rpc.ExportedService, error) {
	return nil, nil
}

func (*mockConnection) GetServiceClient(serviceTopic string) *ninja.ServiceClient {
	return nil
}

func makeService() (error, *PresetsService) {
	service := &PresetsService{
		Model: &model.Presets{
			Version: "1",
			Scenes: []*model.Scene{
				&model.Scene{
					ID:    "existing-uuid",
					Scope: "site:site-id",
				},
			},
		},
		Save: func(m *model.Presets) {
			saved = append(saved, m)
		},
		Conn: &mockConnection{},
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
	id := "not-existing-uuid"
	if scenes, err := s.FetchScenes(&model.Query{ID: &id}); err != nil {
		t.Fatalf("unexpected not found error: %v", err)
	} else if scenes == nil || len(*scenes) != 0 {
		t.Fatalf("scenes was nil or of length other than 0")
	}
}

func TestFetchScene(t *testing.T) {
	err, s := makeService()
	if err != nil {
		t.Fatalf("err was %v but expected nil", err)
	}
	id := "existing-uuid"
	if scene, err := s.FetchScenes(&model.Query{ID: &id}); err != nil {
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
	scope := "site:site-id"
	if scenes, err := s.FetchScenes(&model.Query{Scope: &scope}); err != nil || scenes == nil {
		t.Fatalf("err was: '%v', expected: nil", err)
	} else if len(*scenes) != 1 {
		t.Fatalf("number of results was: %d but expected: 1", len(*scenes))
	}
}

func TestStoreScene(t *testing.T) {
	err, s := makeService()
	if err != nil {
		t.Fatalf("err was %v but expected nil", err)
	}
	id := "existing-uuid"
	if scenes, err := s.FetchScenes(&model.Query{ID: &id}); err != nil {
		t.Fatalf("err was %v but expecting nil", err)
	} else if scenes == nil {
		t.Fatalf("scene was: nil, expected not nil")
	} else {
		saved = make([]*model.Presets, 0)
		_, err := s.StoreScene((*scenes)[0])
		if len(saved) == 0 || err != nil {
			t.Fatalf("save was not called")
		}
	}
}
