package models

import (
	"os"
	"reflect"
	"testing"
)

func TestModels(t *testing.T) {
	user1 := createUser(t, "test@test.com", "First", "Last")
	expected1 := User{1, "test@test.com", "First", "Last", 0}
	assertEqual(t, expected1, user1)

	user2 := createUser(t, "test2@test.com", "Luke", "Skywalker")
	expected2 := User{2, "test2@test.com", "Luke", "Skywalker", 0}
	assertEqual(t, expected2, user2)

	// Put 10 points on the first user.
	transfer1, err := NewTransfer(expected1.UserId, 10)
	if err != nil {
		t.Error(err)
	}
	expectedTransfer1 := Transfer{1, expected1.UserId, 10}
	assertEqual(t, expectedTransfer1, transfer1)

	user1, err = GetUser(expected1.UserId)
	if err != nil {
		t.Error(err)
	}
	expected1.Points = 10
	assertEqual(t, expected1, user1)

	// Try to take out 20 points.
	_, err = NewTransfer(expected1.UserId, -20)
	if err == nil {
		t.Errorf("Transfer should've failed")
	}

	user1, err = GetUser(expected1.UserId)
	if err != nil {
		t.Error(err)
	}
	assertEqual(t, expected1, user1)

	// Now take out only 10.
	transfer2, err := NewTransfer(expected1.UserId, -10)
	if err != nil {
		t.Error(err)
	}
	expectedTransfer2 := Transfer{2, expected1.UserId, -10}
	assertEqual(t, expectedTransfer2, transfer2)

	user1, err = GetUser(expected1.UserId)
	if err != nil {
		t.Error(err)
	}
	expected1.Points = 0
	assertEqual(t, expected1, user1)

	// Try to take out 20 points from the second user
	_, err = NewTransfer(expected2.UserId, -20)
	if err == nil {
		t.Errorf("Transfer should've failed")
	}

	user2, err = GetUser(expected2.UserId)
	if err != nil {
		t.Error(err)
	}
	assertEqual(t, expected2, user2)

}

func assertEqual(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected: %v\n.Actual: %v\n", expected, actual)
	}
}

func createUser(t *testing.T, email, firstName, lastName string) User {
	user, err := NewUser(email, firstName, lastName)
	if err != nil {
		t.Error(err)
	}
	return user
}

func TestMain(m *testing.M) {
	// TODO(tudor): Use different value?
	InitDb("postgres://tudor@localhost/button_test?sslmode=disable", true)
	defer CloseDb()
	os.Exit(m.Run())
}
