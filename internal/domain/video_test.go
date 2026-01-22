package domain_test

import (
	"testing"

	"gooji/internal/domain"
)

func TestNewVideo(t *testing.T) {
	v := domain.NewVideo("Test Title", "test.mp4", "video/mp4")

	if v.ID == "" {
		t.Error("expected ID to be generated")
	}
	if v.Title != "Test Title" {
		t.Errorf("expected title 'Test Title', got '%s'", v.Title)
	}
	if v.Status != domain.StatusPending {
		t.Errorf("expected status Pending, got %s", v.Status)
	}
	if v.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestVideoPublish(t *testing.T) {
	v := domain.NewVideo("Test", "test.mp4", "video/mp4")

	err := v.Publish()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if v.Status != domain.StatusPublished {
		t.Errorf("expected status Published, got %s", v.Status)
	}
	if v.PublishedAt == nil {
		t.Error("expected PublishedAt to be set")
	}
}

func TestVideoReject(t *testing.T) {
	v := domain.NewVideo("Test", "test.mp4", "video/mp4")

	err := v.Reject("Inappropriate content")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if v.Status != domain.StatusRejected {
		t.Errorf("expected status Rejected, got %s", v.Status)
	}
	if v.RejectionReason != "Inappropriate content" {
		t.Errorf("expected rejection reason, got '%s'", v.RejectionReason)
	}
}

func TestVideoCannotPublishIfAlreadyRejected(t *testing.T) {
	v := domain.NewVideo("Test", "test.mp4", "video/mp4")
	v.Reject("Bad")

	err := v.Publish()
	if err == nil {
		t.Error("expected error when publishing rejected video")
	}
}
