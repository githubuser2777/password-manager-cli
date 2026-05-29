\# Ranh giới Bảo mật (Security Guidelines)



\## 1. Quản lý Bộ nhớ (Memory Zeroing)

\- Tuyệt đối không lưu trữ Master Password dưới dạng chuỗi `string` cố định lâu trong bộ nhớ.

\- Sử dụng mảng byte `\[]byte` và ép buộc ghi đè bằng giá trị `0` ngay sau khi hàm mã hóa/giải mã kết thúc.



\## 2. Thư viện Nghiêm cấm (Forbidden Packages)

\- Không sử dụng `math/rand` cho các tác vụ liên quan đến mật mã hoặc sinh chuỗi ngẫu nhiên. Bắt buộc dùng `crypto/rand`.

\- Không sử dụng các hàm băm yếu như MD5, SHA1 cho việc xử lý mật khẩu.



\## 3. Quyền hạn Hệ thống

\- File cơ sở dữ liệu cục bộ bắt buộc phải được giới hạn quyền truy cập chỉ cho chủ sở hữu (`chmod 600`). Agent khi tạo file phải kiểm tra điều kiện này.

