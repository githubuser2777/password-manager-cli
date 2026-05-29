# Changelog

Tất cả các thay đổi lớn của dự án sẽ được ghi nhận tại file này.

## [Phase 2] - 2026-05-29
### Added
- Thêm `internal/core/models.go` định nghĩa cấu trúc `Vault` và `Entry`.
- Viết module mã hóa `internal/crypto/crypto.go` với Argon2id (Key Derivation) và AES-256-GCM (Encrypt/Decrypt).
- Thêm hàm `GenerateRandomPassword` hỗ trợ sinh mật khẩu có tùy chọn ký tự đặc biệt.
- Viết Unit Tests (`crypto_test.go`) cho toàn bộ module mã hóa (Coverage 100% Core Logic).


## [Phase 1] - 2026-05-29
### Added
- Khởi tạo thư mục và Go module (`password-manager-cli`).
- Thiết lập cấu trúc thư mục chuẩn: `cmd/`, `internal/crypto/`, `internal/storage/`, `internal/core/`.
- Cài đặt package `spf13/cobra`.
- Thiết lập file `main.go` và lệnh gốc `cmd/root.go` với các mô tả cơ bản.
- Tích hợp thành công GitHub Spec Kit (`.specify/`, `.github/`).

### Changed
- Cập nhật `rules.md` thêm phần quy định về quy trình Git và Changelog.
