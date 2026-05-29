\# Đặc tả Tính năng (Features Specification)



\## 1. Lệnh `init`

\- \*\*Hành vi:\*\* Yêu cầu người dùng nhập Master Password (ẩn ký tự khi nhập).

\- \*\*Xử lý:\*\* Sinh Salt ngẫu nhiên, khởi tạo file cấu trúc JSON trống, thiết lập quyền file `0600`.



\## 2. Lệnh `add <service\\\_name>`

\- \*\*Hành vi:\*\* Nhận tên dịch vụ từ đối số. Yêu cầu nhập tài khoản và mật khẩu cần lưu.

\- \*\*Xử lý:\*\* Kiểm tra file vault đã tồn tại chưa. Nếu có, giải mã, chèn thêm entry mới, mã hóa lại toàn bộ và ghi đè file.



\## 3. Lệnh `get <service\\\_name>`

\- \*\*Hành vi:\*\* Yêu cầu Master Password, giải mã entry tương ứng và in ra màn hình.

\- \*\*Tùy chọn:\*\* Hỗ trợ flag `--clip` để copy thẳng vào clipboard thay vì in ra terminal (tự động xóa clipboard sau 10 giây).



\## 4. Lệnh `gen`

\- \*\*Hành vi:\*\* Sinh chuỗi ngẫu nhiên dựa trên `crypto/rand`.

\- \*\*Cấu hình mặc định:\*\* Độ dài 16 ký tự, bao gồm \[A-Z], \[a-z], \[0-9], và ký tự đặc biệt.



