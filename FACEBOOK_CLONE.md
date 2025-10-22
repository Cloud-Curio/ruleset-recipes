# Facebook Clone - Modern Social Media Platform

A fully-featured Facebook clone built with Next.js 14, React, TypeScript, and Tailwind CSS, featuring a modern, sleek, and professional design with rich UI components.

## 🎨 Screenshots

### News Feed
![Facebook Feed](https://github.com/user-attachments/assets/052f5599-304d-4e2f-9e07-057467772b38)

### User Profile
![Facebook Profile](https://github.com/user-attachments/assets/9fcd6207-a18c-4db9-a4cf-95b5d0f2153c)

### Friends Page
![Facebook Friends](https://github.com/user-attachments/assets/e275b6de-730c-4f59-8cc0-9db5b9fdc6ff)

### Photos Gallery
![Facebook Photos](https://github.com/user-attachments/assets/9cc21eaf-b054-45ce-a4f1-2e53f40ebb6d)

## ✨ Features

### Social Media Functionality
- 📱 **News Feed** - Scrollable feed with posts from friends
- 👤 **User Profiles** - Complete profile pages with cover photos, bio, and timeline
- 📸 **Stories** - Instagram-style stories component
- 💬 **Posts** - Create, like, comment, and share posts
- 🖼️ **Photo Gallery** - Albums and photo grid with lightbox view
- 👥 **Friends Management** - Friend suggestions, requests, and connections list
- 🔔 **Notifications** - Real-time notification badge
- 💬 **Messaging** - Chat interface placeholder
- 🔍 **Search** - Global search functionality

### UI Components (shadcn/ui inspired)
- **Card** - Flexible card components for posts and content
- **Avatar** - User avatar with fallback support
- **Button** - Multiple variants (default, outline, ghost, etc.)
- **Input** - Form inputs with modern styling
- **Dialog** - Modal dialogs for uploads and interactions
- **Separator** - Visual content separators
- **Navigation** - Responsive navbar with icons

### Design Features
- 🎨 Modern, clean, and professional design
- 🌙 Dark mode support (via next-themes)
- 📱 Fully responsive layout (mobile, tablet, desktop)
- ⚡ Fast performance with Next.js optimizations
- 🎭 Smooth animations with Framer Motion
- 🎯 Accessible components
- 🔄 Real-time interactions (likes, comments)

## 🛠️ Technology Stack

### Frontend
- **Next.js 14** - React framework with server-side rendering
- **React 18** - UI library
- **TypeScript** - Type safety
- **Tailwind CSS** - Utility-first styling
- **Framer Motion** - Animations
- **React Query** - Data fetching and caching
- **Radix UI** - Headless UI components
- **Lucide Icons** - Modern icon library
- **next-themes** - Dark mode support

### UI Component Library
Custom implementation of shadcn/ui patterns:
- Avatar with image fallback
- Buttons with multiple variants
- Cards for content display
- Dialogs for modals
- Input fields with validation
- Separators for visual hierarchy

## 📁 Project Structure

```
frontend/
├── src/
│   ├── components/
│   │   ├── ui/              # Base UI components (shadcn/ui style)
│   │   │   ├── avatar.tsx
│   │   │   ├── button.tsx
│   │   │   ├── card.tsx
│   │   │   ├── dialog.tsx
│   │   │   ├── input.tsx
│   │   │   └── separator.tsx
│   │   └── facebook/        # Facebook-specific components
│   │       ├── Navbar.tsx
│   │       ├── LeftSidebar.tsx
│   │       ├── RightSidebar.tsx
│   │       ├── Stories.tsx
│   │       ├── PostComposer.tsx
│   │       └── Post.tsx
│   ├── contexts/
│   │   └── AuthContext.tsx  # Authentication context
│   ├── lib/
│   │   └── utils.ts         # Utility functions
│   ├── pages/
│   │   ├── _app.tsx         # App wrapper
│   │   ├── index.tsx        # Landing page
│   │   ├── feed.tsx         # Main news feed
│   │   ├── profile.tsx      # User profile
│   │   ├── friends.tsx      # Friends management
│   │   └── photos.tsx       # Photo gallery
│   └── styles/
│       └── globals.css      # Global styles
├── next.config.js
├── tailwind.config.js
└── package.json
```

## 🚀 Quick Start

### Prerequisites
- Node.js 18+
- npm 9+

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd ruleset-recipes
   ```

2. **Install dependencies**
   ```bash
   npm install
   cd frontend && npm install
   ```

3. **Start the development server**
   ```bash
   cd frontend
   npm run dev
   ```

4. **Open your browser**
   Navigate to [http://localhost:3000](http://localhost:3000)

## 📄 Available Pages

- `/feed` - Main news feed with posts and stories
- `/profile` - User profile with timeline and info
- `/friends` - Friend suggestions, requests, and list
- `/photos` - Photo albums and gallery

## 🎨 Design Principles

### Modern & Sleek
- Clean white space
- Subtle shadows and borders
- Smooth transitions and animations
- Consistent spacing and typography

### Professional
- Polished UI components
- Attention to detail
- Accessible design
- Mobile-first approach

### Component-Driven
- Reusable UI components
- Consistent design language
- Type-safe props
- Composable architecture

## 🔧 Development

### Run Tests
```bash
npm run test
```

### Type Checking
```bash
npm run type-check
```

### Linting
```bash
npm run lint
```

### Build for Production
```bash
npm run build
```

## 📱 Responsive Design

The application is fully responsive and optimized for:
- **Mobile** (< 640px)
- **Tablet** (640px - 1024px)
- **Desktop** (> 1024px)

## 🌙 Dark Mode

Dark mode is supported throughout the application using `next-themes`. Users can toggle between light and dark themes, with system preference detection.

## 🎯 Future Enhancements

- [ ] Real-time notifications with WebSockets
- [ ] Video upload and playback
- [ ] Group functionality
- [ ] Marketplace integration
- [ ] Events and calendar
- [ ] Emoji reactions
- [ ] Story creation and viewing
- [ ] Advanced privacy settings
- [ ] Multi-language support

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Design inspired by Facebook
- UI components inspired by [shadcn/ui](https://ui.shadcn.com/)
- Icons from [Lucide](https://lucide.dev/)
- Images from [Picsum](https://picsum.photos/) and [DiceBear](https://dicebear.com/)

---

**Built with ❤️ using Next.js, React, and Tailwind CSS**
