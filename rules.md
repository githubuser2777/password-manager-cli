# Quy tắc Cốt lõi (Core Rules)

Bộ quy tắc này được áp dụng như "hiến pháp" (Constitution) cho mọi Agents hoạt động trong môi trường Antigravity. 

## 1. Hành vi Trò chuyện & Định hình (No Fluff)
- **KHÔNG SỬ DỤNG MẪU CÂU CLICHÉ:** Nghiêm cấm mở đầu prompt hoặc câu trả lời bằng "You are an expert software engineer...", "I am a world-class AI...", "As an AI...". Bắt tay ngay vào giải quyết vấn đề.
- **Ngắn gọn & Súc tích:** Trả lời trực diện. Nếu user yêu cầu viết code, chỉ output logic, code và các thay đổi. Hạn chế xin lỗi hay giải thích dông dài.
- **Linh hoạt (Flexibility):** Luật sinh ra để định hướng, không phải để làm chậm tiến độ. Nếu một thay đổi là cực kỳ nhỏ (ví dụ: sửa một typo, đổi 1 dòng CSS), cho phép bỏ qua quy trình Spec/Plan và thực hiện "Fast Track".

## 2. Xử lý Lỗi & Giả định (Anti-Hallucination)
- Nếu user đưa ra một yêu cầu mơ hồ có tác động lớn tới kiến trúc, **Agent PHẢI dừng lại và đặt câu hỏi**.
- Không tự ý sinh ra các API giả định (mock APIs) hoặc thư viện không tồn tại. Dựa trên các specs trong `docs/plan.md`.
- Nếu có lỗi xảy ra trong quá trình chạy (Error Logs), luôn đọc log một cách kỹ lưỡng, chỉ ra điểm bất thường trước khi đề xuất cách sửa. Tránh spam code sửa lỗi mù quáng.

## 3. Quản lý File & Ngữ cảnh (Context Limits)
- Để tránh tràn bộ nhớ ngữ cảnh (Context Rot), không mở / read quá 5 files một lúc trừ khi thực hiện thao tác tìm kiếm (grep/search).
- Viết code có tính module hóa cao. 

## 4. Tiêu chuẩn Mã nguồn
- Tính rõ ràng ưu tiên hơn độ ngắn gọn. Dùng tên biến mô tả đúng chức năng.
- Tuân thủ DRY (Don't Repeat Yourself) một cách hợp lý; đừng ép DRY nếu nó khiến code bị trừu tượng hóa quá đà và khó đọc.
- Để lại bình luận (comments) ở các hàm có logic nghiệp vụ phức tạp, thuật toán rẽ nhánh, hoặc những đoạn code là "workaround" (cách chữa cháy).

## 5. Quy trình Git & Changelog
- **Git Commit & Push**: Agent phải tự động commit và push sau khi hoàn thành các thay đổi đáng kể (VD: xong 1 phase).
- **Review trước khi Commit**: KHÔNG ĐƯỢC dùng `git add .` một cách mù quáng. Phải kiểm tra file nào thay đổi (`git status`) và chỉ add những file thực sự liên quan đến commit đó.
- **Changelog**: Tự động tạo và cập nhật file `docs/changelog.md` cho mỗi lần có thay đổi lớn. Mỗi thay đổi phải được phân tách rõ ràng (VD: sử dụng Header ngày tháng, version) nhưng lưu chung trong cùng một file.
