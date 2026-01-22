package service_test

import (
	"context"
	"testing"

	"gooji/internal/domain"
	"gooji/internal/repository/sqlite"
	"gooji/internal/service"
)

func setupTestService(t *testing.T) *service.VideoService {
	db, err := sqlite.NewDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	repo := sqlite.NewVideoRepository(db)
	return service.NewVideoService(repo)
}

func TestVideoService_ListPublished(t *testing.T) {
	svc := setupTestService(t)
	ctx := context.Background()

	// Create mix of published and pending
	v1 := domain.NewVideo("Published", "/v1.mp4", "video/mp4")
	v1.Publish()
	v2 := domain.NewVideo("Pending", "/v2.mp4", "video/mp4")

	svc.Create(ctx, v1)
	svc.Create(ctx, v2)

	published, err := svc.ListPublished(ctx)
	if err != nil {
		t.Fatalf("failed to list: %v", err)
	}
	if len(published) != 1 {
		t.Errorf("expected 1 published, got %d", len(published))
	}
}

func TestVideoService_ListPending(t *testing.T) {
	svc := setupTestService(t)
	ctx := context.Background()

	v1 := domain.NewVideo("Pending 1", "/v1.mp4", "video/mp4")
	v2 := domain.NewVideo("Pending 2", "/v2.mp4", "video/mp4")
	v3 := domain.NewVideo("Published", "/v3.mp4", "video/mp4")
	v3.Publish()

	svc.Create(ctx, v1)
	svc.Create(ctx, v2)
	svc.Create(ctx, v3)

	pending, err := svc.ListPending(ctx)
	if err != nil {
		t.Fatalf("failed to list: %v", err)
	}
	if len(pending) != 2 {
		t.Errorf("expected 2 pending, got %d", len(pending))
	}
}

func TestVideoService_Approve(t *testing.T) {
	svc := setupTestService(t)
	ctx := context.Background()

	video := domain.NewVideo("Test", "/test.mp4", "video/mp4")
	svc.Create(ctx, video)

	err := svc.Approve(ctx, video.ID)
	if err != nil {
		t.Fatalf("failed to approve: %v", err)
	}

	approved, _ := svc.GetByID(ctx, video.ID)
	if approved.Status != domain.StatusPublished {
		t.Errorf("expected Published, got %s", approved.Status)
	}
}

func TestVideoService_Reject(t *testing.T) {
	svc := setupTestService(t)
	ctx := context.Background()

	video := domain.NewVideo("Test", "/test.mp4", "video/mp4")
	svc.Create(ctx, video)

	err := svc.Reject(ctx, video.ID, "Inappropriate content")
	if err != nil {
		t.Fatalf("failed to reject: %v", err)
	}

	rejected, _ := svc.GetByID(ctx, video.ID)
	if rejected.Status != domain.StatusRejected {
		t.Errorf("expected Rejected, got %s", rejected.Status)
	}
	if rejected.RejectionReason != "Inappropriate content" {
		t.Errorf("expected reason, got '%s'", rejected.RejectionReason)
	}
}
