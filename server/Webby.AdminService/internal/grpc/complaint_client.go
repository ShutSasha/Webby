package grpc

import (
	"context"
	"fmt"
	"webby/admin-service/internal/grpc/complaintpb"
	"webby/admin-service/internal/models"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type complaintClient struct {
	client complaintpb.ComplaintGrpcServiceClient
	conn   *grpc.ClientConn
}

func NewComplaintClient(address string) (*complaintClient, error) {
	const op = "grpc.NewComplaintClient"

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect to complaint service: %w", op, err)
	}

	client := complaintpb.NewComplaintGrpcServiceClient(conn)

	return &complaintClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *complaintClient) Close() error {
	return c.conn.Close()
}

func (c *complaintClient) GetComplaintByID(ctx context.Context, complaintID uuid.UUID) (*models.Complaint, error) {
	const op = "grpc.complaintClient.GetComplaintByID"

	complaint, err := c.client.GetComplaintByID(ctx, &complaintpb.GetComplaintByIDRequest{
		ComplaintID: complaintID.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	authorID, err := uuid.Parse(complaint.AuthorId)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	targetID, err := uuid.Parse(complaint.TargetId)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	return &models.Complaint{
		ID: complaintID,
		Complainer: models.Complainer{
			ID: authorID,
		},
		Target: models.Target{
			Type: complaint.TargetType,
			ID:   targetID,
		},
		ReasonType:     complaint.ReasonType,
		AdditionalInfo: &complaint.AdditionalInfo,
		CreatedAt:      complaint.CreatedAt.AsTime(),
	}, nil
}
