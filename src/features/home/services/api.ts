import type {
  User,
  Thread,
  ThreadDetail,
  Comment,
  PaginatedResponse,
  ThreadsQueryParams,
  Vote,
} from '@/types'
import { apiClient } from '@/lib/api/client'

// Extend existing ApiClient with microblog-specific methods
class MicroblogApi {
  // Auth
  // async register(name: string, email: string, password: string) {
  //   return apiClient.post<{ user: User }>('/register', {
  //     name,
  //     email,
  //     password,
  //   })
  // }

  // async login(email: string, password: string) {
  //   const response = await apiClient.post<{ token: string }>('/login', {
  //     email,
  //     password,
  //   })

  //   if (response.data?.token) {
  //     localStorage.setItem('token', response.data.token)
  //   }

  //   return response
  // }

  async getMe() {
    return apiClient.get<User>('/users/me')
  }

  // Threads
  async getAllThreads(params?: ThreadsQueryParams) {
    const queryParams = new URLSearchParams()

    if (params?.start !== undefined)
      queryParams.append('start', params.start.toString())
    if (params?.length !== undefined)
      queryParams.append('length', params.length.toString())
    if (params?.sort) queryParams.append('sort', params.sort)
    if (params?.fields) queryParams.append('fields', params.fields)
    if (params?.search) queryParams.append('search', params.search)

    const queryString = queryParams.toString()
    const url = queryString ? `/threads?${queryString}` : '/threads'

    // API returns PaginatedResponse directly (with data array inside)
    const response = await apiClient.get<Thread[]>(url)
    // Transform ApiResponse<Thread[]> to match PaginatedResponse structure
    return response as unknown as PaginatedResponse<Thread>
  }

  async createThread(title: string, body: string, category: string) {
    return apiClient.post<{ thread: Thread }>('/threads', {
      title,
      body,
      category,
    })
  }

  async getThreadDetail(threadId: string) {
    return apiClient.get<ThreadDetail>(`/threads/${threadId}`)
  }

  async deleteThread(threadId: string) {
    return apiClient.delete(`/threads/${threadId}`)
  }
  // Threads
  async getThreadAllComments(threadId: string, params?: ThreadsQueryParams) {
    const queryParams = new URLSearchParams()

    if (params?.start !== undefined)
      queryParams.append('start', params.start.toString())
    if (params?.length !== undefined)
      queryParams.append('length', params.length.toString())
    if (params?.sort) queryParams.append('sort', params.sort)
    if (params?.fields) queryParams.append('fields', params.fields)
    if (params?.search) queryParams.append('search', params.search)

    const queryString = queryParams.toString()
    const url = queryString ? `/threads/${threadId}/comments?${queryString}` : `/threads/${threadId}/comments`

    // API returns PaginatedResponse directly (with data array inside)
    const response = await apiClient.get<Comment[]>(url)
    // Transform ApiResponse<Comment[]> to match PaginatedResponse structure
    return response as unknown as PaginatedResponse<Comment>
  }

  // Comments
  async createComment(threadId: string, content: string) {
    return apiClient.post<{ comment: Comment }>(
      `/threads/${threadId}/comments`,
      { content }
    )
  }

  // Votes
  async upVoteThread(threadId: string) {
    return apiClient.post<Vote>(`/threads/${threadId}/up-vote`)
  }

  async downVoteThread(threadId: string) {
    return apiClient.post<Vote>(`/threads/${threadId}/down-vote`)
  }

  async neutralVoteThread(threadId: string) {
    return apiClient.post<Vote>(`/threads/${threadId}/neutral-vote`)
  }

  // Comment votes
  async upVoteComment(threadId: string, commentId: string) {
    return apiClient.post<Vote>(`/threads/${threadId}/comments/${commentId}/up-vote`)
  }

  async downVoteComment(threadId: string, commentId: string) {
    return apiClient.post<Vote>(`/threads/${threadId}/comments/${commentId}/down-vote`)
  }

  async neutralVoteComment(threadId: string, commentId: string) {
    return apiClient.post<Vote>(`/threads/${threadId}/comments/${commentId}/neutral-vote`)
  }

  logout() {
    localStorage.removeItem('token')
    localStorage.removeItem('firebase_id_token')
  }
}

export const api = new MicroblogApi()
