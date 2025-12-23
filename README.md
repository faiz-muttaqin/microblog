# Microblog - Social Media Platform

A full-stack social media application combining the best features of Twitter, Reddit, and Quora. Built with React + TypeScript + Go, featuring real-time interactions, customizable UI themes, and Firebase authentication.

🔗 **Live Demo**: [https://microblog.faizmuttaqin.com/](https://microblog.faizmuttaqin.com/)

## 📋 Table of Contents

- [Overview](#overview)
- [Tech Stack](#tech-stack)
- [Features](#features)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [API Documentation](#api-documentation)
- [Database Schema](#database-schema)
- [Configuration](#configuration)

---

## 🎯 Overview

**Microblog** is a modern social media platform that allows users to:
- Create and share threads (posts) with titles, body content, and categories
- Comment on threads with nested discussions
- Vote on threads and comments (upvote/downvote system)
- View leaderboards based on user activity and engagement
- Customize the entire UI (themes, fonts, colors, layouts)
- Real-time authentication with Firebase
- Dashboard with analytics (dummy data)
- Chat feature (placeholder for future implementation)

This project demonstrates a production-ready full-stack application with:
- **Frontend**: React 19 + TypeScript + TanStack Router + Vite
- **Backend**: Go + Gin + GORM + PostgreSQL/MySQL/SQLite
- **Authentication**: Firebase Auth with JWT tokens
- **UI Framework**: Tailwind CSS 4 + shadcn/ui + Radix UI
- **State Management**: TanStack Query for server state
- **Real-time Updates**: Optimistic UI updates

---

## 🛠 Tech Stack

### Frontend
- **React 19** - Latest React with concurrent features
- **TypeScript** - Type-safe development
- **Vite** - Lightning-fast build tool
- **TanStack Router** - Type-safe file-based routing
- **TanStack Query** - Server state management
- **TanStack Table** - Powerful table library
- **Tailwind CSS 4** - Utility-first CSS framework
- **shadcn/ui** - Re-usable component library
- **Radix UI** - Accessible component primitives
- **Motion** - Animation library (successor to Framer Motion)
- **Lucide React** - Icon library
- **Firebase SDK** - Authentication client
- **date-fns** - Date manipulation
- **Sonner** - Toast notifications
- **React Hook Form** - Form management
- **Zod** - Schema validation

### Backend
- **Go 1.24** - High-performance backend
- **Gin** - HTTP web framework
- **GORM** - ORM for database operations
- **Firebase Admin SDK** - Token verification
- **PostgreSQL/MySQL/SQLite** - Database support
- **Redis** - Caching layer (optional)
- **JWT** - Token-based authentication
- **CORS** - Cross-origin resource sharing
- **Lumberjack** - Log rotation
- **UUID** - Unique ID generation

### DevOps & Tools
- **Air** - Live reload for Go
- **Commitizen** - Conventional commits
- **ESLint** - JavaScript linting
- **Prettier** - Code formatting
- **Knip** - Unused code detection
- **Concurrently** - Run multiple commands

---

## ✨ Features

### 🔐 Authentication
- Firebase Authentication (Email/Password, Google, etc.)
- JWT token-based API authorization
- Protected routes with automatic redirects
- Persistent auth state with localStorage fallback
- Profile management

### 📝 Thread Management
- **Create threads** with title, body, and category
- **Edit/Delete** your own threads
- **Vote system** (upvote/downvote/neutral)
- **Optimistic UI updates** for instant feedback
- **Daily limit**: 100 threads per user per day (rate limiting)
- **Categories**: general, question, javascript, kotlin, python, react, etc.
- **Real-time vote counts** with database triggers

### 💬 Comments System
- **Nested comments** on threads
- **Edit/Delete** your own comments
- **Vote on comments** (upvote/downvote/neutral)
- **Optimistic updates** for comments
- **User attribution** with avatar and name
- **Timestamps** with relative time display

### 🏆 Leaderboard
- **Score calculation**: 
  - Thread/Comment votes: 5 points each
  - Comments created: 20 points each
- **User rankings** with avatar and profile info
- **Real-time updates** based on activity

### 🎨 UI Customization
- **Theme switching**: Light/Dark/System modes
- **Font selection**: 23+ Google Fonts with search and filter
  - Sans-serif: Inter, Roboto, Open Sans, Poppins, etc.
  - Serif: Merriweather, Playfair Display, Lora, etc.
  - Monospace: JetBrains Mono, Fira Code, Source Code Pro, etc.
- **Color customization**: Primary, background, accent colors
- **Layout options**: Sidebar position, spacing, radius
- **Real-time preview** of changes
- **Persistent settings** in localStorage
- **Advanced settings**: Font weights, line heights, animations

### 📊 Dashboard
- Analytics overview (placeholder with dummy data)
- User activity charts (Recharts integration)
- Quick stats and metrics

### 💬 Chat (Placeholder)
- UI structure for future chat implementation
- Message layout with responsive design
- User list with online indicators

### 🔧 Developer Features
- **Configurable table system** for CRUD operations
- **API client** with automatic token refresh
- **Error handling** with toast notifications
- **Loading states** with skeleton screens
- **Responsive design** for mobile/tablet/desktop
- **Accessibility** with ARIA labels and keyboard navigation
- **SEO-friendly** with proper meta tags

---

## 📁 Project Structure

```
microblog/
├── backend/                      # Go backend
│   ├── backend.go               # Server initialization & CORS
│   ├── docs.json                # Swagger documentation
│   ├── internal/
│   │   ├── database/            # Database connection & migrations
│   │   ├── handler/             # HTTP handlers
│   │   │   ├── google-fonts.go # Google Fonts API proxy
│   │   │   └── ...
│   │   ├── helper/              # Utility functions
│   │   │   └── firebase.go     # Firebase auth verification
│   │   ├── middleware/          # Custom middleware
│   │   ├── model/               # Database models (GORM)
│   │   │   ├── threads.go      # Thread, Comment, Vote models
│   │   │   ├── user.go          # User model
│   │   │   └── ...
│   │   └── routes/              # Route definitions
│   │       └── http_routes.go  # API endpoints
│   └── pkg/                     # Shared packages
│       ├── logger/              # Logging utilities
│       ├── util/                # Helper functions
│       └── ...
│
├── src/                         # React frontend
│   ├── assets/                  # Static assets (images, icons)
│   ├── components/              # Reusable UI components
│   │   ├── ui/                  # shadcn/ui components
│   │   ├── Header.tsx
│   │   ├── ProfileDropdown.tsx
│   │   └── ...
│   ├── config/                  # App configuration
│   │   ├── theme.ts            # Theme presets
│   │   └── fonts.ts            # Font definitions
│   ├── context/                 # React context providers
│   │   ├── theme-provider.tsx
│   │   ├── font-provider.tsx
│   │   └── ...
│   ├── features/                # Feature modules
│   │   ├── home/               # Thread list & creation
│   │   │   ├── components/
│   │   │   │   ├── ThreadCard.tsx
│   │   │   │   ├── CommentCard.tsx
│   │   │   │   ├── CreateThreadDialog.tsx
│   │   │   │   └── ...
│   │   │   ├── context/
│   │   │   └── services/
│   │   ├── auth/               # Authentication
│   │   ├── dashboard/          # Analytics dashboard
│   │   ├── settings/           # User settings
│   │   ├── chats/              # Chat feature (placeholder)
│   │   └── ...
│   ├── hooks/                   # Custom React hooks
│   │   ├── use-auth.ts         # Authentication hook
│   │   ├── use-toast.ts        # Toast notifications
│   │   └── ...
│   ├── lib/                     # Core utilities
│   │   ├── api/
│   │   │   ├── client.ts       # API client with auth
│   │   │   └── storage.ts      # LocalStorage helper
│   │   ├── firebase.ts         # Firebase config
│   │   └── utils.ts            # Helper functions
│   ├── routes/                  # TanStack Router routes
│   │   ├── __root.tsx          # Root layout
│   │   ├── (thread-app)/       # Thread feature routes
│   │   │   ├── index.tsx       # Home page
│   │   │   ├── threads/
│   │   │   │   └── $id.tsx     # Thread detail page
│   │   │   └── ...
│   │   └── ...
│   ├── styles/                  # Global styles
│   ├── types/                   # TypeScript types
│   └── main.tsx                # App entry point
│
├── .env.example                 # Environment variables template
├── package.json                 # Node dependencies
├── go.mod                       # Go dependencies
├── vite.config.ts              # Vite configuration
├── tailwind.config.ts          # Tailwind configuration
└── tsconfig.json               # TypeScript configuration
```

---

## 🚀 Getting Started

### Prerequisites
- **Node.js** 18+ and **pnpm** (or npm/yarn)
- **Go** 1.24+
- **PostgreSQL** / **MySQL** / **SQLite** (choose one)
- **Redis** (optional, for caching)
- **Firebase Project** with Authentication enabled

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/faiz-muttaqin/microblog.git
   cd microblog
   ```

2. **Install dependencies**
   ```bash
   # Frontend dependencies
   pnpm install

   # Backend dependencies
   go mod download
   ```

3. **Configure environment variables**
   ```bash
   cp .env.example .env
   ```

   Edit `.env` and add your credentials:
   ```env
   # Firebase Configuration (Frontend - Public)
   VITE_FIREBASE_API_KEY=your_api_key
   VITE_FIREBASE_AUTH_DOMAIN=your-project.firebaseapp.com
   VITE_FIREBASE_PROJECT_ID=your-project-id
   VITE_FIREBASE_APP_ID=your_app_id
   VITE_FIREBASE_STORAGE_BUCKET=your-project.appspot.com
   VITE_FIREBASE_MESSAGING_SENDER_ID=your_sender_id
   VITE_FIREBASE_MEASUREMENT_ID=your_measurement_id
   
   # Backend API URL (Frontend)
   VITE_BACKEND=/api
   VITE_BASE_PATH=/
   
   # Firebase Admin SDK (Backend - Private)
   FIREBASE_PRIVATE_KEY_JSON={"type":"service_account",...}
   SUPER_USER_EMAIL=admin@example.com
   
   # Database Configuration
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASS=your_password
   DB_NAME=microblog
   
   # Redis Configuration (Optional)
   REDIS_HOST=localhost
   REDIS_PORT=6379
   REDIS_PASSWORD=
   REDIS_DB=0
   
   # Google Fonts API (Optional)
   GOOGLE_FONTS_API_KEY=your_google_fonts_api_key
   
   # Server Configuration
   APP_LOCAL_HOST=:8173
   APP_GIN_MODE=release
   ```

4. **Run the application**

   **Development mode (Frontend + Backend)**:
   ```bash
   pnpm dev:all
   ```
   This runs:
   - Frontend on `http://localhost:18281` (Vite)
   - Backend on `http://localhost:8173` (Go with Air hot-reload)

   **Or run separately**:
   ```bash
   # Terminal 1: Frontend
   pnpm dev

   # Terminal 2: Backend
   air
   # or: go run main.go
   ```

5. **Build for production**
   ```bash
   pnpm build:all
   ```
   This creates:
   - Frontend build in `dist/`
   - Backend binary in `bin/`

6. **Run production build**
   ```bash
   ./bin/main
   # Backend serves the frontend from dist/
   ```

---

## 📚 API Documentation

### Authentication
All protected endpoints require `Authorization: Bearer <firebase_token>` header.

### Endpoints

#### **Users**
- `GET /api/users` - Get all users
- `GET /api/users/me` - Get current user profile

#### **Threads**
- `GET /api/threads` - Get all threads (with filters, pagination)
  - Query params: `category`, `search`, `sort`, `page`, `limit`
- `POST /api/threads` - Create new thread (max 100/day)
- `GET /api/threads/:id` - Get thread details with comments
- `PUT /api/threads/:id` - Update thread (owner only)
- `DELETE /api/threads/:id` - Delete thread (owner only)

#### **Thread Voting**
- `POST /api/threads/:id/up-vote` - Upvote thread
- `POST /api/threads/:id/down-vote` - Downvote thread
- `POST /api/threads/:id/neutral-vote` - Remove vote

#### **Comments**
- `GET /api/threads/:id/comments` - Get thread comments
- `POST /api/threads/:id/comments` - Create comment
- `PUT /api/threads/:id/comments/:commentId` - Update comment (owner only)
- `DELETE /api/threads/:id/comments/:commentId` - Delete comment (owner only)

#### **Comment Voting**
- `POST /api/threads/:id/comments/:commentId/up-vote`
- `POST /api/threads/:id/comments/:commentId/down-vote`
- `POST /api/threads/:id/comments/:commentId/neutral-vote`

#### **Leaderboard**
- `GET /api/leaderboards` - Get user rankings

#### **Google Fonts**
- `GET /api/google-fonts` - Get fonts list (cached 24h)
  - Query params: `q` (search), `category`, `limit`, `offset`

#### **Authentication**
- `GET /api/auth/login` - Verify Firebase token
- `GET /api/auth/logout` - Logout user
- `GET /api/auth/verify` - Verify authentication status

---

## 🗄 Database Schema

### Tables

#### `users`
- `id` (PK, UUID)
- `name`, `email`, `avatar`
- `created_at`, `updated_at`

#### `threads`
- `id` (PK, UUID)
- `title`, `body`, `category`
- `user_id` (FK → users)
- `total_up_votes`, `total_down_votes`, `total_comments`
- `created_at`, `updated_at`

#### `comments`
- `id` (PK, UUID)
- `thread_id` (FK → threads)
- `user_id` (FK → users)
- `content`
- `total_up_votes`, `total_down_votes`
- `created_at`, `updated_at`

#### `thread_votes`
- `id` (PK, UUID)
- `thread_id` (FK → threads)
- `user_id` (FK → users)
- `vote_type` (up/down/neutral)

#### `comment_votes`
- `id` (PK, UUID)
- `comment_id` (FK → comments)
- `user_id` (FK → users)
- `vote_type` (up/down/neutral)

### Database Triggers
- Auto-update `total_up_votes`, `total_down_votes` on vote changes
- Auto-update `total_comments` when comment is created/deleted

---

## ⚙️ Configuration

### CORS Configuration
Backend allows cross-origin requests with:
- All origins: `AllowOrigins: ["*"]`
- All headers: `AllowHeaders: ["*"]` (including Authorization)
- All methods: GET, POST, PUT, PATCH, DELETE, OPTIONS
- Cache preflight: 12 hours

For production, restrict origins:
```go
AllowOrigins: []string{
    "https://microblog.faizmuttaqin.com",
}
```

### Rate Limiting
- **Thread creation**: 100 threads per user per day (00:00-23:59)
- Enforced in `CreateThreadHandler` with database count query

### Caching
- **Google Fonts API**: In-memory cache with 24-hour TTL
- **Firebase tokens**: Stored in localStorage with auto-refresh

### Theme System
The app uses a sophisticated theme system with:
- **CSS Variables**: Dynamic theme colors
- **Font loading**: Google Fonts with fallback
- **localStorage**: Persistent user preferences
- **Real-time preview**: Changes apply instantly

---

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Commit Convention
This project uses Commitizen for conventional commits:
```bash
pnpm commit
```

---

## 📝 License

This project is open source and available under the [MIT License](LICENSE).

---

## 👨‍💻 Author

**Faiz Muttaqin**
- Website: [faizmuttaqin.com](https://faizmuttaqin.com)
- GitHub: [@faiz-muttaqin](https://github.com/faiz-muttaqin)
- Demo: [microblog.faizmuttaqin.com](https://microblog.faizmuttaqin.com)

---

## 🙏 Acknowledgments

- [shadcn/ui](https://ui.shadcn.com/) - Beautiful component library
- [Radix UI](https://www.radix-ui.com/) - Accessible primitives
- [TanStack](https://tanstack.com/) - Excellent React libraries
- [Firebase](https://firebase.google.com/) - Authentication service
- [Gin](https://gin-gonic.com/) - Go web framework
- [GORM](https://gorm.io/) - Go ORM library
