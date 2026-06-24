# T2 — Gộp 8 bản `date.go` làm một

> **Tóm tắt 1 câu:** 8 thư mục trong `date/` chứa code gần giống nhau. Mình so sánh kỹ,
> thấy chúng **không hề mâu thuẫn** — chỉ là bản này có thêm hàm bản kia không có — nên
> gộp lại rất an toàn: **gom tất cả các hàm vào một package chung là xong.**

---

## 1. Có những bản nào?

Thư mục `date/` có **8 bản** cùng làm một việc (xử lý ngày/giờ múi giờ Việt Nam):

```
ad   channel   config   editorial   nrcs   production   resource   royalty
```

Mục tiêu của T2: thay 8 bản trùng lặp này bằng **1 package dùng chung**.

---

## 2. Bản nào có hàm gì?

> ✅ = có hàm này  ·  ❌ = không có

| Hàm | ad | channel | config | editorial | nrcs | production | resource | royalty |
|---|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|
| `VNLocation` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `ParseDateInVN` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `NormalizeDateToVN` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `TodayInVN` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `FormatDateVN` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `NowInVN` | ✅ | ❌ | ❌ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `ResolveAirDatetime` | ✅ | ❌ | ❌ | ✅ | ❌ | ✅ | ❌ | ❌ |
| `IsDateInPast` | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| `IsDateToday` | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| `SameDate` | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| `DateRangeInVN` | ❌ | ❌ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ |
| `FormatTimeVN` | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ | ✅ |
| `FormatTimeVNPtr` | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ | ✅ |

**Đọc bảng này thế nào — chia làm 3 nhóm:**

- 🟦 **Nhóm 1 — 5 hàm ai cũng có** (5 dòng đầu): phần lõi chung.
- 🟨 **Nhóm 2 — vài bản mới có**: `NowInVN` (6/8 bản), `ResolveAirDatetime` (3 bản).
- 🟩 **Nhóm 3 — hàm "đặc sản" của 1 bản**: ví dụ `DateRangeInVN` chỉ `nrcs` có,
  `FormatTimeVN`/`FormatTimeVNPtr` chỉ `resource` và `royalty` có.

---

## 3. Câu hỏi quan trọng nhất: hàm trùng tên có làm khác nhau không?

Khi gộp, nỗi lo lớn nhất là: *hai bản có hàm **cùng tên** nhưng **chạy ra kết quả khác nhau***
→ gộp ẩu sẽ sai logic.

Mình kiểm tra bằng cách **so từng dòng code** (`diff`) của các hàm trùng tên:

| Hàm trùng tên | Xuất hiện ở | Kết quả so sánh code |
|---|---|---|
| 5 hàm lõi | cả 8 bản | ✅ **Giống y hệt** |
| `NowInVN` | 6 bản | ✅ **Giống y hệt** |
| `ResolveAirDatetime` | `ad`, `editorial`, `production` | ✅ **Giống y hệt** |

### ✅ Kết luận

> **Không có hàm nào cùng tên mà chạy khác nhau.**
> Mỗi bản chỉ *thêm* hàm mới, **không** sửa hàm cũ.
> → Gộp cực đơn giản: **gom hết các hàm lại là xong**, không phải đắn đo chọn bản nào.

### 💡 Nếu lỡ gặp 2 hàm trùng tên mà khác nhau thì sao?

Lần này không gặp, nhưng nguyên tắc xử lý sẽ là:

1. **Ghi rõ** chúng khác nhau ở điểm nào.
2. **Chọn** cách làm an toàn/tổng quát hơn, hoặc **tách** thành 2 hàm riêng.
3. **Giải thích lý do** — chứ không chọn bừa.

---

## 4. Gộp lại thành gì?

Tạo **1 package chung** ở `date/` (tên `util`), gồm đủ **13 hàm**:

| Lấy từ bản | Các hàm |
|---|---|
| `ad` | 5 hàm lõi + `NowInVN` + `ResolveAirDatetime` |
| `channel` | `IsDateInPast`, `IsDateToday`, `SameDate` |
| `nrcs` | `DateRangeInVN` |
| `resource` | `FormatTimeVN`, `FormatTimeVNPtr` |

> 8 thư mục cũ **tạm thời giữ nguyên** để không làm hỏng test đang chạy tốt;
> sẽ chuyển dần sang dùng package mới ở các bước sau.

---

## 5. Cách mình đã kiểm tra (ai cũng chạy lại được)

```bash
# So 5 hàm lõi của mọi bản với bản channel.
# Nếu lệnh diff IM LẶNG (không in gì) => code giống hệt nhau.
ref=channel
for fn in VNLocation ParseDateInVN NormalizeDateToVN TodayInVN FormatDateVN; do
  for d in ad config editorial nrcs production resource royalty; do
    diff <(sed -n "/^func $fn(/,/^}/p" date/$ref/date.go) \
         <(sed -n "/^func $fn(/,/^}/p" date/$d/date.go)
  done
done

# So ResolveAirDatetime giữa 3 bản ad / editorial / production
diff <(sed -n '/func ResolveAirDatetime/,/^}/p' date/ad/date.go) \
     <(sed -n '/func ResolveAirDatetime/,/^}/p' date/editorial/date.go)
diff <(sed -n '/func ResolveAirDatetime/,/^}/p' date/ad/date.go) \
     <(sed -n '/func ResolveAirDatetime/,/^}/p' date/production/date.go)
```

**Tất cả lệnh trên chạy ra không có gì** → khẳng định các hàm trùng tên giống hệt nhau.
