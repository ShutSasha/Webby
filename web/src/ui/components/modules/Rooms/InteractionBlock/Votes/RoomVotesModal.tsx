'use client'

import { useState } from 'react'

import { useRoomStore } from '@/stores/room.store'
import Modal from '@/ui/components/shared/Modal'

import CreateRightChoiceVote from './CreateRightChoiceVote'
import RightChoiceVoteDetails from './RightChoiceVoteDetails'
import VotesList from './VotesList'

type Props = {
  roomId: string
  isHost: boolean
}

export type VoteViewState = 'LIST' | 'CREATE_RIGHT_CHOICE' | 'DETAILS_RIGHT_CHOICE'

export default function RoomVotesModal({ roomId, isHost }: Props) {
  const isVotesModalOpen = useRoomStore(state => state.isVotesModalOpen)
  const setVotesModalOpen = useRoomStore(state => state.setVotesModalOpen)

  const [view, setView] = useState<VoteViewState>('LIST')
  const [selectedVoteId, setSelectedVoteId] = useState<string | null>(null)

  const handleClose = () => {
    setVotesModalOpen(false)
    setTimeout(() => {
      setView('LIST')
      setSelectedVoteId(null)
    }, 300)
  }

  const handleSelectVote = (voteId: string, type: 'RIGHT_CHOICE') => {
    setSelectedVoteId(voteId)
    if (type === 'RIGHT_CHOICE') setView('DETAILS_RIGHT_CHOICE')
  }

  return (
    <Modal isOpen={isVotesModalOpen} onClose={handleClose} modalClasses="max-w-[480px]">
      {view === 'LIST' && (
        <VotesList
          roomId={roomId}
          isHost={isHost}
          onViewChange={setView}
          onSelectVote={id => handleSelectVote(id, 'RIGHT_CHOICE')}
        />
      )}

      {view === 'CREATE_RIGHT_CHOICE' && <CreateRightChoiceVote roomId={roomId} onBack={() => setView('LIST')} />}

      {view === 'DETAILS_RIGHT_CHOICE' && selectedVoteId && (
        <RightChoiceVoteDetails
          roomId={roomId}
          voteId={selectedVoteId}
          isHost={isHost}
          onBack={() => setView('LIST')}
        />
      )}
    </Modal>
  )
}
