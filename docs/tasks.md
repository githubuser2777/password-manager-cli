# Execution Checklist

Quy định: 
1. Khi hoàn thành một task, Code Agent sẽ đánh dấu `[x]` tương ứng. Nếu một task quá phức tạp, chủ động chia nhỏ thêm trong file này.
2. Tuân thủ tuyệt đối `rules.md`: **Không mở/đọc quá 5 file cùng lúc**; **Bắt buộc viết comment** cho các logic mã hóa phức tạp.

## Phase 1: Foundation
- [x] Khởi tạo thư mục và Go module (`go mod init github.com/yourusername/password-manager-cli` hoặc tương tự).
- [x] Thiết lập cấu trúc thư mục: `cmd/`, `internal/crypto/`, `internal/storage/`, `internal/core/`.
- [x] Cài đặt package `spf13/cobra` và thiết lập file `main.go`, `cmd/root.go`.

## Phase 2: Domain Models & Cryptography Module
- [x] (internal/core): Định nghĩa struct `Vault` và `Entry`.
- [x] (internal/crypto): Cài đặt hàm `DeriveKey(masterPassword string, salt []byte) []byte` sử dụng Argon2id. Cài đặt hàm `GenerateSalt() []byte`.
- [x] (internal/crypto): Cài đặt hàm `Encrypt(plaintext, key, nonce []byte) ([]byte, error)` và `Decrypt(ciphertext, key, nonce []byte) ([]byte, error)` sử dụng AES-256-GCM.
- [x] (internal/crypto): Cài đặt hàm `GenerateRandomPassword(length int, includeSpecial bool) string`.
- [x] Viết Unit Tests kiểm tra tính chính xác của `crypto` module (Mã hóa xong giải mã phải ra dữ liệu gốc).

## Phase 3: Storage Module
- [x] (internal/storage): Xây dựng hàm `SaveVault(path string, masterPassword string, vault *core.Vault) error` dựa theo Storage Flow (ghi binary Salt + Nonce + Ciphertext).
- [x] (internal/storage): Xây dựng hàm `LoadVault(path string, masterPassword string) (*core.Vault, error)`.
- [x] Viết Unit Tests cho việc lưu và tải vault từ file (dùng thư mục tạm).

## Phase 4: CLI Commands
- [x] (cmd/init.go): Implement lệnh `init`. (Kiểm tra file vault đã tồn tại chưa, hỏi người dùng nhập Master Password ẩn 2 lần để xác nhận, tạo vault rỗng với Salt mới và lưu lại).
- [x] (cmd/add.go): Implement lệnh `add <service>`. (Nhập username, nhập password hoặc tự động tạo nếu dùng flag `--generate`).
- [x] (cmd/get.go): Implement lệnh `get <service>`. (Hiển thị username và password. Hỗ trợ flag `--copy` để tự động copy vào clipboard thay vì in ra màn hình).
- [x] (cmd/list.go): Implement lệnh `list`. (In ra tất cả các service đã lưu).
- [x] (cmd/generate.go): Implement lệnh `generate`. (Sinh mật khẩu ngẫu nhiên để dùng ngay).

## Phase 5: Polish & Build
- [ ] Đóng gói và viết script build (`Makefile` hoặc `build.sh`).
- [ ] Viết `README.md` hướng dẫn sử dụng.
