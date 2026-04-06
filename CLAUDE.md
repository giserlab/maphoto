# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Maphoto** (also called Transfeed) is a map-based photo sharing application. Users can create places on a map and attach photos to them. It supports both anonymous viewing (public map) and authenticated management (admin dashboard).

- **Backend**: Go 1.22.4 with Echo web framework
- **Database**: SQLite (via GORM)
- **Frontend**: Vue 3 + TypeScript + Vite + Element Plus + OpenLayers + exifr
- **Authentication**: JWT (JSON Web Tokens)
- **API Documentation**: Swagger/OpenAPI

## Development Commands

### Backend Development

```bash
# Build for macOS arm64 and Linux amd64
make build

# Clean build artifacts
make clean

# Run backend in development mode
go run main.go serve -g -x http://127.0.0.1:8090

# Production-like setup
maphoto serve -p 8090 -d ./maphoto.db -a admin -w admin1234 -u /prefix -s false -g true -x http://127.0.0.1:8090
```

### Frontend Development

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies
npm install

# Run development server (with proxy to backend)
npm run dev

# Type checking
npm run type-check

# Build for production (outputs to internal/app/web/assets/dist)
npm run build

# Linting
npm run lint
```

### Testing

```bash
# Run backend tests
go test ./...

# Run specific test file
go test ./internal/util/uid_test.go

# Run frontend e2e tests
cd frontend && npm run test:e2e
```

### CLI Configuration Options

- `-port` / `-p`: Port to listen on (default: 8090)
- `-database-url` / `-d`: Database path (default: ./maphoto.db)
- `-admin-name` / `-a`: Default admin username (default: "admin")
- `-admin-password` / `-w`: Default admin password (default: "admin1234")
- `-url-prefix` / `-u`: URL prefix for routing
- `-disable-swagger` / `-s`: Disable Swagger UI
- `-debug` / `-g`: Enable debug mode
- `-domain` / `-x`: Domain for the application

### Frontend Environment Variables

- `VITE_API_BASE_URL`: API base URL (default: `/api/v1`)
- `VITE_DEFAULT_USERNAME`: Default username for anonymous viewing (default: `admin`)

## Architecture

### Directory Structure

```
frontend/                 # Vue 3 frontend application
├── src/
│   ├── api/             # API client (Axios)
│   ├── assets/          # Static assets and styles
│   ├── components/      # Vue components
│   │   ├── common/      # Common UI components
│   │   ├── map/         # Map components (OpenLayers)
│   │   ├── place/       # Place-related components
│   │   └── photo/       # Photo-related components
│   ├── composables/     # Vue composables
│   ├── layouts/         # Layout components
│   ├── router/          # Vue Router configuration
│   ├── stores/          # Pinia stores
│   ├── types/           # TypeScript type definitions
│   ├── utils/           # Utility functions
│   └── views/           # Page components
├── package.json
├── tsconfig.json
└── vite.config.ts

