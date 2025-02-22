package tests

import (
	"os"
	"testing"

	"github.com/EmelinDanila/task-manager-api/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockDB struct for a mock object that implements the Database interface
type MockDB struct {
	mock.Mock
}

// Close mock implementation of the Close method
func (m *MockDB) Close() error {
	args := m.Called()
	return args.Error(0)
}

// GetDB mock implementation of the GetDB method
func (m *MockDB) GetDB() *gorm.DB {
	return nil // In tests, we don't use a real database
}

/*
// Test for checking environment variable loading
func TestLoadEnv(t *testing.T) {
	os.Clearenv()

	config.LoadEnvVars()

	expected := "127.0.0.1"
	actual := os.Getenv("DB_HOST")
	assert.Equal(t, expected, actual, "DB_HOST should match the expected value")
}
*/

// Test for checking interaction with a mock database object
func TestMockConnectDatabase(t *testing.T) {
	// Create a mock object for the database
	mockDB := new(MockDB)

	// Set up the mock object for the Close method
	mockDB.On("Close").Return(nil)

	// Check that the mock object is not nil
	assert.NotNil(t, mockDB, "MockDB should not be nil")

	// Call the Close method and check that it doesn't return an error
	err := mockDB.Close()
	assert.NoError(t, err, "Close should not return an error")

	// Check that the mock Close method was called
	mockDB.AssertExpectations(t)
}

// Test for checking real database connection
func TestConnectDatabase(t *testing.T) {
	// Load environment variables
	oldEnv := os.Getenv("GO_ENV")
	os.Setenv("GO_ENV", "test")
	//config.LoadEnvVars()

	// Attempt to connect to the database
	db, err := config.ConnectDatabase()
	assert.NoError(t, err, "ConnectDatabase should not return an error")

	// Check that the database object is not nil
	assert.NotNil(t, db, "DB should not be nil")

	// Close the database connection
	err = db.Close()
	assert.NoError(t, err, "Close should not return an error")

	os.Setenv("GO_ENV", oldEnv)
}
