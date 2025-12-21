export function getErrorMessage(error: unknown): string {
    if (!error) return 'An unexpected error occurred'
    if (typeof error === 'string') return error
    if (error instanceof Error) return error.message
    try {
        // Try to read common shapes like { message } or { error: string }
        const err = error as Record<string, unknown>
        if (typeof err.message === 'string') return err.message
        if (typeof err.error === 'string') return err.error
        if (typeof err.data === 'string') return err.data
    } catch {
      // ignore
    }
    return 'An unexpected error occurred'
}
