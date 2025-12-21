export interface User {
  id: string
  name: string
  email: string
  avatar?: string
  first_name?: string
  last_name?: string
  username?: string
  phone_number?: string
  role?: string
  role_id?: number
  status?: string
  external_id?: string
  verification_status?: string
  session?: string
  last_login?: string
  created_at?: string
  updated_at?: string
}

export interface Thread {
  id: string
  title: string
  body: string
  category: string
  created_at: string
  updated_at: string
  user_id: string
  user: User
  total_comments: number
  total_up_votes: number
  total_down_votes: number
  up_voted_by_me: boolean
  down_voted_by_me: boolean
  comments?: Comment[]
  votes?: Vote[]
  upVotesBy?: string[]
  downVotesBy?: string[]
}

export interface Comment {
  id: string
  content: string
  createdAt: string
  updatedAt: string
  user_id: string
  user: User
  thread_id: string
  total_up_votes: number
  total_down_votes: number
  up_voted_by_me: boolean
  down_voted_by_me: boolean
  votes?: CommentVote[]
  upVotesBy?: string[]
  downVotesBy?: string[]
}

export interface ThreadDetail extends Thread {
  comments: Comment[]
}

export interface Vote {
  id: string
  thread_id?: string
  comment_id?: string
  user_id: string
  vote_type: 'up' | 'down' | 'neutral'
  total_up_votes?:   number,
  total_down_votes?: number,
  total_comments?:   number,
  up_voted_by_me?:   boolean,
  down_voted_by_me?: boolean,
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

export interface PaginatedResponse<T> {
  success: boolean
  draw?: number
  recordsTotal: number
  recordsFiltered: number
  data: T[]
}

export interface ThreadsQueryParams {
  start?: number
  length?: number
  sort?: string
  fields?: string
  search?: string
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
