package handlers

type handler struct {
	roomCreator          RoomCreator
	roomUpdater          RoomUpdater
	roomDeleter          RoomDeleter
	roomDetailsRetriever RoomDetailsRetriever
	roomLister           RoomLister
	roomMemberManager    RoomMemberManager
	playbackSynchronizer PlaybackSynchronizer
}

func New(
	roomCreator RoomCreator,
	roomUpdater RoomUpdater,
	roomDeleter RoomDeleter,
	roomDetailsRetriever RoomDetailsRetriever,
	roomLister RoomLister,
	roomMemberManager RoomMemberManager,
	playbackSynchronizer PlaybackSynchronizer,
) handler {
	return handler{
		roomCreator:          roomCreator,
		roomUpdater:          roomUpdater,
		roomDeleter:          roomDeleter,
		roomDetailsRetriever: roomDetailsRetriever,
		roomLister:           roomLister,
		roomMemberManager:    roomMemberManager,
		playbackSynchronizer: playbackSynchronizer,
	}
}
