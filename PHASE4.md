# Phase 4: React Web Client

Build a Mattermost-style web client for Nexus-MM.

## Tech Stack
- React 18 + TypeScript
- Vite for build
- TailwindCSS for styling
- React Router for navigation
- Zustand for state management
- Axios for API calls

## Directory: web/

```
web/
  package.json
  vite.config.ts
  tsconfig.json
  tailwind.config.js
  index.html
  src/
    main.tsx
    App.tsx
    api/
      client.ts          - Axios instance with JWT interceptor
      types.ts           - TypeScript types matching Go models
      users.ts           - User API calls
      teams.ts           - Team API calls  
      channels.ts        - Channel API calls
      messages.ts        - Message API calls
      search.ts          - Search API
    store/
      auth.ts            - Auth state (Zustand)
      team.ts            - Current team/channel state
      messages.ts        - Messages state
    components/
      Layout.tsx         - Main app layout (sidebar + main)
      Sidebar/
        Sidebar.tsx      - Team/channel sidebar
        ChannelList.tsx   - Channel list with categories
        TeamSwitcher.tsx  - Team selector
      Chat/
        ChatView.tsx     - Main chat view
        MessageList.tsx  - Message list with infinite scroll
        MessageItem.tsx  - Single message (avatar, content, reactions, thread count)
        MessageInput.tsx - Message composer with slash commands
        ThreadPanel.tsx  - Thread/reply side panel
      Search/
        SearchBar.tsx    - Global search input
        SearchResults.tsx - Search results panel
      Auth/
        LoginPage.tsx    - Login form
        RegisterPage.tsx - Register form
      common/
        Avatar.tsx
        Button.tsx
        Modal.tsx
        Spinner.tsx
    hooks/
      useWebSocket.ts   - WuKongIM WS connection for real-time messages
    styles/
      globals.css       - Tailwind imports + custom styles
```

## UI Design (Mattermost-inspired)
- Dark sidebar (slate-800/900) with channel list
- White/light main content area
- Fixed header with channel name + search
- Message list with user avatars, timestamps
- Thread panel slides in from right
- Reactions shown below messages
- Pin indicator on pinned messages
- Channel categories collapsible in sidebar

## Key Features
1. Login/Register flow
2. Team switching
3. Channel list with categories (collapsible)
4. Real-time messaging (WebSocket to WuKongIM)
5. Thread/reply panel
6. Message reactions (click to add/remove)
7. Pin/unpin messages
8. Global search with results
9. Slash command autocomplete
10. Responsive design

## API Integration
- Base URL configurable via env: VITE_API_URL
- JWT stored in localStorage
- Axios interceptor adds Bearer token
- 401 responses redirect to login

## WebSocket
- Connect to WuKongIM WebSocket for real-time
- URL configurable: VITE_WS_URL
- Auto-reconnect on disconnect
- Incoming messages update Zustand store

## Build & Serve
- `npm run dev` for development
- `npm run build` for production
- Production build served by Go server at /web/

## Important
- Make it look professional and polished
- Dark mode sidebar, light content area
- Use proper TypeScript types everywhere
- Handle loading/error states
- No placeholder "TODO" - implement everything
- Make sure `npm run build` succeeds

When completely finished:
1. Run `npm run build` in web/ to verify
2. Commit and push to origin master
3. Run: openclaw system event --text "Done: Nexus-MM Phase 4 complete - React web client with real-time messaging, threads, search, reactions, pins" --mode now
