'use client'

import { useState } from 'react'

import { useRoomStore } from '@/stores/room.store'
import Modal from '@/ui/components/shared/Modal'

import CreateNextVideoVote from './CreateNextVideoVote'
import CreateRightChoiceVote from './CreateRightChoiceVote'
import NextVideoVoteDetails from './NextVideoVoteDetails'
import RightChoiceVoteDetails from './RightChoiceVoteDetails'
import SelectVoteType from './SelectVoteType'
import VotesList from './VotesList'

type Props = {
  roomId: string
  isHost: boolean
}

export type VoteViewState =
  | 'LIST'
  | 'SELECT_TYPE'
  | 'CREATE_RIGHT_CHOICE'
  | 'CREATE_NEXT_VIDEO'
  | 'DETAILS_RIGHT_CHOICE'
  | 'DETAILS_NEXT_VIDEO'

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

  const handleSelectVote = (voteId: string, type: 'RIGHT_CHOICE' | 'NEXT_VIDEO') => {
    setSelectedVoteId(voteId)
    if (type === 'RIGHT_CHOICE') setView('DETAILS_RIGHT_CHOICE')
    if (type === 'NEXT_VIDEO') setView('DETAILS_NEXT_VIDEO')
  }

  return (
    <Modal isOpen={isVotesModalOpen} onClose={handleClose} modalClasses="max-w-[480px]">
      {view === 'LIST' && (
        <VotesList roomId={roomId} isHost={isHost} onViewChange={setView} onSelectVote={handleSelectVote} />
      )}

      {view === 'SELECT_TYPE' && <SelectVoteType onSelect={setView} onBack={() => setView('LIST')} />}

      {view === 'CREATE_NEXT_VIDEO' && (
        <CreateNextVideoVote
          roomId={roomId}
          onBack={() => setView('SELECT_TYPE')}
          onSuccess={() => setView('DETAILS_NEXT_VIDEO')}
        />
      )}

      {view === 'CREATE_RIGHT_CHOICE' && (
        <CreateRightChoiceVote
          roomId={roomId}
          onBack={() => setView('SELECT_TYPE')}
          onSuccess={() => setView('LIST')}
        />
      )}

      {view === 'DETAILS_NEXT_VIDEO' && <NextVideoVoteDetails roomId={roomId} onBack={() => setView('LIST')} />}
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
