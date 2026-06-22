# tverp-utils-sandbox

Bộ bài tập thử việc cho vị trí **Go Developer**.

Đây là một module Go độc lập (`go 1.24`), gồm các hàm tiện ích thuần (pure functions):
xử lý ngày/giờ theo múi giờ Việt Nam, sinh mã, chuẩn hóa chuỗi tiếng Việt, ước lượng
thời lượng đọc, và lược bỏ markdown. **Không cần kiến thức nghiệp vụ** — chỉ cần Go vững.

## Bắt đầu

```bash
go version          # cần Go >= 1.24
go mod download     # tải dependency (testify, golang.org/x/text)
make build          # phải build sạch
make test           # chạy toàn bộ test hiện có
```

Quy trình nộp bài: tạo branch riêng cho mỗi task → commit → mở Pull Request →
chờ review. **Mỗi task một branch + một PR.**

## Quy tắc chung

- Mỗi task có file/thư mục riêng, độc lập nhau — làm theo thứ tự khuyến nghị bên dưới.
- `go vet ./...` phải sạch; `gofmt`/`goimports` đã format.
- Code phải đọc được như code xung quanh: đặt tên, comment, idiom đồng nhất.
- Khi sửa code có sẵn: **không được làm hỏng test đang xanh**.
- Được dùng AI hỗ trợ, nhưng bạn phải **hiểu và bảo vệ được** mọi dòng mình nộp
  (sẽ có buổi trao đổi trực tiếp: giải thích quyết định + mở rộng code tại chỗ).

## Thứ tự khuyến nghị

T1 → T2 → T4 → T5 → T6 → T7 → T3

---

## T1 — Viết unit test cho các hàm xử lý ngày (`date/`)

Thư mục `date/` có **8 biến thể** của cùng một bộ hàm tiện ích ngày/giờ VN
(`ad`, `channel`, `config`, `editorial`, `nrcs`, `production`, `resource`, `royalty`).
Chỉ `date/channel/` đã có sẵn file test (`date_test.go`) — dùng làm **mẫu tham khảo**.

**Việc cần làm:** viết test (table-driven) cho **7 biến thể còn lại** (tất cả trừ `channel`).

Yêu cầu:
- Bao phủ mọi hàm export trong từng file (`VNLocation`, `ParseDateInVN`,
  `NormalizeDateToVN`, `TodayInVN`, ... — mỗi biến thể có thể khác nhau đôi chút).
- Cân nhắc edge case: năm nhuận (29/02), biên ngày, định dạng sai, chuỗi rỗng,
  chuyển đổi múi giờ (UTC ↔ VN, lệch ngày khi qua nửa đêm).
- Đặt cùng package (`package util`) với file nguồn, đặt tên `date_test.go`.

**Đánh giá:** nền tảng Go, idiom table-driven, tư duy edge-case, tính tỉ mỉ.

---

## T2 — Hợp nhất 8 biến thể `date.go` thành một module dùng chung

8 file trong `date/` gần như trùng lặp nhưng **không giống hệt** — một số biến thể có
thêm hàm (vd `channel` có `IsDateInPast`/`SameDate`, `nrcs` có range nửa mở, `config`
tối giản hơn). Đây là nợ kỹ thuật cần khử.

**Việc cần làm:** thiết kế **một** package dùng chung (vd `date/`), gộp toàn bộ chức năng
của 8 biến thể (hợp của tất cả các hàm), kèm test đầy đủ.

Yêu cầu:
- Trước khi gộp: **lập bảng so sánh** 8 biến thể (hàm nào có ở đâu, chữ ký/hành vi
  có khác nhau không) — đính kèm trong mô tả PR.
- Nếu hai biến thể có cùng tên hàm nhưng **hành vi khác nhau**, nêu rõ và đề xuất
  cách xử lý (đừng chọn bừa — giải thích trade-off).
- Giữ test xanh; bổ sung test cho phần hợp nhất.

**Đánh giá:** đọc-hiểu code lạ, tư duy DRY, phát hiện khác biệt tinh vi, giao tiếp trade-off.

---

## T3 — CLI sinh mã (`codegen/`)

`codegen/editorial/` và `codegen/nrcs/` chứa các hàm sinh mã dạng
`PREFIX-YYYYMMDD-XXXXXX` (hậu tố ngẫu nhiên base32, dùng `crypto/rand`).

**Việc cần làm:** viết một CLI nhỏ `cmd/gencode` tái sử dụng các hàm này.

Ví dụ giao diện mong muốn (bạn được tự thiết kế, miễn hợp lý + có `--help`):

```bash
gencode document            # in một document code
gencode story --count 5     # in 5 story code
gencode package --json      # in dạng JSON
```

