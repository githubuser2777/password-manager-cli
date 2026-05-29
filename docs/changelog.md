# Changelog

Tất cả các thay đổi lớn của dự án sẽ được ghi nhận tại file này.

## [Phase 5] - 2026-05-29
### Added
- Thêm tệp `LICENSE` (MIT License) để công khai dự án dạng mã nguồn mở.
- Di chuyển và viết lại `SECURITY.md` thành `.github/SECURITY.md` theo chuẩn Security Policy của GitHub.
- Đã thêm `README.md` với đầy đủ hướng dẫn sử dụng công cụ.
- Tạo script build dành cho Windows (`build.ps1`) và Linux/Mac (`Makefile`).
- Hoàn tất dự án phiên bản 1.0.



## [Phase 4] - 2026-05-29
### Added
- Tích hợp package `golang.org/x/term` để nhập Master Password bảo mật, không hiển thị trên màn hình.
- Tích hợp package `github.com/atotto/clipboard` để hỗ trợ copy password nhanh.
- Hoàn thiện 5 bộ lệnh (commands) chính thông qua Cobra:
  - `passmgr init`: Khởi tạo kho chứa mới với mật khẩu chủ.
  - `passmgr add <service>`: Thêm tài khoản mới, hỗ trợ cờ `-g` sinh password ngẫu nhiên.
  - `passmgr get <service>`: Lấy password, hỗ trợ cờ `-c` sao chép trực tiếp vào Clipboard.
  - `passmgr list`: Xem danh sách tất cả các service đã lưu.
  - `passmgr generate`: Tiện ích tạo mật khẩu nhanh ở màn hình ngoài.


## [Phase 3] - 2026-05-29
### Added
- Thêm module `internal/storage/storage.go` quản lý việc IO đọc/ghi file `vault.enc`.
- Xây dựng cơ chế Load/Save an toàn (Atomic Write): Ghi ra file tạm trước khi rename để tránh lỗi hỏng file nếu bị crash giữa chừng.
- Đóng gói dữ liệu nhị phân với cấu trúc `Salt(16) + Nonce(12) + Ciphertext`.
- Viết Unit Tests (`storage_test.go`) giả lập luồng ghi/đọc bằng thư mục tạm thời.


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
