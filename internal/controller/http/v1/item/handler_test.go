package item_test

import (
	"encoding/json"
	"errors"
	v1 "github.com/antoniokichaev/hezzl-collector/internal/controller/http/v1"
	"github.com/antoniokichaev/hezzl-collector/internal/controller/http/v1/item"
	"github.com/antoniokichaev/hezzl-collector/internal/entity/items"
	"github.com/antoniokichaev/hezzl-collector/internal/repo"
	"github.com/antoniokichaev/hezzl-collector/internal/repo/mocks"
	"github.com/antoniokichaev/hezzl-collector/internal/repo/pgdb"
	"github.com/antoniokichaev/hezzl-collector/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func setupServer(deps *service.SDependencies) *httptest.Server {

	repos := service.NewServices(deps, nil)
	eng := gin.Default()
	v1.New(eng, repos)

	return httptest.NewServer(eng)

}

func TestHandlers_GetItem(t *testing.T) {
	itemMocks := mocks.NewItem(t)
	campaignMocks := mocks.NewCampaign(t)
	deps := service.SDependencies{Repos: &repo.Repositories{
		Item:     itemMocks,
		Campaign: campaignMocks,
	}}
	srv := setupServer(&deps)
	defer srv.Close()

	tests := []struct {
		name          string
		args          []any
		returnValues  []any
		wantItem      items.Item
		body          items.Item
		queryUrlParam string
	}{
		// TODO: Add test cases.
		{
			name:          "valid",
			args:          []any{mock.Anything, 1, 1},
			returnValues:  []any{items.Item{ID: 1, CampaignID: 1, Priority: 1, Removed: false, Name: "item 1"}, nil},
			wantItem:      items.Item{ID: 1, CampaignID: 1, Priority: 1, Removed: false, Name: "item 1"},
			body:          items.Item{ID: 1, CampaignID: 1, Priority: 1, Removed: false, Name: "item 1"},
			queryUrlParam: "?campaignId=1&id=1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			itemMocks.On("GetItem", tt.args...).Return(tt.returnValues...).Once()
			request := resty.New().R()
			request.Method = http.MethodGet

			target, err := url.JoinPath(srv.URL, "/v1/item/")
			assert.NoError(t, err)
			request.URL = target + tt.queryUrlParam
			request.SetBody(tt.body)

			response, err := request.Send()
			assert.NoError(t, err)
			actual := items.Item{}
			err = json.Unmarshal(response.Body(), &actual)
			assert.NoError(t, err)
			assert.Equal(t, tt.body, actual)

		})
	}
}

func TestHandlers_Create(t *testing.T) {
	itemMocks := mocks.NewItem(t)
	campaignMocks := mocks.NewCampaign(t)
	deps := service.SDependencies{Repos: &repo.Repositories{
		Item:     itemMocks,
		Campaign: campaignMocks,
	}}
	srv := setupServer(&deps)
	defer srv.Close()

	tests := []struct {
		name          string
		args          []any
		returnValues  []any
		wantItem      items.Item
		wantErr       error
		body          item.RequestCreate
		queryUrlParam string
		statusCode    int
	}{
		// TODO: Add test cases.
		{
			name:          "valid",
			args:          []any{mock.Anything, "item 1", 1},
			returnValues:  []any{items.Item{ID: 1, CampaignID: 1, Priority: 1, Removed: false, Name: "item 1"}, nil},
			wantItem:      items.Item{ID: 1, CampaignID: 1, Priority: 1, Removed: false, Name: "item 1"},
			wantErr:       nil,
			body:          item.RequestCreate{Name: "item 1"},
			queryUrlParam: "?campaignId=1",
			statusCode:    http.StatusOK,
		},
		{
			name:          "bad companiId",
			args:          []any{mock.Anything, "item 1", 99},
			returnValues:  []any{items.Item{}, errors.New("err")},
			wantItem:      items.Item{},
			wantErr:       errors.New("err"),
			body:          item.RequestCreate{Name: "item 1"},
			queryUrlParam: "?campaignId=99",
			statusCode:    http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			itemMocks.On("CreateItem", tt.args...).Return(tt.returnValues...).Once()
			request := resty.New().R()
			request.Method = http.MethodPost

			target, err := url.JoinPath(srv.URL, "/v1/item/create/")
			assert.NoError(t, err)
			request.URL = target + tt.queryUrlParam
			request.SetBody(tt.body)

			response, err := request.Send()
			assert.NoError(t, err)
			assert.Equal(t, tt.statusCode, response.StatusCode())
			actual := items.Item{}
			err = json.Unmarshal(response.Body(), &actual)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantItem, actual)

		})
	}
}

func TestHandlers_Update(t *testing.T) {
	itemMocks := mocks.NewItem(t)
	campaignMocks := mocks.NewCampaign(t)
	deps := service.SDependencies{Repos: &repo.Repositories{
		Item:     itemMocks,
		Campaign: campaignMocks,
	}}
	srv := setupServer(&deps)
	defer srv.Close()

	tests := []struct {
		name          string
		args          []any
		returnValues  []any
		wantItem      items.Item
		body          item.RequestUpdate
		queryUrlParam string
		statusCode    int
	}{
		// TODO: Add test cases.
		{
			name:          "valid",
			args:          []any{mock.Anything, "item 1", "", 1, 1},
			returnValues:  []any{items.Item{ID: 1, CampaignID: 1, Priority: 1, Removed: false, Name: "item 1"}, nil},
			wantItem:      items.Item{ID: 1, CampaignID: 1, Priority: 1, Removed: false, Name: "item 1"},
			body:          item.RequestUpdate{RequestCreate: item.RequestCreate{Name: "item 1"}, Description: ""},
			queryUrlParam: "?campaignId=1&id=1",
			statusCode:    http.StatusOK,
		},
		{
			name:          "bad companiId",
			args:          []any{mock.Anything, "item 1", "", 1, 99},
			returnValues:  []any{items.Item{}, pgdb.ErrNotFoundItem},
			wantItem:      items.Item{},
			body:          item.RequestUpdate{RequestCreate: item.RequestCreate{Name: "item 1"}, Description: ""},
			queryUrlParam: "?campaignId=99&id=1",
			statusCode:    http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			itemMocks.On("UpdateItem", tt.args...).Return(tt.returnValues...).Once()
			request := resty.New().R()
			request.Method = http.MethodPatch

			target, err := url.JoinPath(srv.URL, "/v1/item/update/")
			assert.NoError(t, err)
			request.URL = target + tt.queryUrlParam
			request.SetBody(tt.body)

			response, err := request.Send()
			assert.NoError(t, err)
			assert.Equal(t, tt.statusCode, response.StatusCode())
			actual := items.Item{}
			err = json.Unmarshal(response.Body(), &actual)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantItem, actual)

		})
	}
}
