package api

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/ccbrown/cloud-snitch/backend/app"
	"github.com/ccbrown/cloud-snitch/backend/app/apptest"
	"github.com/ccbrown/cloud-snitch/backend/model"
)

func TestMain(m *testing.M) {
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
	m.Run()
}

type TestAPI struct {
	*API
	app              *apptest.TestApp
	T                *testing.T
	AnonymousContext context.Context
}

func NewTestAPI(t *testing.T) *TestAPI {
	a := apptest.NewTestApp(t)
	ret := &TestAPI{
		API: New(a.App, Config{}),
		app: a,
		T:   t,
	}
	ret.AnonymousContext = context.WithValue(context.Background(), sessionContextKey, ret.app.NewAnonymousSession())
	return ret
}

func (a *TestAPI) NewTestUser(emailAddress string, role model.UserRole) (*model.User, context.Context) {
	user, sess := a.app.NewTestUser(emailAddress, role)
	ctx := context.WithValue(context.Background(), sessionContextKey, sess)
	return user, ctx
}

func (a *TestAPI) NewTestTeamWithSubscription(sess context.Context, tier app.TeamSubscriptionTier) *model.Team {
	team := a.app.NewTestTeamWithSubscription(ctxSession(sess), tier)
	return team
}

func TestGetHealthCheck(t *testing.T) {
	api := NewTestAPI(t)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/health-check", nil)
	require.NoError(t, err)
	api.ServeHTTP(w, r)

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "%v", string(body))
}
