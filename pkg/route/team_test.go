package route

import (
	"defense-dashboard/model"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestTeamHandlerBindsIDAsData(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	sqlDB.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = sqlDB.Close() })
	if err := db.AutoMigrate(&model.Team{}, &model.Event{}); err != nil {
		t.Fatal(err)
	}
	for _, team := range []model.Team{{ID: 1, Name: "Alpha"}, {ID: 2, Name: "Beta"}} {
		if err := db.Create(&team).Error; err != nil {
			t.Fatal(err)
		}
	}
	if err := db.Create(&model.Event{TeamID: 2, Point: 7, Log: "Beta event"}).Error; err != nil {
		t.Fatal(err)
	}

	router := gin.New()
	router.SetHTMLTemplate(template.Must(template.New("index.html").Parse("{{ .team }}|{{ .score }}|{{ .logs }}")))
	router.GET("/team/:id", Controller{DB: db}.TeamHandler)

	for _, tc := range []struct {
		name string
		id   string
		code int
		body string
	}{
		{"first team", "1", http.StatusOK, "Alpha|0|"},
		{"second team and events", "2", http.StatusOK, "Beta|7|"},
		{"missing team", "999", http.StatusNotFound, "Team not found"},
		{"SQL boolean expression", "1 OR 1=1", http.StatusNotFound, "Team not found"},
		{"SQL column expression", "name = 'Beta'", http.StatusNotFound, "Team not found"},
		{"SQL comment", "1 --", http.StatusNotFound, "Team not found"},
		{"non-numeric ID", "missing", http.StatusNotFound, "Team not found"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/team/"+url.PathEscape(tc.id), nil))
			if response.Code != tc.code {
				t.Fatalf("status = %d, want %d; body = %q", response.Code, tc.code, response.Body.String())
			}
			if !strings.Contains(response.Body.String(), tc.body) {
				t.Fatalf("body = %q, want substring %q", response.Body.String(), tc.body)
			}
			if tc.id == "2" && !strings.Contains(response.Body.String(), "Beta event") {
				t.Fatal("team events were not loaded")
			}
		})
	}
}
