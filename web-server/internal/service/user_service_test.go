package service

import (
	"testing"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/model"
)

func TestRequireGlobalPermission_HasPermission(t *testing.T) {
	user := &model.User{
		PermManageAllUser: true,
		PermCreateUser:    false,
		PermCreateBucket:  true,
	}

	// Should pass for permissions user has
	if err := RequireGlobalPermission(user, "MANAGE_ALL_USER"); err != nil {
		t.Errorf("Expected no error for MANAGE_ALL_USER, got %v", err)
	}

	if err := RequireGlobalPermission(user, "CREATE_BUCKET"); err != nil {
		t.Errorf("Expected no error for CREATE_BUCKET, got %v", err)
	}

	// Should pass when checking multiple permissions user has
	if err := RequireGlobalPermission(user, "MANAGE_ALL_USER", "CREATE_BUCKET"); err != nil {
		t.Errorf("Expected no error for multiple permissions, got %v", err)
	}
}

func TestRequireGlobalPermission_MissingPermission(t *testing.T) {
	user := &model.User{
		PermManageAllUser: false,
		PermCreateUser:    false,
		PermCreateBucket:  true,
	}

	// Should fail for permissions user doesn't have
	err := RequireGlobalPermission(user, "CREATE_USER")
	if err == nil {
		t.Error("Expected error for missing CREATE_USER permission, got nil")
	}

	// Should fail if any permission is missing
	err = RequireGlobalPermission(user, "CREATE_BUCKET", "CREATE_USER")
	if err == nil {
		t.Error("Expected error when one permission is missing, got nil")
	}
}

func TestRequireGlobalPermission_ErrorMessage(t *testing.T) {
	user := &model.User{
		PermManageAllUser: false,
		PermCreateUser:    false,
		PermCreateBucket:  false,
	}

	err := RequireGlobalPermission(user, "MANAGE_ALL_USER")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedMsg := "You do not have the required permissions. This action requires the \"MANAGE_ALL_USER\" permission."
	if err.Error() != expectedMsg {
		t.Errorf("Expected message '%s', got '%s'", expectedMsg, err.Error())
	}
}
