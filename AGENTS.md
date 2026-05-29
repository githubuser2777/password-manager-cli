# Hệ thống Agents cho Antigravity (Spec-Driven Development)

Tài liệu này định hướng cách các Agents trong hệ thống **Antigravity** tương tác, phối hợp và thực thi công việc dựa trên triết lý **Spec-Driven Development (SDD)** (tương tự GSD / GitHub Spec Kit). 

Mục tiêu cốt lõi: Ngăn chặn tình trạng "Context Rot" (thất thoát ngữ cảnh khi chat quá dài) và "Vibe Coding" (code cảm tính không có kế hoạch).

---

## 1. Cơ chế Định hướng Agent (Agent Routing)
Các agents không được gắn với các persona sáo rỗng. Phân quyền của chúng dựa trên **Phase (Giai đoạn)** của dự án và **Context (Ngữ cảnh)** hiện tại.

### A. Tác tử Đặc tả (Spec Agent)
- **Hành vi:** Thu thập dữ liệu đầu vào từ user (ý tưởng, file cũ, báo cáo bug) và đặt câu hỏi để chốt chặt phạm vi (scope).
- **Trách nhiệm:** Tập trung 100% vào "What" (Chúng ta làm cái gì?) và "Why" (Tại sao phải làm?). KHÔNG viết code hay thiết kế kiến trúc.
- **Đầu ra (Output):** Tạo/Cập nhật file `docs/spec.md`.

### B. Tác tử Kiến trúc (Planning Agent)
- **Hành vi:** Tiêu thụ file `docs/spec.md` để phác thảo thiết kế kỹ thuật (System Blueprint).
- **Trách nhiệm:** Đưa ra quyết định "How" (Làm như thế nào?). Xác định stack, cấu trúc thư mục, API schemas, Data models, và thư viện cần dùng.
- **Đầu ra (Output):** Tạo/Cập nhật file `docs/plan.md`.

### C. Tác tử Điều phối (Task Agent)
- **Hành vi:** Chuyển đổi `docs/plan.md` thành một danh sách công việc (checklist) có thể thực thi độc lập.
- **Trách nhiệm:** Chia nhỏ khối lượng công việc. Nếu một task mất quá nhiều token/bước để xử lý, nó phải được chia nhỏ hơn nữa.
- **Đầu ra (Output):** Tạo/Cập nhật file `docs/tasks.md`.

### D. Tác tử Lập trình (Code Agent / Editor)
- **Hành vi:** Đọc file `docs/tasks.md` và thực thi từng task một cách tuần tự.
- **Trách nhiệm:** Mở file, viết code, sửa lỗi, chạy test. Sau khi xong một task, tự động check `[x]` vào `docs/tasks.md` trước khi sang task mới. Tránh đọc toàn bộ source code nếu không cần thiết.

---

## 2. Giao thức Tương tác (Communication Protocol)
- **Chuyển giao trạng thái (Handoff):** Khi một Agent kết thúc Phase của mình, nó phải gọi rõ output file để Agent tiếp theo tiếp quản (VD: "Spec hoàn tất tại `docs/spec.md`. Sẵn sàng cho Phase Planning").
- **Kích hoạt công cụ (Tools):** Khuyến khích sử dụng Command Line Tools (bash/zsh/go/npx) để tự động hóa các thao tác lint, test, hoặc compile ngay trong terminal của Antigravity.