internal/
├── app/                 # Application logic
│   ├── form/           # Form structures and validation
│   ├── model/          # Data models (GORM entities)
│   ├── store/          # Database layer (SQLite with GORM)
│   └── web/            # Web layer
│       ├── api/        # REST API handlers
│       ├── jwt/        # JWT authentication middleware
│       ├── server/     # Echo server setup and routing
│       ├── view/       # View layer (HTML templates)
│       └── assets/     # Static assets (compiled from frontend)
├── cli/                # CLI command definitions (Cobra)
├── env/                # Environment configuration
└── util/               # Utility functions (logging, encryption, etc.)
```

### Key Components

#### Backend

1. **CLI Layer** (`internal/cli/`): Uses Cobra for command-line interface
   - `root.go`: Main command with global flags
   - `serve.go`: Web server command

2. **Web Server** (`internal/app/web/server/`): Echo framework setup
   - CORS configured to allow all origins (`*`) in development
   - JWT middleware for protected routes
   - Global middleware injects configuration into context

3. **Database Layer** (`internal/app/store/`): SQLite with GORM
   - Automatic initialization of admin user
   - Lightweight SQLite database file

4. **API Layer** (`internal/app/web/api/`): REST endpoints
   - Protected routes require JWT authentication
   - Swagger documentation available at `/swagger/*`

#### Frontend

5. **Vue 3 Application** (`frontend/`)
   - **Framework**: Vue 3 with Composition API and `<script setup>`
   - **Language**: TypeScript with strict mode
   - **Build Tool**: Vite
   - **UI Library**: Element Plus
   - **State Management**: Pinia
   - **Routing**: Vue Router
   - **Maps**: OpenLayers (ol package)

6. **Key Frontend Components**
   - `MapContainer.vue`: Core map component with OpenLayers
   - `PlaceCard.vue`: Place display card component
   - `PlaceForm.vue`: Form for creating/editing places
   - `PhotoUpload.vue`: Photo upload component
   - `AdminLayout.vue`: Admin dashboard layout

7. **Frontend Features**
   - **Anonymous Mode**: Public map viewing without login
   - **Management Mode**: Authenticated place/photo management
   - **Responsive Design**: Mobile and desktop support
   - **Group Filtering**: Filter places by group/category

### Configuration Flow

1. CLI parses flags and creates `env.Options`
2. Database initialized with `store.InitDB()`
3. Admin user created with `store.InitAdmin()`
4. Echo server started with `server.RunServer()`
5. Configuration injected into request context via middleware

### Frontend Build Process

1. Develop frontend in `frontend/` directory
2. Run `npm run build` to compile for production
3. Build output goes to `internal/app/web/assets/dist/`
4. Go embeds static files using `//go:embed` directives
5. Single binary deployment with embedded frontend

## Development Notes

### Database

- Uses SQLite for simplicity and easy deployment
- Database file created at runtime if it doesn't exist
- GORM handles migrations and schema management

### Authentication

- JWT tokens issued on login
- Protected routes use `echojwt` middleware
- Token secret generated at runtime with `util.ShortUID(12)`
- Token stored in localStorage and attached to API requests via Authorization header

### Frontend Architecture

- **Development**: Run `npm run dev` in `frontend/` directory, Vite dev server proxies API requests to backend
- **Production**: Build output embedded into Go binary via `//go:embed`
- **Routing**: Vue Router handles client-side routing with route guards for authentication
- **State**: Pinia stores manage global state (auth, places, user config)
- **Maps**: OpenLayers with OpenStreetMap tiles, supports vector markers, feature selection, and coordinate picking

### API Endpoints

#### Public Endpoints
- `POST /api/v1/user/login` - User login
- `GET /api/v1/share/:username` - Get public places for a user (supports `?group=` filter)

#### Protected Endpoints (require JWT)
- `GET /api/v1/place/all` - Get all places for current user
- `POST /api/v1/place/add` - Create new place
- `POST /api/v1/place/update/:id` - Update place
- `GET /api/v1/place/del/:id` - Delete place
- `POST /api/v1/place/cover` - Update place cover image
- `POST /api/v1/place/pic/add` - Add photo to place
- `POST /api/v1/place/pic/del` - Remove photo from place
- `GET /api/v1/user/config` - Get user config
- `POST /api/v1/user/config/update` - Update user config

### API Documentation

- Swagger UI available at `/swagger/index.html` (unless disabled with `-s`)
- API documentation generated from Go comments
- See `docs/` directory for OpenAPI specifications

### Deployment Considerations

- When using Nginx reverse proxy, set `-url-prefix` to match location block
- Application supports running behind subpaths (e.g., `/maphoto`)
- Single binary deployment with embedded assets
- Frontend build must be run before `go build` to include latest UI changes

## 近期变更记录 (2026-04-04)

### 安全改进
1. **地点管理权限控制强化**：确保用户只能管理自己创建的地点
   - 改进 `FindPlace` 函数，改为直接数据库查询并检查 `user_id` 字段
   - 改进 `PlaceDel` 函数，直接删除地点记录并级联删除照片
   - 所有地点操作API（添加、更新、删除、修改封面、添加/删除照片）均已确保严格的用户权限检查

## 近期变更记录 (2026-02-19)

### 功能改进
1. **EXIF GPS 坐标读取**：创建/修改地点时支持从本地照片读取 GPS 坐标
   - 新增 `exifr` 依赖（支持 JPEG, TIFF, HEIC/HEIF 格式）
   - `PlaceForm.vue` 坐标区域添加"从照片读取"按钮
   - 自动解析照片 EXIF 信息并填充经纬度坐标
   - 地图自动定位到读取的坐标位置

## 近期变更记录 (2026-02-15)

### 功能改进
1. **照片管理优化**：修改地点绑定照片从文件上传改为 URL 输入，与后端接口对齐 (`internal/app/form/place_attach.go`)
2. **创建地点支持批量添加照片**：`PlaceForm.vue` 支持输入多个照片 URL，通过 `photos` 数组提交
3. **地点列表排序**：`PlaceAll` 接口返回结果按创建时间倒序排列（新的在前）
4. **支持修改地点分组**：`PlaceUpdateForm` 前后端均添加 `group` 字段

### Bug 修复
1. **JWT Token 修复**：开发环境 Token 和用户登录 Token 添加 `id` 字段，修复 `ParseJWT` 解析失败
2. **后端接口修复**：
   - `PlacePicAdd`：修复 photo URL 为空的问题
   - `PlaceUpdate`：修复查询语句错误（使用 `.First()` 正确查询），列名修正为 `user_id`
   - 添加/删除照片、更新地点后重新查询返回完整数据

### 构建优化
1. **Makefile 改进**：`make build` 先执行 `frontend-build`（npm install + build），再编译 Go 后端

### 前端优化
1. **路由缓存**：`App.vue` 添加 `<keep-alive>`，首页、分享页、管理页等标记 `meta.keepAlive: true`
2. **组件更新**：`PhotoUpload.vue` 改为 URL 输入组件，`PlaceForm.vue` 支持照片列表管理

---

## 自我维护规则（Self-Maintenance Rule）

- 每次完成重大变更后（新模块、新页面、架构调整、引入新库/技术、目录结构变化、重要约定修改等），**必须**把"最后一步"设为：
  "检查 CLAUDE.md 是否需要同步更新，并输出更新建议或直接修改它"
- 重大变更的定义包括但不限于：
  - 新增/重构主要目录
  - 引入新框架/状态管理/构建工具
  - 重大路由/组件树变化
  - 规范/风格/linter 规则变更
- 保持本文件始终是项目的"活文档"和"单点真相"