Yêu cầu: parse flag chuẩn, có `--help`, xử lý lỗi rõ ràng, kèm `README` ngắn cho CLI +
test cho phần logic bạn tự viết.

**Đánh giá:** thiết kế từ yêu cầu mơ hồ, tính tự chủ, tài liệu, hoàn thiện end-to-end.

---

## T4 — Test + tìm edge-case cho `textcode/program_code.go`

`textcode/program_code.go` sinh/so sánh mã từ tiêu đề tiếng Việt: bỏ dấu, chuẩn hóa,
lấy ký tự. File này **chưa có test**.

**Việc cần làm:** viết test toàn diện + chỉ ra (bằng test) các edge-case dễ sai.

Gợi ý nơi dễ có vấn đề:
- Chữ **Đ/đ** — ký tự này **không phân rã** qua chuẩn hóa Unicode NFD như các chữ có
  dấu khác. Hãy kiểm tra kỹ kết quả với từ chứa Đ/đ.
- Tiêu đề có nhiều khoảng trắng, ký tự đặc biệt, chuỗi rỗng, chỉ toàn dấu.

**Đánh giá:** tính tỉ mỉ với Unicode/chuỗi tiếng Việt, tư duy phản-ví-dụ.

---

## T5 — Truy lỗi (bug hunt) trong `readingpace/`

`readingpace/reading_pace.go` ước lượng thời lượng đọc (giây) của một đoạn text tiếng
Việt: chuẩn hóa (bỏ markup, mở rộng số/viết tắt sang dạng đọc) rồi tính theo WPM.

**Báo lỗi từ QA (mơ hồ, như thực tế):**
> "Một số bản tin hiển thị thời lượng đọc bị sai — có vẻ liên quan tới tin có **số**
> và tin gán **từ điển phát âm** cho người dẫn. Nhờ kiểm tra lại."

**Việc cần làm:**
1. Viết test để **tái hiện** lỗi (file này hiện không kèm test — bạn tự viết).
2. Tìm **nguyên nhân gốc**, sửa.
3. Để lại **regression test** chứng minh đã sửa đúng.
4. Mô tả trong PR: bạn đã khoanh vùng lỗi thế nào, vì sao fix đó đúng.

> Gợi ý: đừng chỉ test các trường hợp "tròn trịa". Lỗi nằm ở những giá trị mà
> test hiển nhiên dễ bỏ qua.

**Đánh giá:** năng lực debug/lập luận trong code có sẵn — kỹ năng quan trọng nhất khi
vào dự án thật. (Đây là task phân loại mạnh nhất.)

---

## T6 — Benchmark + tối ưu `markdown/`

`markdown/inline_markdown.go` có `StripInlineMarkdown` — bộ parse một lượt (single-pass)
lược bỏ `**đậm**`, `*nghiêng*`, `__gạch chân__` và xử lý escape. Đã có test.

**Việc cần làm:**
1. Viết benchmark (`testing.B`) cho `StripInlineMarkdown` với vài dạng input
   (text dài, nhiều mark, nhiều escape).
2. Đo, tìm điểm tốn kém, **tối ưu** mà **không đổi hành vi** (test cũ phải xanh).
3. Báo cáo số liệu before/after trong PR (kèm cách đo).

**Đánh giá:** ý thức hiệu năng, kỹ năng benchmark Go, kỷ luật "không phá behavior".

---

## T7 — Fuzz test (`markdown/` + `date/`)

Go có fuzzing gốc (`testing.F` / `f.Fuzz`).

**Việc cần làm:** viết fuzz test cho:
- `StripInlineMarkdown` (`markdown/`) — tìm panic, vòng lặp vô hạn, hoặc kết quả bất nhất
  (vd: strip hai lần phải ra cùng kết quả với input không có mark).
- `ParseDateInVN` (chọn một biến thể `date/`) — input rác không được panic.

Yêu cầu: nêu **invariant** mà fuzz kiểm chứng; nếu phát hiện crash, ghi lại corpus +
mô tả.

**Đánh giá:** tư duy robustness, biết dùng công cụ Go hiện đại.

---

## Tổng quan đánh giá

| Task | Trục năng lực chính |
|---|---|
| T1 | Go cơ bản, table-driven test, edge-case |
| T2 | Đọc code lạ, DRY, trade-off, giao tiếp |
| T3 | Thiết kế end-to-end, tài liệu, tự chủ |
| T4 | Tỉ mỉ Unicode/chuỗi tiếng Việt |
| T5 | Debug/lập luận (quan trọng nhất) |
| T6 | Hiệu năng, benchmark |
| T7 | Robustness, fuzzing |

Không bắt buộc làm hết — báo trước nếu thiếu thời gian. Chất lượng + độ hiểu > số lượng.
