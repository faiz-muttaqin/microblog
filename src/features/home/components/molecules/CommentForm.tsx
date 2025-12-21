import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Send, Image, Paperclip, Smile, AtSign, Hash, Code } from 'lucide-react'

interface CommentFormProps {
    onSubmit: (content: string) => Promise<void>
    disabled?: boolean
}

export const CommentForm = ({ onSubmit, disabled }: CommentFormProps) => {
    const [content, setContent] = useState('')
    const [isSubmitting, setIsSubmitting] = useState(false)

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        if (!content.trim() || isSubmitting) return

        setIsSubmitting(true)
        try {
            await onSubmit(content.trim())
            setContent('')
        } finally {
            setIsSubmitting(false)
        }
    }

    return (
        <form onSubmit={handleSubmit} className="space-y-3">
            <Textarea
                placeholder="Write a comment..."
                value={content}
                onChange={(e) => setContent(e.target.value)}
                disabled={disabled || isSubmitting}
                className="min-h-24 resize-none rounded"
            />

            {/* Icon Tray */}
            <div className="flex items-center justify-between gap-1 px-3 py-2 bg-muted/50 rounded-md">
                <div>

                    <Button
                        variant="ghost"
                        size="sm"
                        className="h-8 w-8 p-0"
                        title="Add image"
                    >
                        <Image className="h-4 w-4" />
                    </Button>
                    <Button
                        variant="ghost"
                        size="sm"
                        className="h-8 w-8 p-0"
                        title="Attach file"
                    >
                        <Paperclip className="h-4 w-4" />
                    </Button>
                    <Button
                        variant="ghost"
                        size="sm"
                        className="h-8 w-8 p-0"
                        title="Add emoji"
                    >
                        <Smile className="h-4 w-4" />
                    </Button>
                    <Button
                        variant="ghost"
                        size="sm"
                        className="h-8 w-8 p-0"
                        title="Mention user"
                    >
                        <AtSign className="h-4 w-4" />
                    </Button>
                    <Button
                        variant="ghost"
                        size="sm"
                        className="h-8 w-8 p-0"
                        title="Add hashtag"
                    >
                        <Hash className="h-4 w-4" />
                    </Button>
                    <Button
                        variant="ghost"
                        size="sm"
                        className="h-8 w-8 p-0"
                        title="Insert code"
                    >
                        <Code className="h-4 w-4" />
                    </Button>
                </div>
                <Button
                    type="submit"
                    disabled={!content.trim() || disabled || isSubmitting}
                    size="sm"
                >
                    <Send className="h-4 w-4 mr-2" />
                    {isSubmitting ? 'Posting...' : 'Post Comment'}
                </Button>
            </div>
        </form>
    )
}
