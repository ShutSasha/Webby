package services

import (
	"context"
	"fmt"
	"webby/admin-service/internal/models"

	"github.com/google/uuid"
)

type repository interface {
	ListComplaints(ctx context.Context, offset, limit int) ([]models.Complaint, int, error)
	ResolveComplaint(ctx context.Context, complaintID, userID uuid.UUID, isAccepted bool) error
}

type complaintGetter interface {
	GetComplaintByID(ctx context.Context, complaintID uuid.UUID) (*models.Complaint, error)
}

type videoInfoGetter interface {
	GetVideoByID(ctx context.Context, videoID uuid.UUID) (*models.Video, error)
	GetAuthorIDByVideoID(ctx context.Context, videoID string) (uuid.UUID, error)
}

type notificationSender interface {
	SendNotificationToUser(ctx context.Context, userID, targetID uuid.UUID, title, message string, complaintType string) error
}

type service struct {
	repository         repository
	complaintGetter    complaintGetter
	videoInfoGetter    videoInfoGetter
	notificationSender notificationSender
}

func New(repository repository, complaintGetter complaintGetter, videoInfoGetter videoInfoGetter, notificationSender notificationSender) *service {
	return &service{
		repository:         repository,
		complaintGetter:    complaintGetter,
		videoInfoGetter:    videoInfoGetter,
		notificationSender: notificationSender,
	}
}

func (s *service) ListComplaints(ctx context.Context, page, limit int) ([]models.Complaint, int, error) {
	const op = "service.ListComplaints"

	offset := (page - 1) * limit
	result, total, err := s.repository.ListComplaints(ctx, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	return result, total, nil
}

func (s *service) AcceptComplaint(ctx context.Context, complaintID, userID uuid.UUID) error {
	const op = "service.AcceptComplaint"

	complaint, err := s.complaintGetter.GetComplaintByID(ctx, complaintID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	violaterID := uuid.Nil
	message := ""
	switch complaint.TargetType {
	case "User":
		violaterID = complaint.TargetID
		message = "You were banned due to content restrictions"
	case "Video":
		video, err := s.videoInfoGetter.GetVideoByID(ctx, complaint.TargetID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		videoAuthorID, err := s.videoInfoGetter.GetAuthorIDByVideoID(ctx, video.ID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		violaterID = videoAuthorID
		message = fmt.Sprintf("Your video \"%s\" violates our platform rules. We have decided to ban it.", video.Title)
	}

	err = s.repository.ResolveComplaint(ctx, complaintID, userID, true)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.notificationSender.SendNotificationToUser(
		ctx,
		violaterID,
		complaint.TargetID,
		"Content violations",
		message,
		complaint.TargetType,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *service) DenyComplaint(ctx context.Context, complaintID, userID uuid.UUID, reason string) error {
	const op = "service.DenyComplaint"

	complaint, err := s.complaintGetter.GetComplaintByID(ctx, complaintID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.repository.ResolveComplaint(ctx, complaintID, userID, false)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.notificationSender.SendNotificationToUser(
		ctx,
		complaint.AuthorID,
		complaint.TargetID,
		"Content violations",
		reason,
		complaint.TargetType,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
