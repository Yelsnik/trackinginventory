package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/Yelsnik/trackinginventory/db/mock"
	"github.com/Yelsnik/trackinginventory/token"
	"github.com/gin-gonic/gin"

	db "github.com/Yelsnik/trackinginventory/db/sqlc"
	"github.com/Yelsnik/trackinginventory/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func requireBodyMatchInventory(t *testing.T, body *bytes.Buffer, inventory db.Inventory) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotInventory db.Inventory
	err = json.Unmarshal(data, &gotInventory)

	fmt.Println("INVENTORY:", inventory, "GOTINVENTORY", gotInventory)
	require.NoError(t, err)
	require.Equal(t, inventory, gotInventory)
}

func randomInventory(ownerID int64) db.Inventory {
	return db.Inventory{
		ID:           util.RandomInt(1, 1000),
		Item:         util.RandomString(2),
		SerialNumber: util.RandomString(4),
		Price:        util.RandomInt(1, 1000),
		Owner:        ownerID,
	}
}

func TestCreateInventoryApi(t *testing.T) {
	user, _ := randomUser(t)
	inventory := randomInventory(user.ID)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"item":     inventory.Item,
				"serialno": inventory.SerialNumber,
				"price":    inventory.Price,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateInventoryParams{
					Item:         inventory.Item,
					SerialNumber: inventory.SerialNumber,
					Price:        inventory.Price,
					Owner:        inventory.Owner,
				}

				store.EXPECT().
					CreateInventory(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(inventory, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchInventory(t, recorder.Body, inventory)
			},
		},
		{
			name: "InvalidItem",
			body: gin.H{
				"item":     "",
				"serialno": inventory.SerialNumber,
				"price":    inventory.Price,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateInventory(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidSerialNo",
			body: gin.H{
				"item":     inventory.Item,
				"serialno": "",
				"price":    inventory.Price,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateInventory(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPrice",
			body: gin.H{
				"item":     inventory.Item,
				"serialno": inventory.SerialNumber,
				"price":    0,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateInventory(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"item":     inventory.Item,
				"serialno": inventory.SerialNumber,
				"price":    inventory.Price,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateInventory(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Inventory{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)

			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/inventory"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetInventoryApi(t *testing.T) {

	user, _ := randomUser(t)
	inventory := randomInventory(user.ID)

	testCases := []struct {
		name          string
		inventoryID   int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name:        "OK",
			inventoryID: inventory.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetInventory(gomock.Any(), gomock.Eq(inventory.ID)).
					Times(1).
					Return(inventory, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchInventory(t, recorder.Body, inventory)
			},
		},
	}

	fmt.Println("INventory", inventory)
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// start test server and send request
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/inventory/%d", tc.inventoryID)
			fmt.Println(tc.inventoryID, url)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}
}
