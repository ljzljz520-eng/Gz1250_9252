package people

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestAdminPeopleLifecycle(t *testing.T) {
	service := NewService(NewMemoryRepository(nil), FixedClock(time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)))
	handler := NewHandler(service)

	create := CreateInput{Name: "Ava", Phone: "13800000111", Email: "ava@example.test", Role: RolePhotographer, Status: StatusActive}
	response := requestJSON(handler, http.MethodPost, "/admin/people", create)
	if response.Code != http.StatusCreated {
		t.Fatalf("create status = %d", response.Code)
	}
	var created Person
	decodeResponse(t, response, &created)
	if created.ID != "person-001" || created.CreatedAt.Format(time.RFC3339) != "2026-01-02T03:04:05Z" {
		t.Fatalf("created person = %#v", created)
	}

	response = requestJSON(handler, http.MethodPatch, "/admin/people/"+created.ID, StatusInput{Status: StatusInactive})
	if response.Code != http.StatusOK {
		t.Fatalf("update status = %d", response.Code)
	}
	var updated Person
	decodeResponse(t, response, &updated)
	if updated.Status != StatusInactive {
		t.Fatalf("updated status = %q", updated.Status)
	}

	response = requestJSON(handler, http.MethodGet, "/admin/people", nil)
	if response.Code != http.StatusOK {
		t.Fatalf("list status = %d", response.Code)
	}
	var listed []Person
	decodeResponse(t, response, &listed)
	if len(listed) != 1 || listed[0].Phone != create.Phone || listed[0].Role != create.Role {
		t.Fatalf("listed people = %#v", listed)
	}

	response = requestJSON(handler, http.MethodDelete, "/admin/people/"+created.ID, nil)
	if response.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d", response.Code)
	}
	response = requestJSON(handler, http.MethodGet, "/admin/people", nil)
	var empty []Person
	decodeResponse(t, response, &empty)
	if len(empty) != 0 {
		t.Fatalf("people after delete = %#v", empty)
	}
}

func TestAdminPeopleListLoadsAllFixtureRoles(t *testing.T) {
	file, err := os.Open("../../fixtures/people.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	clock := FixedClock(time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC))
	initial, err := LoadFixture(file, clock)
	if err != nil {
		t.Fatal(err)
	}
	handler := NewHandler(NewService(NewMemoryRepository(initial), clock))
	response := requestJSON(handler, http.MethodGet, "/admin/people", nil)
	if response.Code != http.StatusOK {
		t.Fatalf("list status = %d", response.Code)
	}
	var listed []Person
	decodeResponse(t, response, &listed)
	if len(listed) != 4 {
		t.Fatalf("listed people = %#v", listed)
	}
	wantRoles := []Role{RolePhotographer, RolePhotoEditor, RoleMakeupArtist, RoleCustomerService}
	for index, person := range listed {
		if person.Role != wantRoles[index] || person.Phone == "" || person.CreatedAt.Format(time.RFC3339) != "2026-01-02T03:04:05Z" {
			t.Fatalf("listed person %d = %#v", index, person)
		}
	}
}

func TestDuplicateCustomerPhoneReturnsConflictAndDoesNotMutate(t *testing.T) {
	clock := FixedClock(time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC))
	service := NewService(NewMemoryRepository(nil), clock)
	handler := NewHandler(service)
	photographer := CreateInput{Name: "Photographer", Phone: "13800000999", Role: RolePhotographer, Status: StatusActive}
	if response := requestJSON(handler, http.MethodPost, "/admin/people", photographer); response.Code != http.StatusCreated {
		t.Fatalf("seed status = %d", response.Code)
	}

	customer := CreateInput{Name: "Customer", Phone: photographer.Phone, Role: RoleCustomerService, Status: StatusActive}
	response := requestJSON(handler, http.MethodPost, "/admin/people", customer)
	var problem map[string]string
	decodeResponse(t, response, &problem)
	if response.Code != http.StatusConflict {
		t.Errorf("conflict status = %d", response.Code)
	}
	if problem["code"] != "contact_conflict" {
		errorf := problem["code"]
		t.Errorf("conflict code = %q", errorf)
	}

	response = requestJSON(handler, http.MethodGet, "/admin/people", nil)
	var listed []Person
	decodeResponse(t, response, &listed)
	if len(listed) != 1 || listed[0].Name != photographer.Name || listed[0].Role != photographer.Role {
		t.Fatalf("people after conflict = %#v", listed)
	}
}

func TestAdminPeopleReturnsStableValidationCode(t *testing.T) {
	service := NewService(NewMemoryRepository(nil), FixedClock(time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)))
	response := requestJSON(NewHandler(service), http.MethodPost, "/admin/people", CreateInput{Name: "Alex", Phone: "13800000666", Role: Role("producer")})
	var problem map[string]string
	decodeResponse(t, response, &problem)
	if response.Code != http.StatusBadRequest || problem["code"] != "invalid_role" {
		t.Fatalf("validation response = %d %#v", response.Code, problem)
	}
}

func requestJSON(handler http.Handler, method string, path string, value any) *httptest.ResponseRecorder {
	var body bytes.Buffer
	if value != nil {
		_ = json.NewEncoder(&body).Encode(value)
	}
	request := httptest.NewRequest(method, path, &body)
	request = request.WithContext(context.Background())
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	return recorder
}

func decodeResponse(t *testing.T, response *httptest.ResponseRecorder, destination any) {
	t.Helper()
	if err := json.NewDecoder(response.Body).Decode(destination); err != nil {
		t.Fatalf("decode response: %v", err)
	}
}
