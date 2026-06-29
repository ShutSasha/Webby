package services

import (
	"context"
	"fmt"
	"webby/admin-service/internal/apperrors"
	"webby/admin-service/internal/models"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type complaintsRepository interface {
	ListComplaints(ctx context.Context, offset, limit int) ([]models.Complaint, int, error)
	ResolveComplaint(ctx context.Context, complaintID, userID uuid.UUID, isAccepted bool) error
	IsResolved(ctx context.Context, complaintID uuid.UUID) (bool, error)
}

type complaintGetter interface {
	GetComplaintByID(ctx context.Context, complaintID uuid.UUID) (*models.Complaint, error)
}

type mediaRetriever interface {
	GetVideoByID(ctx context.Context, videoID uuid.UUID) (*models.Video, error)
	GetVideosBatch(ctx context.Context, ids []string) (map[string]string, error)
	GetAuthorIDByVideoID(ctx context.Context, videoID string) (uuid.UUID, error)
}

type videoBanner interface {
	BanVideo(ctx context.Context, videoID uuid.UUID) error
}

type userManager interface {
	BanUser(ctx context.Context, userID, requestUserID uuid.UUID) error
	GetUsersByIDs(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID]models.Complainer, error)
}

type notificationSender interface {
	SendNotificationToUser(ctx context.Context, userID, targetID uuid.UUID, title, message string, complaintType string) error
}

type complaintsSErvice struct {
	repository         complaintsRepository
	complaintGetter    complaintGetter
	mediaRetriever     mediaRetriever
	videoBanner        videoBanner
	userManager        userManager
	notificationSender notificationSender
}

func NewComplaintsService(repository complaintsRepository, complaintGetter complaintGetter, mediaRetriever mediaRetriever, videoBanner videoBanner, userManager userManager, notificationSender notificationSender) *complaintsSErvice {
	return &complaintsSErvice{
		repository:         repository,
		complaintGetter:    complaintGetter,
		mediaRetriever:     mediaRetriever,
		videoBanner:        videoBanner,
		userManager:        userManager,
		notificationSender: notificationSender,
	}
}

func (s *complaintsSErvice) ListComplaints(ctx context.Context, page, limit int) ([]models.Complaint, int, error) {
	const op = "service.ListComplaints"
	const webbyVideoPrefix = "wb_"

	offset := (page - 1) * limit
	complaints, total, err := s.repository.ListComplaints(ctx, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	userIDsSet := make(map[uuid.UUID]struct{})
	videoIDsSet := make(map[string]struct{})
	for _, complaint := range complaints {
		userIDsSet[complaint.Complainer.ID] = struct{}{}

		switch complaint.Target.Type {
		case "User":
			userIDsSet[complaint.Target.ID] = struct{}{}
		case "Video":
			videoIDsSet[webbyVideoPrefix+complaint.Target.ID.String()] = struct{}{}
		}
	}

	userIDs := make([]uuid.UUID, 0, len(userIDsSet))
	for id := range userIDsSet {
		userIDs = append(userIDs, id)
	}

	videoIDs := make([]string, 0, len(videoIDsSet))
	for id := range videoIDsSet {
		videoIDs = append(videoIDs, id)
	}

	complainers := make(map[uuid.UUID]models.Complainer)
	videos := make(map[string]string)

	g, insideContext := errgroup.WithContext(ctx)
	g.Go(func() error {
		var localErr error
		complainers, localErr = s.userManager.GetUsersByIDs(insideContext, userIDs)
		if localErr != nil {
			return localErr
		}

		return nil
	})

	g.Go(func() error {
		var localErr error
		videos, localErr = s.mediaRetriever.GetVideosBatch(insideContext, videoIDs)
		if localErr != nil {
			return localErr
		}

		return nil
	})
	if err := g.Wait(); err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	for i := range complaints {
		complaints[i].Complainer.Username = complainers[complaints[i].Complainer.ID].Username

		if complaints[i].Target.Type == "User" {
			complaints[i].Target.Name = complainers[complaints[i].Target.ID].Username
		} else {
			complaints[i].Target.Name = videos[webbyVideoPrefix+complaints[i].Target.ID.String()]
		}
	}

	return complaints, total, nil
}

func (s *complaintsSErvice) AcceptComplaint(ctx context.Context, complaintID, userID uuid.UUID) error {
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
	switch complaint.Target.Type {
	case "User":
		err := s.userManager.BanUser(ctx, complaint.Target.ID, userID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		violaterID = complaint.Target.ID
		message = "You were banned due to content restrictions"
	case "Video":
		video, err := s.mediaRetriever.GetVideoByID(ctx, complaint.Target.ID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		videoAuthorID, err := s.mediaRetriever.GetAuthorIDByVideoID(ctx, video.ID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		violaterID = videoAuthorID
		message = fmt.Sprintf("Your video \"%s\" violates our platform rules. We have decided to ban it.", video.Title)

		err = s.videoBanner.BanVideo(ctx, complaint.Target.ID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	err = s.repository.ResolveComplaint(ctx, complaintID, userID, true)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.notificationSender.SendNotificationToUser(
		ctx,
		violaterID,
		complaint.Target.ID,
		"Content violations",
		message,
		complaint.Target.Type,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *complaintsSErvice) DenyComplaint(ctx context.Context, complaintID, userID uuid.UUID, reason string) error {
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
		complaint.Complainer.ID,
		complaint.Target.ID,
		"Content violations",
		fmt.Sprintf("Your complaint has been denied. Reason: %s", reason),
		complaint.Target.Type,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
