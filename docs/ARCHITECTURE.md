\# Kiến trúc Hệ thống (System Architecture)



\## 1. Luồng xử lý Mật mã học (Cryptography Pipeline)



Hệ thống sử dụng cơ chế mã hóa đối xứng để bảo vệ Vault. Dữ liệu không bao giờ được lưu dưới dạng plain-text.



Master Password ──> \[ Argon2id + Salt ] ──> Master Key (32 bytes)

&#x20;                                                  │

Plaintext Data  ──> \[ AES-256-GCM + Nonce ] <──────┘

&#x20;                                                  │

&#x20;                                                  ▼

&#x20;                                        \[ Encrypted Vault File ]



\## 2. Lưu trữ Dữ liệu (Storage Mechanism)



\- Vị trí mặc định: `\~/.config/pwdmgr/vault.json` (đối với hệ điều hành Linux/macOS).

\- Định dạng cấu trúc file Vault:



```json

{

&#x20; "salt": "chuỗi\_base64\_ngẫu\_nhiên",

&#x20; "entries": {

&#x20;   "tên\_dịch\_vụ": "dữ\_liệu\_mã\_hóa\_gcm\_base64"

&#x20; }

}

```

\## 3. Cấu trúc Lệnh CLI (Cobra Command Tree)



root (pwdmgr)

├── init (Khởi tạo kho chứa)

├── add  (Thêm tài khoản mới)

├── get  (Lấy và giải mã)

└── gen  (Sinh mật khẩu ngẫu nhiên)

