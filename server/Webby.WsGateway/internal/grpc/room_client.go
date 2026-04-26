package clients

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	roompb "webby-wsgateway/internal/grpc/roompb"
)

// RoomActivity is the gateway-internal representation of "who is online in which room"
// that we ship to RoomService periodically.
type RoomActivity struct {
	RoomID  string
	UserIDs []string
}

type RoomClient struct {
	conn *grpc.ClientConn
	c    roompb.RoomGrpcServiceClient
}

func NewRoomClient(addr string) (*RoomClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &RoomClient{conn: conn, c: roompb.NewRoomGrpcServiceClient(conn)}, nil
}

func (c *RoomClient) Close() error { return c.conn.Close() }

func (c *RoomClient) AwardActiveUsers(ctx context.Context, activity []RoomActivity) error {
	rooms := make([]*roompb.RoomActivity, 0, len(activity))
	for _, a := range activity {
		rooms = append(rooms, &roompb.RoomActivity{
			RoomId:  a.RoomID,
			UserIds: a.UserIDs,
		})
	}
	_, err := c.c.AwardActiveUsers(ctx, &roompb.AwardActiveUsersRequest{Rooms: rooms})
	return err
}
