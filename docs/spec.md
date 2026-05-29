# Product Specification: Password Manager CLI

## 1. Overview
Công cụ dòng lệnh (CLI) giúp tạo mật khẩu ngẫu nhiên có độ bảo mật cao và lưu trữ chúng một cách an toàn bằng ngôn ngữ Go. Toàn bộ kho mật khẩu được bảo vệ và giải mã thông qua một Master Password duy nhất.

## 2. Target Audience
- Cá nhân, quản trị viên hệ thống, lập trình viên cần một trình quản lý mật khẩu nhanh chóng, an toàn hoạt động hoàn toàn trên terminal.

## 3. Core Features
- `init`: Khởi tạo kho lưu trữ (vault) ban đầu, thiết lập Master Password.
- `add`: Thêm một tài khoản (Website/App, Username, Password). Tùy chọn tự động sinh mật khẩu cho tài khoản này.
- `get`: Lấy và giải mã mật khẩu dựa trên tên (yêu cầu xác thực Master Password). Hỗ trợ copy thẳng vào clipboard.
- `list`: Xem danh sách các dịch vụ đã lưu (không hiển thị mật khẩu thực tế).
- `generate`: Sinh mật khẩu ngẫu nhiên an toàn, cho phép tùy chỉnh độ dài và loại ký tự.

## 4. Security Requirements
- **Encryption**: Dữ liệu phải được mã hóa bằng AES-256 ở chế độ GCM (đảm bảo cả tính bảo mật và tính toàn vẹn của dữ liệu).
- **Key Derivation**: Sử dụng thuật toán `Argon2id` để băm Master Password, chống lại các cuộc tấn công Brute-force và Dictionary attack.
- **Storage**: Mật khẩu được mã hóa nguyên khối và lưu cục bộ tại máy tính của người dùng (vd: `~/.passmgr/vault.enc`).
- **Memory Security**: Khi nhập Master Password trên terminal, ký tự không được phép hiển thị ra màn hình.
