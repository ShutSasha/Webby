package clients

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	votespb "webby-wsgateway/internal/grpc/votespb"
)

type VotesClient struct {
	conn *grpc.ClientConn
	c    votespb.VotesGrpcServiceClient
}

func NewVotesClient(addr string) (*VotesClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &VotesClient{conn: conn, c: votespb.NewVotesGrpcServiceClient(conn)}, nil
}

func (c *VotesClient) Close() error { return c.conn.Close() }

func (c *VotesClient) CastVote(ctx context.Context, roomID, userID, voteID, choiceID string) error {
	_, err := c.c.CastVote(ctx, &votespb.CastVoteRequest{
		RoomId:   roomID,
		UserId:   userID,
		VoteId:   voteID,
		ChoiceId: choiceID,
	})
	return err
}
