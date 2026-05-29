# Execution Checklist

Quy định: 
1. Khi hoàn thành một task, Code Agent sẽ đánh dấu `[x]` tương ứng. Nếu một task quá phức tạp, chủ động chia nhỏ thêm trong file này.
2. Tuân thủ tuyệt đối `rules.md`: **Không mở/đọc quá 5 file cùng lúc**; **Bắt buộc viết comment** cho các logic mã hóa phức tạp.

## Phase 1: Foundation
- [ ] Khởi tạo thư mục và Go module (`go mod init github.com/yourusername/password-manager-cli` hoặc tương tự).
- [ ] Thiết lập cấu trúc thư mục: `cmd/`, `internal/crypto/`, `internal/storage/`, `internal/core/`.
- [ ] Cài đặt package `spf13/cobra` và thiết lập file `main.go`, `cmd/root.go`.

## Phase 2: Domain Models & Cryptography Module
- [ ] (internal/core): Định nghĩa struct `Vault` và `Entry`.
- [ ] (internal/crypto): Cài đặt hàm `DeriveKey(masterPassword string, salt []byte) []byte` sử dụng Argon2id. Cài đặt hàm `GenerateSalt() []byte`.
- [ ] (internal/crypto): Cài đặt hàm `Encrypt(plaintext, key, nonce []byte) ([]byte, error)` và `Decrypt(ciphertext, key, nonce []byte) ([]byte, error)` sử dụng AES-256-GCM.
- [ ] (internal/crypto): Cài đặt hàm `GenerateRandomPassword(length int, includeSpecial bool) string`.
- [ ] Viết Unit Tests kiểm tra tính chính xác của `crypto` module (Mã hóa xong giải mã phải ra dữ liệu gốc).

## Phase 3: Storage Module
- [ ] (internal/storage): Xây dựng hàm `SaveVault(path string, masterPassword string, vault *core.Vault) error` dựa theo Storage Flow (ghi binary Salt + Nonce + Ciphertext).
- [ ] (internal/storage): Xây dựng hàm `LoadVault(path string, masterPassword string) (*core.Vault, error)`.
- [ ] Viết Unit Tests cho việc lưu và tải vault từ file (dùng thư mục tạm).

## Phase 4: CLI Commands
- [ ] (cmd/init.go): Implement lệnh `init`. (Kiểm tra file vault đã tồn tại chưa, hỏi người dùng nhập Master Password ẩn 2 lần để xác nhận, tạo vault rỗng với Salt mới và lưu lại).
- [ ] (cmd/add.go): Implement lệnh `add <service>`. (Yêu cầu Master Password, load vault, hỏi username/password, lưu lại. Thêm cờ `--generate` để tự sinh password).
- [ ] (cmd/get.go): Implement lệnh `get <service>`. (Load vault, tìm service, in ra password hoặc sử dụng clipboard để copy).
- [ ] (cmd/list.go): Implement lệnh `list`. (Load vault, in danh sách các dịch vụ hiện có trong vault).
- [ ] (cmd/generate.go): Implement lệnh `generate`. (Chỉ sinh password và in ra màn hình, không cần vault).

## Phase 5: Polish & Build
- [ ] Đóng gói và viết script build (`Makefile` hoặc `build.sh`).
- [ ] Viết `README.md` hướng dẫn sử dụng.
