\# Bản đồ Ngữ cảnh (Source Code Mapping)



\## 1. Định vị Module

\- `internal/crypto/`: Chứa logic băm Argon2id và mã hóa AES-GCM. Nếu cần sửa thuật toán mã hóa, chỉ can thiệp vùng này.

\- `internal/vault/`: Chứa các hàm `Read()`, `Write()`, `Modify()` tương tác với hệ thống tệp.

\- `cmd/`: Chứa định nghĩa giao diện CLI của Cobra.



\## 2. Quy tắc Đọc/Ghi dành cho Agent

\- Khi có yêu cầu thay đổi tính năng: Cập nhật `docs/FEATURES.md` trước, sau đó mới tiến hành sửa code ở thư mục `cmd/`.

\- Khi cần tối ưu hóa hiệu năng hoặc bảo mật: Đối chiếu với `docs/SECURITY.md` để đảm bảo không vi phạm các nguyên tắc cốt lõi.

