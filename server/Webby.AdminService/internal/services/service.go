package services

import (
	"context"
	"fmt"
	"webby/admin-service/internal/apperrors"
	"webby/admin-service/internal/models"

	"github.com/google/uuid"
)

type repository interface {
	ListComplaints(ctx context.Context, offset, limit int) ([]models.Complaint, int, error)
	ResolveComplaint(ctx context.Context, complaintID, userID uuid.UUID, isAccepted bool) error
	IsResolved(ctx context.Context, complaintID uuid.UUID) (bool, error)
}

type complaintGetter interface {
	GetComplaintByID(ctx context.Context, complaintID uuid.UUID) (*models.Complaint, error)
}

type mediaRetriever interface {
	GetVideoByID(ctx context.Context, videoID uuid.UUID) (*models.Video, error)
	GetAuthorIDByVideoID(ctx context.Context, videoID string) (uuid.UUID, error)
}

type videoBanner interface {
	BanVideo(ctx context.Context, videoID uuid.UUID) error
}

type userBanner interface {
	BanUser(ctx context.Context, userID uuid.UUID) error
}

type notificationSender interface {
	SendNotificationToUser(ctx context.Context, userID, targetID uuid.UUID, title, message string, complaintType string) error
}

type service struct {
	repository         repository
	complaintGetter    complaintGetter
	mediaRetriever     mediaRetriever
	videoBanner        videoBanner
	userBanner         userBanner
	notificationSender notificationSender
}

func New(repository repository, complaintGetter complaintGetter, mediaRetriever mediaRetriever, videoBanner videoBanner, userBanner userBanner, notificationSender notificationSender) *service {
	return &service{
		repository:         repository,
		complaintGetter:    complaintGetter,
		mediaRetriever:     mediaRetriever,
		videoBanner:        videoBanner,
		userBanner:         userBanner,
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

	isAlreadyResolved, err := s.repository.IsResolved(ctx, complaintID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if isAlreadyResolved {
		return fmt.Errorf("%s: %w", op, apperrors.ErrAlreadyResolved)
	}

	complaint, err := s.complaintGetter.GetComplaintByID(ctx, complaintID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	violaterID := uuid.Nil
	message := ""
	switch complaint.TargetType {
	case "User":
		err := s.userBanner.BanUser(ctx, complaint.TargetID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		violaterID = complaint.TargetID
		message = "You were banned due to content restrictions"
	case "Video":
		err := s.videoBanner.BanVideo(ctx, complaint.TargetID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		video, err := s.mediaRetriever.GetVideoByID(ctx, complaint.TargetID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		videoAuthorID, err := s.mediaRetriever.GetAuthorIDByVideoID(ctx, video.ID)
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

	isAlreadyResolved, err := s.repository.IsResolved(ctx, complaintID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if isAlreadyResolved {
		return fmt.Errorf("%s: %w", op, apperrors.ErrAlreadyResolved)
	}

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
		fmt.Sprintf("Your complaint has been denied. Reason: %s", reason),
		complaint.TargetType,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
