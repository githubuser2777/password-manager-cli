# Technical Blueprint: Password Manager CLI

## 1. Tech Stack
- **Ngôn ngữ**: Go (Golang)
- **CLI Framework**: `github.com/spf13/cobra` cho cấu trúc lệnh và `golang.org/x/term` để ẩn password khi nhập.
- **Cryptography**: 
  - Khởi tạo khóa: `golang.org/x/crypto/argon2`
  - Mã hóa/Giải mã: `crypto/aes` và `crypto/cipher` (thư viện chuẩn Go).
  - Sinh số ngẫu nhiên an toàn: `crypto/rand`
- **Lưu trữ**: File cục bộ, nội dung là cấu trúc JSON đã được mã hóa toàn bộ bằng AES-256.
- **Tiện ích (Optional)**: `github.com/atotto/clipboard` để copy vào bộ nhớ đệm.

## 2. Architecture & Modules
Dự án áp dụng cấu trúc thư mục chuẩn Go (`Standard Go Project Layout`):
- `cmd/`: Chứa định nghĩa và entrypoint của các lệnh CLI (`root.go`, `init.go`, `add.go`, `get.go`, `list.go`, `generate.go`).
- `internal/crypto/`: Đóng gói toàn bộ logic mật mã (AES-GCM, Argon2id, Random Password).
- `internal/storage/`: Quản lý IO (đọc/ghi file cục bộ) và thao tác với JSON.
- `internal/core/`: Chứa các domain models (Entity) cốt lõi của ứng dụng.

## 3. Data Models

### Entry (Một tài khoản)
```go
type Entry struct {
    Username  string `json:"username"`
    Password  string `json:"password"`
    CreatedAt string `json:"created_at"` // ISO8601 string
}
```

### Vault (Kho chứa tổng)
```go
type Vault struct {
    // Salt dùng chung cho việc băm Argon2 (cần thiết để sinh lại Key). Sẽ được sinh ngẫu nhiên khi init.
    Salt    []byte           `json:"salt"` 
    Entries map[string]Entry `json:"entries"` // Key là tên service (vd: "github", "google")
}
```

## 4. Storage Flow
- **Cấu trúc File `vault.enc`**:
  Do `Salt` cần được biết trước để băm Master Password, nên ta có hai hướng:
  - Hướng 1: `vault.enc` chứa phần `Salt` ở đầu file không mã hóa (plaintext), tiếp theo là `Nonce` và Dữ liệu đã mã hóa (Ciphertext).
  - Hướng 2: Lưu metadata (Salt, Nonce) và Ciphertext dưới dạng một JSON bọc bên ngoài.
  => **Lựa chọn Hướng 1** (Binary format) để tối ưu và bảo mật (Salt (16 bytes) + Nonce (12 bytes) + Ciphertext).

- **Đọc (Load)**:
  1. Đọc file `vault.enc`.
  2. Tách 16 bytes đầu làm `Salt`, 12 bytes tiếp theo làm `Nonce`, phần còn lại là `Ciphertext`.
  3. Dùng `Salt` + Master Password truyền qua Argon2id -> `Key 256-bit`.
  4. Giải mã `Ciphertext` với `Key` và `Nonce` bằng AES-GCM -> JSON string.
  5. JSON Unmarshal thành struct `Vault`.

- **Ghi (Save)**:
  1. JSON Marshal struct `Vault`.
  2. Sinh `Nonce` mới (12 bytes).
  3. Mã hóa JSON string bằng AES-GCM -> `Ciphertext`.
  4. Ghi file: `Salt` + `Nonce` + `Ciphertext`.

## 5. Coding Standards & Conventions (Theo rules.md)
- **Quản lý Context**: Code Agent không được đọc quá 5 files cùng lúc.
- **Tiêu chuẩn Mã nguồn**: Viết code rõ ràng, ưu tiên dễ hiểu hơn là code quá ngắn.
- **Bình luận (Comments)**: Bắt buộc phải có comment giải thích chi tiết cho các logic nghiệp vụ phức tạp, đặc biệt là trong module `internal/crypto` (xử lý Argon2id, AES-GCM).
