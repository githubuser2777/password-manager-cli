# Password Manager CLI

Một công cụ quản lý mật khẩu dòng lệnh (CLI) được viết bằng Go, tập trung vào bảo mật và tính tiện dụng.

## Tính Năng Nổi Bật
- **Mã hóa AES-256-GCM**: Mã hóa xác thực, tiêu chuẩn cao nhất hiện nay.
- **Sinh khóa Argon2id**: Ngăn chặn tối đa các cuộc tấn công brute-force vào Master Password.
- **Lưu trữ Cục bộ (Local Storage)**: Dữ liệu của bạn 100% nằm trên máy của bạn (mặc định tại `~/.passmgr/vault.enc`).
- **Clipboard Integration**: Hỗ trợ copy mật khẩu vào bộ nhớ đệm mà không cần in ra màn hình.
- **Sinh mật khẩu ngẫu nhiên**: Có sẵn tính năng tạo mật khẩu siêu mạnh.

## Cài Đặt

Yêu cầu đã cài đặt [Go](https://go.dev/dl/).
Mở terminal và chạy lệnh sau để build mã nguồn:
```bash
go build -o passmgr.exe
```
*(Trên Linux/Mac thì dùng `go build -o passmgr`)*

## Hướng Dẫn Sử Dụng

### 1. Khởi tạo kho lưu trữ (Vault)
```bash
passmgr init
```
Lệnh này sẽ yêu cầu bạn tạo Master Password. Master Password là chìa khóa duy nhất để mở vault, nếu quên bạn sẽ mất toàn bộ dữ liệu.

### 2. Thêm mật khẩu mới
```bash
passmgr add github.com
```
Lệnh này sẽ hỏi bạn Username và Password. Bạn cũng có thể để công cụ tự sinh password ngẫu nhiên bằng cờ `-g`:
```bash
passmgr add github.com -g
```

### 3. Lấy mật khẩu
```bash
passmgr get github.com
```
Để copy trực tiếp vào clipboard (bảo mật hơn vì không in ra màn hình):
```bash
passmgr get github.com -c
```

### 4. Liệt kê các tài khoản đã lưu
```bash
passmgr list
```

### 5. Tiện ích sinh mật khẩu ngẫu nhiên
```bash
passmgr generate 20
```
(Tạo mật khẩu dài 20 ký tự, có thể bỏ ký tự đặc biệt bằng cờ `-n`).

## Build Scripts
Bạn có thể sử dụng các script có sẵn để build và test tự động:
- **Windows**: Chạy file `.\build.ps1`
- **Linux/Mac**: Chạy lệnh `make build` (yêu cầu cài đặt Make)
