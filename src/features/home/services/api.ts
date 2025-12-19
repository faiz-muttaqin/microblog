import { apiClient } from '@/lib/api/client'
import type { User, Thread, ThreadDetail, Comment } from '@/types'

// Extend existing ApiClient with microblog-specific methods
class MicroblogApi {
  // Auth
  async register(name: string, email: string, password: string) {
    return apiClient.post<{ user: User }>('/register', { name, email, password })
  }

  async login(email: string, password: string) {
    const response = await apiClient.post<{ token: string }>('/login', { email, password })
    
    if (response.data?.token) {
      localStorage.setItem('token', response.data.token)
    }
    
    return response
  }

  async getMe() {
    return apiClient.get<{ user: User }>('/users/me')
  }

  // Threads
  async getAllThreads() {
    return apiClient.get<{ threads: Thread[] }>('/threads')
  }

  async createThread(title: string, body: string, category: string) {
    return apiClient.post<{ thread: Thread }>('/threads', { title, body, category })
  }

  async getThreadDetail(threadId: string) {
    return apiClient.get<{ detailThread: ThreadDetail }>(`/threads/${threadId}`)
  }

  async deleteThread(threadId: string) {
    return apiClient.delete(`/threads/${threadId}`)
  }

  // Comments
  async createComment(threadId: string, content: string) {
    return apiClient.post<{ comment: Comment }>(`/threads/${threadId}/comments`, { content })
  }

  // Votes
  async upVoteThread(threadId: string) {
    return apiClient.post(`/threads/${threadId}/up-vote`)
  }

  async downVoteThread(threadId: string) {
    return apiClient.post(`/threads/${threadId}/down-vote`)
  }

  async neutralVoteThread(threadId: string) {
    return apiClient.post(`/threads/${threadId}/neutral-vote`)
  }

  logout() {
    localStorage.removeItem('token')
    localStorage.removeItem('firebase_id_token')
  }
}

export const api = new MicroblogApi()
