export interface User {
  id: string
  name: string
  email: string
  avatar?: string
}

export interface Thread {
  id: string
  title: string
  body: string
  category: string
  createdAt: string
  updatedAt: string
  owner_id: string
  owner: User
  comments?: Comment[]
  votes?: ThreadVote[]
  totalComments?: number
  upVotesBy?: string[]
  downVotesBy?: string[]
}

export interface Comment {
  id: string
  content: string
  createdAt: string
  updatedAt: string
  owner_id: string
  owner: User
  thread_id: string
  votes?: CommentVote[]
  upVotesBy?: string[]
  downVotesBy?: string[]
}

export interface ThreadDetail extends Thread {
  comments: Comment[]
}

export interface ThreadVote {
  id: string
  thread_id: string
  user_id: string
  vote_type: 'up' | 'down' | 'neutral'
}

export interface CommentVote {
  id: string
  comment_id: string
  user_id: string
  vote_type: 'up' | 'down' | 'neutral'
}

export interface LeaderboardUser extends User {
  score: number
}

export interface ApiResponse<T> {
  status: string
  message?: string
  data: T
}

export interface RegisterPayload {
  name: string
  email: string
  password: string
}

export interface LoginPayload {
  email: string
  password: string
}

export interface CreateThreadPayload {
  title: string
  body: string
  category: string
}

export interface CreateCommentPayload {
  content: string
}
