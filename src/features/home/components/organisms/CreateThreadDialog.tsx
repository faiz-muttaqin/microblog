import { useState } from 'react'
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select'
import { Separator } from '@/components/ui/separator'
import { Badge } from '@/components/ui/badge'
import { Switch } from '@/components/ui/switch'
import {
    ImageIcon,
    Link2,
    Smile,
    Code,
    AtSign,
    Hash,
    Video,
    FileText,
    Bold,
    Italic,
    List,
    AlignLeft,
    Settings2
} from 'lucide-react'
import { toast } from 'sonner'

interface CreateThreadDialogProps {
    open: boolean
    onClose: () => void
    onSubmit: (title: string, body: string, category: string) => void
}

const categories = [
    'general',
    'technology',
    'programming',
    'design',
    'business',
    'science',
    'entertainment',
    'sports',
    'other',
]

export const CreateThreadDialog = ({ open, onClose, onSubmit }: CreateThreadDialogProps) => {
    const [title, setTitle] = useState('')
    const [body, setBody] = useState('')
    const [category, setCategory] = useState('general')
    const [tags, setTags] = useState<string[]>([])
    const [tagInput, setTagInput] = useState('')
    const [allowComments, setAllowComments] = useState(true)
    const [notifyReplies, setNotifyReplies] = useState(true)
    const [showAdvanced, setShowAdvanced] = useState(false)

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault()
        if (title.trim() && body.trim()) {
            onSubmit(title, body, category)
            // Reset form
            setTitle('')
            setBody('')
            setCategory('general')
            setTags([])
            setTagInput('')
            setAllowComments(true)
            setNotifyReplies(true)
            setShowAdvanced(false)
            onClose()
        }
    }

    const handleAddTag = () => {
        if (tagInput.trim() && !tags.includes(tagInput.trim()) && tags.length < 5) {
            setTags([...tags, tagInput.trim()])
            setTagInput('')
        } else if (tags.length >= 5) {
            toast.error('Maximum 5 tags allowed')
        }
    }

    const handleRemoveTag = (tagToRemove: string) => {
        setTags(tags.filter(tag => tag !== tagToRemove))
    }

    const handleDummyAction = (action: string) => {
        toast.info(`${action} feature coming soon!`)
    }

    return (
        <Dialog open={open} onOpenChange={onClose}>
            <DialogContent className="sm:max-w-175 max-h-[90vh] overflow-y-auto">
                <DialogHeader>
                    <DialogTitle>Create New Thread</DialogTitle>
                </DialogHeader>

                <form onSubmit={handleSubmit} className="space-y-4">
                    {/* Title Input */}
                    <div className="space-y-2">
                        <Label htmlFor="title">Title *</Label>
                        <Input
                            id="title"
                            placeholder="What's on your mind?"
                            value={title}
                            onChange={(e) => setTitle(e.target.value)}
                            required
                            className="text-lg"
                        />
                    </div>



                    {/* Content Textarea */}
                    <div className="space-y-2">
                        <Label htmlFor="body">Content *</Label>

                        <Textarea
                            id="body"
                            placeholder="Share your thoughts... (Markdown supported)"
                            rows={8}
                            value={body}
                            onChange={(e) => setBody(e.target.value)}
                            required
                            className="resize-none"
                        />
                        {/* Formatting Toolbar (Dummy) */}
                        <div className="flex items-center justify-between">
                            <div className="bg-muted/30">
                                <div className="flex items-center gap-1 flex-wrap">
                                    <Button
                                        type="button"
                                        variant="ghost"
                                        size="sm"
                                        className="h-4 w-4 p-0"
                                        onClick={() => handleDummyAction('Bold')}
                                    >
                                        <Bold className="h-4 w-4" />
                                    </Button>
                                    <Button
                                        type="button"
                                        variant="ghost"
                                        size="sm"
                                        className="h-4 w-4 p-0"
                                        onClick={() => handleDummyAction('Italic')}
                                    >
                                        <Italic className="h-4 w-4" />
                                    </Button>
                                    <Button
                                        type="button"
                                        variant="ghost"
                                        size="sm"
                                        className="h-4 w-4 p-0"
                                        onClick={() => handleDummyAction('List')}
                                    >
                                        <List className="h-4 w-4" />
                                    </Button>
                                    <Button
                                        type="button"
                                        variant="ghost"
                                        size="sm"
                                        className="h-4 w-4 p-0"
                                        onClick={() => handleDummyAction('Align')}
                                    >
                                        <AlignLeft className="h-4 w-4" />
                                    </Button>
                                    <Separator orientation="vertical" className="h-6 mx-1" />
                                    <Button
                                        type="button"
                                        variant="ghost"
                                        size="sm"
                                        className="h-4 w-4 p-0"
                                        onClick={() => handleDummyAction('Code block')}
                                    >
                                        <Code className="h-4 w-4" />
                                    </Button>
                                    <Button
                                        type="button"
                                        variant="ghost"
                                        size="sm"
                                        className="h-4 w-4 p-0"
                                        onClick={() => handleDummyAction('Add link')}
                                    >
                                        <Link2 className="h-4 w-4" />
                                    </Button>
                                    <Button
                                        type="button"
                                        variant="ghost"
                                        size="sm"
                                        className="h-4 w-4 p-0"
                                        onClick={() => handleDummyAction('Mention user')}
                                    >
                                        <AtSign className="h-4 w-4" />
                                    </Button>
                                    <Separator orientation="vertical" className="h-6 mx-1" />
                                    <Button
                                        type="button"
                                        variant="ghost"
                                        size="sm"
                                        className="h-4 w-4 p-0"
                                        onClick={() => handleDummyAction('Upload image')}
                                    >
                                        <ImageIcon className="h-4 w-4" />
                                    </Button>
                                    <Button
                                        type="button"
                                        variant="ghost"
                                        size="sm"
                                        className="h-4 w-4 p-0"
                                        onClick={() => handleDummyAction('Attach video')}
                                    >
                                        <Video className="h-4 w-4" />
                                    </Button>
                                    <Button
                                        type="button"
                                        variant="ghost"
                                        size="sm"
                                        className="h-4 w-4 p-0"
                                        onClick={() => handleDummyAction('Attach file')}
                                    >
                                        <FileText className="h-4 w-4" />
                                    </Button>
                                    <Button
                                        type="button"
                                        variant="ghost"
                                        size="sm"
                                        className="h-4 w-4 p-0"
                                        onClick={() => handleDummyAction('Add emoji')}
                                    >
                                        <Smile className="h-4 w-4" />
                                    </Button>
                                </div>
                            </div>
                            <p className="text-xs text-muted-foreground">
                                {body.length} characters · {body.split(/\s+/).filter(w => w).length} words
                            </p>

                        </div>
                    </div>

                    {/* Category & Tags Row */}
                    <div className="flex flex-col sm:flex-row sm:items-center sm:gap-4 sm:justify-between">
                        <div className="space-y-2">
                            <Label htmlFor="category">Category *</Label>
                            <Select value={category} onValueChange={setCategory}>
                                <SelectTrigger>
                                    <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                    {categories.map((cat) => (
                                        <SelectItem key={cat} value={cat}>
                                            {cat.charAt(0).toUpperCase() + cat.slice(1)}
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                        </div>

                        <div className="flex-1 space-y-2">
                            <Label htmlFor="tags">Tags (Optional)</Label>
                            <div className="flex gap-2">
                                <Input
                                    id="tags"
                                    placeholder="Add tag..."
                                    value={tagInput}
                                    onChange={(e) => setTagInput(e.target.value)}
                                    onKeyDown={(e) => {
                                        if (e.key === 'Enter') {
                                            e.preventDefault()
                                            handleAddTag()
                                        }
                                    }}
                                />
                                <Button
                                    type="button"
                                    variant="outline"
                                    size="sm"
                                    onClick={handleAddTag}
                                >
                                    Add
                                </Button>
                            </div>
                        </div>
                    </div>

                    {/* Tags Display */}
                    {tags.length > 0 && (
                        <div className="flex flex-wrap gap-2">
                            {tags.map((tag) => (
                                <Badge
                                    key={tag}
                                    variant="secondary"
                                    className="cursor-pointer hover:bg-destructive/20"
                                    onClick={() => handleRemoveTag(tag)}
                                >
                                    <Hash className="h-3 w-3 mr-1" />
                                    {tag}
                                    <span className="ml-1 text-xs">×</span>
                                </Badge>
                            ))}
                        </div>
                    )}
                    {/* Advanced Settings Toggle */}
                    <Button
                        type="button"
                        variant="ghost"
                        size="sm"
                        onClick={() => setShowAdvanced(!showAdvanced)}
                        className="w-full"
                    >
                        <Settings2 className="h-4 w-4 mr-2" />
                        {showAdvanced ? 'Hide' : 'Show'} Advanced Settings
                    </Button>

                    {/* Advanced Settings (Dummy) */}
                    {showAdvanced && (
                        <div className="space-y-4 p-4 border rounded-md bg-muted/20">
                            <div className="flex items-center justify-between">
                                <div className="space-y-0.5">
                                    <Label htmlFor="comments">Allow Comments</Label>
                                    <p className="text-xs text-muted-foreground">
                                        Let others comment on your thread
                                    </p>
                                </div>
                                <Switch
                                    id="comments"
                                    checked={allowComments}
                                    onCheckedChange={setAllowComments}
                                />
                            </div>

                            <Separator />

                            <div className="flex items-center justify-between">
                                <div className="space-y-0.5">
                                    <Label htmlFor="notify">Notify on Replies</Label>
                                    <p className="text-xs text-muted-foreground">
                                        Get notified when someone replies
                                    </p>
                                </div>
                                <Switch
                                    id="notify"
                                    checked={notifyReplies}
                                    onCheckedChange={setNotifyReplies}
                                />
                            </div>

                            <Separator />

                            <div className="space-y-2">
                                <Label>Visibility (Coming Soon)</Label>
                                <Select defaultValue="public" disabled>
                                    <SelectTrigger>
                                        <SelectValue />
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectItem value="public">Public</SelectItem>
                                        <SelectItem value="followers">Followers Only</SelectItem>
                                        <SelectItem value="private">Private</SelectItem>
                                    </SelectContent>
                                </Select>
                            </div>
                        </div>
                    )}

                    <Separator />

                    {/* Action Buttons */}
                    <div className="flex justify-between items-center">
                        <p className="text-xs text-muted-foreground">
                            * Required fields
                        </p>
                        <div className="flex gap-2">
                            <Button type="button" variant="outline" onClick={onClose}>
                                Cancel
                            </Button>
                            <Button type="submit" disabled={!title.trim() || !body.trim()}>
                                Post Thread
                            </Button>
                        </div>
                    </div>
                </form>
            </DialogContent>
        </Dialog>
    )
}
