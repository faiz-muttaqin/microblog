import { VoteButton } from '../atoms/VoteButton'

interface VoteActionsProps {
  upVotesCount: number
  downVotesCount: number
  isUpVoted: boolean
  isDownVoted: boolean
  onUpVote: () => void
  onDownVote: () => void
  disabled?: boolean
}

export const VoteActions = ({
  upVotesCount,
  downVotesCount,
  isUpVoted,
  isDownVoted,
  onUpVote,
  onDownVote,
  disabled
}: VoteActionsProps) => {
  const score = upVotesCount - downVotesCount

  return (
    <div className="flex items-center gap-1">
      <VoteButton
        isActive={isUpVoted}
        onClick={onUpVote}
        variant="up"
        disabled={disabled}
      />
      <span className="text-sm font-medium min-w-[2rem] text-center">{score}</span>
      <VoteButton
        isActive={isDownVoted}
        onClick={onDownVote}
        variant="down"
        disabled={disabled}
      />
    </div>
  )
}
