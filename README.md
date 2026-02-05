# Drive Course Viewer

Xem video từ Google Drive (Share with me) dưới dạng khóa học như Udemy hoặc YouTube playlist.

## Cấu hình

### 1. Google Cloud Console

1. Tạo project tại [Google Cloud Console](https://console.cloud.google.com/)
2. Bật **Google Drive API**
3. Tạo **OAuth 2.0 Client ID** (ứng dụng Web)
4. Thêm Authorized redirect URI: `http://localhost:8080/api/auth/callback`
5. Lấy **Client ID** và **Client Secret**

### 2. Biến môi trường

**Backend** (`backend/`):

```
GOOGLE_CLIENT_ID=your-client-id
GOOGLE_CLIENT_SECRET=your-client-secret
BACKEND_URL=http://localhost:8080
FRONTEND_URL=http://localhost:3000
```

**Frontend** (tạo file `frontend/.env.local`):

```
NEXT_PUBLIC_BACKEND_URL=http://localhost:8080
```

### 3. Chạy ứng dụng

**Terminal 1 - Backend:**
```bash
cd backend
go run ./cmd/server
```

**Terminal 2 - Frontend:**
```bash
cd frontend
npm install
npm run dev
```

Mở http://localhost:3000

## Cách sử dụng

1. Đăng nhập với Google
2. Nhấn "Tải danh sách" để xem các thư mục Shared with me
3. Chọn các thư mục chứa video
4. Nhấn "Tạo khóa học"
5. Chọn bài học và xem video
